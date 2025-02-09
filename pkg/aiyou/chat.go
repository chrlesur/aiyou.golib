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
// File: pkg/aiyou/chat.go

package aiyou

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// StreamReader helps read and process the streaming response
type StreamReader struct {
	reader *bufio.Reader
	closer io.Closer
	logger Logger
}

// NewStreamReader creates a new StreamReader
func NewStreamReader(r io.ReadCloser, logger Logger) *StreamReader {
	return &StreamReader{
		reader: bufio.NewReader(r),
		closer: r,
		logger: logger,
	}
}

// ChatCompletion attempts non-streaming first and falls back to aggregated streaming if needed
func (c *Client) ChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	c.logger.Debugf("Starting ChatCompletion request")

	// Première tentative en mode non-streaming
	req.Stream = false
	jsonData, err := json.Marshal(req)
	if err != nil {
		c.logger.Errorf("Failed to marshal request: %v", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	c.logger.Debugf("Attempting non-streaming request first")
	resp, err := c.AuthenticatedRequest(ctx, "POST", "/api/v1/chat/completions", bytes.NewReader(jsonData))
	if err != nil {
		c.logger.Errorf("Non-streaming request failed: %v", err)
		// Vérifier si c'est une erreur de rate limit avant de faire le fallback
		if _, isRateLimit := err.(*RateLimitError); isRateLimit {
			return nil, err // Propager directement l'erreur de rate limit
		}
		return c.fallbackToStreamingAggregation(ctx, req)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Errorf("Failed to read response body: %v", err)
		return c.fallbackToStreamingAggregation(ctx, req)
	}

	// Vérifier si c'est une réponse d'erreur qui nécessite le fallback
	var errorResp struct {
		Object  string `json:"object"`
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    int    `json:"code"`
	}

	if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Object == "error" {
		if strings.Contains(errorResp.Message, "Stream options") {
			c.logger.Debugf("Detected streaming options error, falling back to streaming aggregation")
			return c.fallbackToStreamingAggregation(ctx, req)
		}
		// Autres erreurs API
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("API Error: %s - %s", errorResp.Type, errorResp.Message),
		}
	}

	// Si on arrive ici, le mode non-streaming a fonctionné
	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		c.logger.Errorf("Failed to decode response: %v", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	c.logger.Infof("ChatCompletion request successful using non-streaming mode")
	return &chatResp, nil
}

// fallbackToStreamingAggregation handles the fallback to streaming mode with aggregation
func (c *Client) fallbackToStreamingAggregation(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	c.logger.Infof("Falling back to streaming aggregation mode")

	req.Stream = true
	stream, err := c.ChatCompletionStream(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to start stream in fallback: %w", err)
	}
	defer stream.Close()

	var fullResponse ChatCompletionResponse
	var aggregatedContent strings.Builder

	for {
		chunk, err := stream.ReadChunk()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading stream in fallback: %w", err)
		}
		if chunk == nil || len(chunk.Choices) == 0 {
			continue
		}

		if fullResponse.ID == "" {
			fullResponse = *chunk
		}

		choice := chunk.Choices[0]
		if choice.Delta != nil && choice.Delta.Content != "" {
			aggregatedContent.WriteString(choice.Delta.Content)
		}
	}

	if len(fullResponse.Choices) > 0 {
		fullResponse.Choices[0].Message = Message{
			Role: "assistant",
			Content: []ContentPart{
				{
					Type: "text",
					Text: aggregatedContent.String(),
				},
			},
		}
	}

	c.logger.Infof("Successfully completed request using streaming aggregation fallback")
	return &fullResponse, nil
}

// ChatCompletionStream sends a streaming chat completion request
func (c *Client) ChatCompletionStream(ctx context.Context, req ChatCompletionRequest) (*StreamReader, error) {
	req.Stream = true

	jsonData, err := json.Marshal(req)
	if err != nil {
		c.logger.Errorf("Failed to marshal request: %v", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.AuthenticatedRequest(ctx, "POST", "/api/v1/chat/completions", bytes.NewReader(jsonData))
	if err != nil {
		c.logger.Errorf("ChatCompletionStream request failed: %v", err)
		// Propager directement l'erreur
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		c.logger.Warnf("Unexpected status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return NewStreamReader(resp.Body, c.logger), nil
}

// ReadChunk reads and processes a single chunk from the stream
// Dans pkg/aiyou/chat.go
func (sr *StreamReader) ReadChunk() (*ChatCompletionResponse, error) {
	line, err := sr.reader.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		sr.logger.Errorf("Error reading stream: %v", err)
		return nil, err
	}

	// Nettoyer la ligne
	line = bytes.TrimSpace(line)
	if len(line) == 0 {
		return sr.ReadChunk() // Continuer à lire si la ligne est vide
	}

	// Vérifier si c'est la fin du stream
	if string(line) == "[DONE]" {
		return nil, io.EOF
	}

	// Vérifier et traiter le préfixe "data: "
	if bytes.HasPrefix(line, []byte("data: ")) {
		line = bytes.TrimPrefix(line, []byte("data: "))
	} else {
		sr.logger.Debugf("Skipping non-data line: %s", string(line))
		return sr.ReadChunk() // Ignorer les lignes sans préfixe "data: "
	}

	// Vérifier si la ligne est vide après le trim
	if len(bytes.TrimSpace(line)) == 0 {
		return sr.ReadChunk()
	}

	var chunk ChatCompletionResponse
	if err := json.Unmarshal(line, &chunk); err != nil {
		sr.logger.Errorf("Failed to unmarshal chunk: %v, raw data: %s", err, string(line))
		// Option 1 : continuer à lire
		return sr.ReadChunk()
		// Option 2 : retourner l'erreur
		// return nil, fmt.Errorf("failed to unmarshal chunk: %w, raw data: %s", err, string(line))
	}

	return &chunk, nil
}

// SetLogger sets a custom logger for the StreamReader
func (sr *StreamReader) SetLogger(logger Logger) {
	sr.logger = logger
}

// Close closes the underlying reader
func (sr *StreamReader) Close() error {
	return sr.closer.Close()
}
