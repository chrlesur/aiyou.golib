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
// File: pkg/aiyou/conversation.go
package aiyou

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// SaveConversation sauvegarde une conversation dans le système
func (c *Client) SaveConversation(ctx context.Context, req SaveConversationRequest) (*SaveConversationResponse, error) {
	endpoint := "/api/v1/save"
	c.logger.Debugf("Saving conversation with title: %s", req.Title)

	jsonData, err := json.Marshal(req)
	if err != nil {
		c.logger.Errorf("Failed to marshal save conversation request: %v", err)
		return nil, fmt.Errorf("failed to marshal save conversation request: %w", err)
	}

	resp, err := c.AuthenticatedRequest(ctx, "POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Errorf("Failed to save conversation: %v", err)
		return nil, fmt.Errorf("failed to save conversation: %w", err)
	}
	defer resp.Body.Close()

	var saveResp SaveConversationResponse
	if err := json.NewDecoder(resp.Body).Decode(&saveResp); err != nil {
		c.logger.Errorf("Failed to decode save conversation response: %v", err)
		return nil, fmt.Errorf("failed to decode save conversation response: %w", err)
	}

	c.logger.Infof("Successfully saved conversation with thread ID: %s", saveResp.Thread.ID)
	return &saveResp, nil
}

// GetConversation récupère une conversation spécifique par son ID
func (c *Client) GetConversation(ctx context.Context, threadID string) (*ConversationThread, error) {
	endpoint := fmt.Sprintf("/api/v1/threads/%s", threadID)
	c.logger.Debugf("Fetching conversation thread: %s", threadID)

	resp, err := c.AuthenticatedRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		c.logger.Errorf("Failed to fetch conversation: %v", err)
		return nil, fmt.Errorf("failed to fetch conversation: %w", err)
	}
	defer resp.Body.Close()

	var thread ConversationThread
	if err := json.NewDecoder(resp.Body).Decode(&thread); err != nil {
		c.logger.Errorf("Failed to decode conversation thread: %v", err)
		return nil, fmt.Errorf("failed to decode conversation thread: %w", err)
	}

	c.logger.Infof("Successfully retrieved conversation thread: %s", thread.ID)
	return &thread, nil
}
