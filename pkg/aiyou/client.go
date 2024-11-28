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
// File: pkg/aiyou/client.go

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
type ClientOption func(*Client)

// Nouvelle option de configuration
func WithRateLimiter(config RateLimiterConfig) ClientOption {
	return func(c *Client) {
		c.rateLimiter = NewRateLimiter(config, c.logger)
	}
}

// NewClient creates a new instance of Client with the given email and password.

func NewClient(email, password string, options ...ClientOption) (*Client, error) {
	client := &Client{
		baseURL:      "https://ai.dragonflygroup.fr",
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		maxRetries:   3,
		initialDelay: time.Second,
		logger:       NewDefaultLogger(os.Stderr),
	}

	for _, option := range options {
		option(client)
	}

	client.safeLog = SafeLog(client.logger)

	auth := NewJWTAuthenticator(email, password, client.baseURL, client.httpClient, client.logger)
	client.auth = auth

	return client, nil
}

// WithLogger sets a custom logger for the client.
// If not provided, a default logger writing to os.Stderr will be used.
func WithLogger(logger Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithBaseURL sets the base URL for API requests.
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

// AuthenticatedRequest performs an authenticated request to the API.
// It handles authentication, retries, and logs the request and response details
// while masking sensitive information.
func (c *Client) AuthenticatedRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	c.safeLog(DEBUG, "Preparing authenticated request: %s %s", method, path)
   
	// Appliquer le rate limiting avant la tentative de requête
	if c.rateLimiter != nil {
	if err := c.rateLimiter.Wait(ctx); err != nil {
	c.safeLog(WARN, "Client-side rate limit exceeded: %v", err)
	waitTime := c.rateLimiter.GetWaitTime()
	return nil, &RateLimitError{
	RetryAfter: int(waitTime.Seconds()),
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
   
	// Gestion du rate limiting côté serveur
	if resp.StatusCode == http.StatusTooManyRequests {
	c.safeLog(WARN, "Server-side rate limit exceeded, retrying after 60 seconds")
	return &RateLimitError{
	RetryAfter: 60,
	IsClientSide: false,
	}
	}
   
	c.safeLog(INFO, "Request completed with status: %v", resp.StatusCode)
	return nil
	})
   
	if err != nil {
	return nil, err
	}
   
	return resp, nil
   }
// SetBaseURL sets the base URL for API requests.
func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

// SetLogger sets the logger for the client.
func (c *Client) SetLogger(logger Logger) {
	c.logger = logger
}

// CreateChatCompletion is a convenience method for creating a chat completion
func (c *Client) CreateChatCompletion(ctx context.Context, messages []Message, assistantID string) (*ChatCompletionResponse, error) {
	req := ChatCompletionRequest{
		Messages: messages,
		AssistantID: assistantID,
		Temperature: 0.7,
		TopP: 1,
	}
	return c.ChatCompletion(ctx, req)
}

// CreateChatCompletionStream is a convenience method for creating a streaming chat completion
func (c *Client) CreateChatCompletionStream(ctx context.Context, messages []Message, assistantID string) (*StreamReader, error) {
	req := ChatCompletionRequest{
		Messages:    messages,
		AssistantID: assistantID,
		// Set default values for other fields
		Temperature:  0.7,
		TopP:         1,
		PromptSystem: "",
		Stream:       true,
	}
	return c.ChatCompletionStream(ctx, req)
}

// WithRetry configure les options de retry pour le client
func WithRetry(maxRetries int, initialDelay time.Duration) ClientOption {
	return func(c *Client) {
		c.maxRetries = maxRetries
		c.initialDelay = initialDelay
	}
}
