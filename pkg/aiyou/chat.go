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

// ChatCompletion sends a chat completion request and returns the response
func (c *Client) ChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	c.logger.Debugf("Starting ChatCompletion request")

	jsonData, err := json.Marshal(req)
	if err != nil {
		c.logger.Errorf("Failed to marshal request: %v", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.AuthenticatedRequest(ctx, "POST", "/api/v1/chat/completions", bytes.NewReader(jsonData))
	if err != nil {
		c.logger.Errorf("ChatCompletion request failed: %v", err)
		return nil, fmt.Errorf("chat completion request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Warnf("Unexpected status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var chatResp ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		c.logger.Errorf("Failed to decode response: %v", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	c.logger.Infof("ChatCompletion request successful")
	return &chatResp, nil
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
		return nil, fmt.Errorf("chat completion stream request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		c.logger.Warnf("Unexpected status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return NewStreamReader(resp.Body, c.logger), nil
}

// ReadChunk reads and processes a single chunk from the stream
func (sr *StreamReader) ReadChunk() (*ChatCompletionResponse, error) {
	line, err := sr.reader.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			sr.logger.Infof("End of stream reached")
		} else {
			sr.logger.Errorf("Error reading stream: %v", err)
		}
		return nil, err
	}

	line = bytes.TrimPrefix(line, []byte(""))
	line = bytes.TrimSpace(line)

	if len(line) == 0 {
		return nil, nil
	}

	var chunk ChatCompletionResponse
	if err := json.Unmarshal(line, &chunk); err != nil {
		sr.logger.Errorf("Failed to unmarshal chunk: %v", err)
		return nil, fmt.Errorf("failed to unmarshal chunk: %w", err)
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
