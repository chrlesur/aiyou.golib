/*
Copyright (C) 2024 Cloud Temple

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

// Package aiyou provides a client for interacting with the AI.YOU API from Cloud Temple.
package aiyou

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Client represents a client for the AI.YOU API.
type Client struct {
	baseURL      string
	httpClient   *http.Client
	auth         Authenticator
	maxRetries   int
	initialDelay time.Duration
	logger       Logger
	safeLog      func(level LogLevel, format string, args ...interface{})
	rateLimiter  *RateLimiter
}

// ClientOption is a function type to modify Client.
type ClientOption func(*Client) error

// NewClient creates a new instance of Client with the given options.
// At least one authentication method (email/password or bearer token) must be provided.
func NewClient(options ...ClientOption) (*Client, error) {
	client := &Client{
		baseURL:      "https://ai.dragonflygroup.fr",
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		maxRetries:   3,
		initialDelay: time.Second,
		logger:       NewDefaultLogger(os.Stderr),
	}

	var err error
	for _, option := range options {
		if err = option(client); err != nil {
			return nil, fmt.Errorf("failed to apply client option: %w", err)
		}
	}

	client.safeLog = SafeLog(client.logger)

	// Vérifier qu'une méthode d'authentification a été configurée
	if client.auth == nil {
		return nil, fmt.Errorf("no authentication method provided: use WithEmailPassword or WithBearerToken")
	}

	return client, nil
}

// WithEmailPassword configures the client to use email/password authentication
func WithEmailPassword(email, password string) ClientOption {
	return func(c *Client) error {
		if email == "" || password == "" {
			return fmt.Errorf("email and password cannot be empty")
		}
		c.auth = NewJWTAuthenticator(email, password, c.baseURL, c.httpClient, c.logger)
		c.logger.Debugf("Configured client with email/password authentication for: %s", maskSensitiveInfo(email))
		return nil
	}
}

// WithBearerToken configures the client to use bearer token authentication
func WithBearerToken(token string) ClientOption {
	return func(c *Client) error {
		if token == "" {
			return fmt.Errorf("bearer token cannot be empty")
		}
		c.auth = NewBearerAuthenticator(token, c.logger)
		c.logger.Debugf("Configured client with bearer token authentication")
		return nil
	}
}

// WithLogger sets a custom logger for the client.
func WithLogger(logger Logger) ClientOption {
	return func(c *Client) error {
		if logger == nil {
			return fmt.Errorf("logger cannot be nil")
		}
		c.logger = logger
		if c.auth != nil {
			switch auth := c.auth.(type) {
			case *JWTAuthenticator:
				auth.SetLogger(logger)
			case *BearerAuthenticator:
				auth.SetLogger(logger)
			}
		}
		return nil
	}
}

// WithBaseURL sets the base URL for API requests.
func WithBaseURL(url string) ClientOption {
	return func(c *Client) error {
		if url == "" {
			return fmt.Errorf("base URL cannot be empty")
		}
		c.baseURL = url
		return nil
	}
}

// WithRetry configures retry options for the client
func WithRetry(maxRetries int, initialDelay time.Duration) ClientOption {
	return func(c *Client) error {
		if maxRetries < 0 {
			return fmt.Errorf("maxRetries cannot be negative")
		}
		if initialDelay < 0 {
			return fmt.Errorf("initialDelay cannot be negative")
		}
		c.maxRetries = maxRetries
		c.initialDelay = initialDelay
		return nil
	}
}

// WithRateLimiter configures the rate limiter for the client
func WithRateLimiter(config RateLimiterConfig) ClientOption {
	return func(c *Client) error {
		c.rateLimiter = NewRateLimiter(config, c.logger)
		return nil
	}
}

// SetBearerToken updates the bearer token if using bearer token authentication
func (c *Client) SetBearerToken(token string) error {
	if token == "" {
		return fmt.Errorf("bearer token cannot be empty")
	}

	auth, ok := c.auth.(*BearerAuthenticator)
	if !ok {
		return fmt.Errorf("client is not configured for bearer token authentication")
	}

	auth.SetToken(token)
	c.logger.Infof("Bearer token has been updated")
	return nil
}

// AuthenticatedRequest performs an authenticated request to the API.
func (c *Client) AuthenticatedRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	c.safeLog(DEBUG, "Preparing authenticated request: %s %s", method, path)

	if c.rateLimiter != nil {
		if err := c.rateLimiter.Wait(ctx); err != nil {
			c.safeLog(WARN, "Client-side rate limit exceeded: %v", err)
			waitTime := c.rateLimiter.GetWaitTime()
			return nil, &RateLimitError{
				RetryAfter:   int(waitTime.Seconds()),
				IsClientSide: true,
			}
		}
	}

	var resp *http.Response
	err := retryOperation(ctx, c.logger, c.maxRetries, c.initialDelay, func() error {
		if err := c.auth.Authenticate(ctx); err != nil {
			c.safeLog(ERROR, "Authentication failed: %v", err)
			return &AuthenticationError{Message: err.Error()}
		}

		req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
		if err != nil {
			c.safeLog(ERROR, "Failed to create request: %v", err)
			return fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+c.auth.Token())
		req.Header.Set("Content-Type", "application/json")

		c.safeLog(DEBUG, "Sending request to %s", req.URL)
		resp, err = c.httpClient.Do(req)
		if err != nil {
			c.safeLog(ERROR, "Request failed: %v", err)
			return &NetworkError{Err: err}
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			c.safeLog(WARN, "Server-side rate limit exceeded, retrying after 60 seconds")
			return &RateLimitError{
				RetryAfter:   60,
				IsClientSide: false,
			}
		}

		c.safeLog(INFO, "Request completed with status: %d", resp.StatusCode)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// SetBaseURL sets the base URL for API requests
func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

// SetLogger sets the logger for the client
func (c *Client) SetLogger(logger Logger) {
	c.logger = logger
	if c.auth != nil {
		switch auth := c.auth.(type) {
		case *JWTAuthenticator:
			auth.SetLogger(logger)
		case *BearerAuthenticator:
			auth.SetLogger(logger)
		}
	}
	c.safeLog = SafeLog(logger)
}

// CreateChatCompletion is a helper method that wraps ChatCompletion
func (c *Client) CreateChatCompletion(ctx context.Context, messages []Message, assistantID string) (*ChatCompletionResponse, error) {
	req := ChatCompletionRequest{
		Messages:    messages,
		AssistantID: assistantID,
		Stream:      false,
	}
	return c.ChatCompletion(ctx, req)
}

// CreateChatCompletionStream is a helper method that wraps ChatCompletionStream
func (c *Client) CreateChatCompletionStream(ctx context.Context, messages []Message, assistantID string) (*StreamReader, error) {
	req := ChatCompletionRequest{
		Messages:    messages,
		AssistantID: assistantID,
		Stream:      true,
	}
	return c.ChatCompletionStream(ctx, req)
}
