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
	var resp *ChatCompletionResponse
	err := retryOperation(ctx, c.logger, c.maxRetries, c.initialDelay, func() error {
		var err error
		endpoint := "/api/v1/chat/completions"
		jsonData, err := json.Marshal(req)
		if err != nil {
			c.logger.Errorf("Failed to marshal request: %v", err)
			return fmt.Errorf("failed to marshal request: %w", err)
		}

		c.logger.Debugf("Sending ChatCompletion request to %s", endpoint)
		httpResp, err := c.AuthenticatedRequest(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonData))
		if err != nil {
			c.logger.Errorf("ChatCompletion request failed: %v", err)
			return &NetworkError{Err: err}
		}
		defer httpResp.Body.Close()

		if httpResp.StatusCode != http.StatusOK {
			c.logger.Warnf("Unexpected status code: %d", httpResp.StatusCode)
			return &APIError{StatusCode: httpResp.StatusCode, Message: fmt.Sprintf("unexpected status code: %d", httpResp.StatusCode)}
		}

		resp = &ChatCompletionResponse{}
		if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
			c.logger.Errorf("Failed to decode response: %v", err)
			return fmt.Errorf("failed to decode response: %w", err)
		}

		c.logger.Infof("ChatCompletion request successful")
		return nil
	})

	return resp, err
}

// ChatCompletionStream sends a streaming chat completion request
func (c *Client) ChatCompletionStream(ctx context.Context, req ChatCompletionRequest) (*StreamReader, error) {
	c.logger.Debugf("Starting ChatCompletionStream request")
	var streamReader *StreamReader
	err := retryOperation(ctx, c.logger, c.maxRetries, c.initialDelay, func() error {
		endpoint := "/api/v1/chat/completions"
		req.Stream = true
		jsonData, err := json.Marshal(req)
		if err != nil {
			c.logger.Errorf("Failed to marshal request: %v", err)
			return fmt.Errorf("failed to marshal request: %w", err)
		}

		c.logger.Debugf("Sending ChatCompletionStream request to %s", endpoint)
		resp, err := c.AuthenticatedRequest(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonData))
		if err != nil {
			c.logger.Errorf("ChatCompletionStream request failed: %v", err)
			return &NetworkError{Err: err}
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			c.logger.Warnf("Unexpected status code: %d", resp.StatusCode)
			return &APIError{StatusCode: resp.StatusCode, Message: fmt.Sprintf("unexpected status code: %d", resp.StatusCode)}
		}

		streamReader = NewStreamReader(resp.Body, c.logger)
		c.logger.Infof("ChatCompletionStream request successful, streaming started")
		return nil
	})

	return streamReader, err
}

// ReadChunk reads and processes a single chunk from the stream
func (sr *StreamReader) ReadChunk() (*ChatCompletionResponse, error) {
    line, err := sr.reader.ReadBytes('\n')
    if err != nil {
        if err == io.EOF {
            if sr.logger != nil {
                sr.logger.Infof("End of stream reached")
            }
        } else {
            if sr.logger != nil {
                sr.logger.Errorf("Error reading stream: %v", err)
            }
        }
        return nil, err
    }

    if sr.logger != nil {
        sr.logger.Debugf("Raw chunk data: %s", string(line))
    }

    line = bytes.TrimPrefix(line, []byte("data: "))
    line = bytes.TrimSpace(line)

    if len(line) == 0 {
        return nil, nil // Skip empty lines
    }

    var chunk ChatCompletionResponse
    if err := json.Unmarshal(line, &chunk); err != nil {
        if sr.logger != nil {
            sr.logger.Errorf("Failed to unmarshal chunk: %v", err)
        }
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
