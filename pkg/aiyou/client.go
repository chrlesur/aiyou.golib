/*
Copyright (C) 2023 Cloud Temple

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
	"log"
	"net/http"
	"time"
    "bufio"
    "bytes"
)

// Client represents a client for the AI.YOU API.
type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     *log.Logger
	auth       Authenticator
}

// NewClient creates a new instance of Client with the given email and password.
func NewClient(email, password string, options ...ClientOption) (*Client, error) {
	client := &Client{
		baseURL:    "https://ai.dragonflygroup.fr", // URL de base par défaut
		httpClient: &http.Client{Timeout: 30 * time.Second},
		logger:     log.New(log.Writer(), "aiyou: ", log.LstdFlags),
	}

	for _, option := range options {
		option(client)
	}

	auth := NewJWTAuthenticator(email, password, client.baseURL, client.httpClient)
	client.auth = auth

	return client, nil
}

// AuthenticatedRequest performs an authenticated request to the API.
func (c *Client) AuthenticatedRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	if err := c.auth.Authenticate(ctx); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.auth.Token())

	// Ajoutez ici la logique pour ajouter le corps de la requête si nécessaire

	return c.httpClient.Do(req)
}

// SetBaseURL sets the base URL for API requests.
func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

// SetLogger sets the logger for the client.
func (c *Client) SetLogger(logger *log.Logger) {
	c.logger = logger
}

// CreateChatCompletion is a convenience method for creating a chat completion
func (c *Client) CreateChatCompletion(ctx context.Context, messages []Message, assistantID string) (*ChatCompletionResponse, error) {
    req := ChatCompletionRequest{
        Messages:    messages,
        AssistantID: assistantID,
        // Set default values for other fields
        Temperature:  0.7,
        TopP:         1,
        PromptSystem: "",
        Stream:       false,
    }
    return c.ChatCompletion(ctx, req)
}

// CreateChatCompletionStream is a convenience method for creating a streaming chat completion
func (c *Client) CreateChatCompletionStream(ctx context.Context, messages []Message, assistantID string) (io.ReadCloser, error) {
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

