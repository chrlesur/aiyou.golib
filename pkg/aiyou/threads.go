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

// File: pkg/aiyou/threads.go
package aiyou

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// GetUserThreads récupère la liste des threads de l'utilisateur avec pagination et filtrage
func (c *Client) GetUserThreads(ctx context.Context, params *UserThreadsParams) (*UserThreadsOutput, error) {
	endpoint := "/api/v1/user/threads"
	if params != nil {
		query := url.Values{}
		if params.Page > 0 {
			query.Set("page", strconv.Itoa(params.Page))
		}
		if params.ItemsPerPage > 0 {
			query.Set("itemsPerPage", strconv.Itoa(params.ItemsPerPage))
		}
		if len(query) > 0 {
			endpoint = fmt.Sprintf("%s?%s", endpoint, query.Encode())
		}
	}

	resp, err := c.AuthenticatedRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var threadsOutput UserThreadsOutput
	if err := json.NewDecoder(resp.Body).Decode(&threadsOutput); err != nil {
		if err == io.EOF {
			c.logger.Errorf("Empty response body")
			return nil, err // Retourne directement l'erreur EOF
		}
		c.logger.Errorf("Failed to decode threads response: %v", err)
		return nil, fmt.Errorf("failed to decode threads response: %w", err)
	}

	c.logger.Infof("Successfully retrieved %d threads", len(threadsOutput.Threads))
	return &threadsOutput, nil
}

// DeleteThread supprime un thread spécifique
func (c *Client) DeleteThread(ctx context.Context, threadID string) error {
	endpoint := fmt.Sprintf("/api/v1/threads/%s", threadID)
	c.logger.Debugf("Deleting thread: %s", threadID)

	resp, err := c.AuthenticatedRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		c.logger.Errorf("Failed to delete thread: %v", err)
		return fmt.Errorf("failed to delete thread: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	c.logger.Infof("Successfully deleted thread: %s", threadID)
	return nil
}
