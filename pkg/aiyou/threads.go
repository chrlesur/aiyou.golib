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
	"net/url"
	"time"
)

// GetUserThreads récupère la liste des threads de l'utilisateur avec pagination et filtrage
func (c *Client) GetUserThreads(ctx context.Context, filter *ThreadFilter) (*ThreadsResponse, error) {
	endpoint := "/api/v1/user/threads"
	c.logger.Debugf("Fetching user threads with filter: %+v", filter)

	// Construction des paramètres de requête
	query := url.Values{}
	if filter != nil {
		if filter.AssistantID != "" {
			query.Set("assistantId", filter.AssistantID)
		}
		if !filter.StartDate.IsZero() {
			query.Set("startDate", filter.StartDate.Format(time.RFC3339))
		}
		if !filter.EndDate.IsZero() {
			query.Set("endDate", filter.EndDate.Format(time.RFC3339))
		}
		if filter.Page > 0 {
			query.Set("page", fmt.Sprintf("%d", filter.Page))
		}
		if filter.Limit > 0 {
			query.Set("limit", fmt.Sprintf("%d", filter.Limit))
		}
	}

	// Ajout des paramètres à l'URL si présents
	if len(query) > 0 {
		endpoint = fmt.Sprintf("%s?%s", endpoint, query.Encode())
	}

	resp, err := c.AuthenticatedRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		c.logger.Errorf("Failed to fetch threads: %v", err)
		return nil, fmt.Errorf("failed to fetch threads: %w", err)
	}
	defer resp.Body.Close()

	var threadsResp ThreadsResponse
	if err := json.NewDecoder(resp.Body).Decode(&threadsResp); err != nil {
		c.logger.Errorf("Failed to decode threads response: %v", err)
		return nil, fmt.Errorf("failed to decode threads response: %w", err)
	}

	c.logger.Infof("Successfully retrieved %d threads (total: %d)",
		len(threadsResp.Threads), threadsResp.Total)
	return &threadsResp, nil
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
