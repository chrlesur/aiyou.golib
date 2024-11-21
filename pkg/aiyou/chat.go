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

// Fichier: pkg/aiyou/chat.go

package aiyou

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
	"bufio"
    "bytes"
)

// StreamReader helps read and process the streaming response
type StreamReader struct {
    reader *bufio.Reader
}

// ChatCompletion sends a chat completion request and returns the response
func (c *Client) ChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
    endpoint := "/api/v1/chat/completions"
    jsonData, err := json.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }

    resp, err := c.AuthenticatedRequest(ctx, http.MethodPost, endpoint, strings.NewReader(string(jsonData)))
    if err != nil {
        return nil, fmt.Errorf("failed to send chat completion request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    var chatResp ChatCompletionResponse
    if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    return &chatResp, nil
}

// NewStreamReader creates a new StreamReader
func NewStreamReader(r io.Reader) *StreamReader {
    return &StreamReader{reader: bufio.NewReader(r)}
}

// ReadChunk reads and processes a single chunk from the stream
func (sr *StreamReader) ReadChunk() (*ChatCompletionResponse, error) {
    line, err := sr.reader.ReadBytes('\n')
    if err != nil {
        return nil, err
    }

    // Remove "data: " prefix if present
    line = bytes.TrimPrefix(line, []byte("data: "))

    var chunk ChatCompletionResponse
    if err := json.Unmarshal(line, &chunk); err != nil {
        return nil, fmt.Errorf("failed to unmarshal chunk: %w", err)
    }

    return &chunk, nil
}

// Modification de la fonction ChatCompletionStream existante
func (c *Client) ChatCompletionStream(ctx context.Context, req ChatCompletionRequest) (*StreamReader, error) {
    endpoint := "/api/v1/chat/completions"
    req.Stream = true
    jsonData, err := json.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }

    resp, err := c.AuthenticatedRequest(ctx, http.MethodPost, endpoint, strings.NewReader(string(jsonData)))
    if err != nil {
        return nil, fmt.Errorf("failed to send chat completion stream request: %w", err)
    }

    if resp.StatusCode != http.StatusOK {
        resp.Body.Close()
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    return NewStreamReader(resp.Body), nil
}