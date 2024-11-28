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
	"io/ioutil"
	"net/http"
)

// SaveConversation sauvegarde une conversation dans le système
func (c *Client) SaveConversation(ctx context.Context, req SaveConversationRequest) (*SaveConversationResponse, error) {
	endpoint := "/api/v1/save"
	c.logger.Debugf("Saving conversation with assistant ID: %s", req.AssistantID)

	// Validation basique
	if req.AssistantID == "" {
		return nil, fmt.Errorf("assistantId is required")
	}
	if req.Conversation == "" {
		return nil, fmt.Errorf("conversation is required")
	}

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

	// Accepter à la fois 200 et 201 comme codes de succès
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var saveResp SaveConversationResponse
	if err := json.NewDecoder(resp.Body).Decode(&saveResp); err != nil {
		c.logger.Errorf("Failed to decode save conversation response: %v", err)
		return nil, fmt.Errorf("failed to decode save conversation response: %w", err)
	}

	c.logger.Infof("Successfully saved conversation with thread ID: %s", saveResp.ID)
	return &saveResp, nil
}

// GetConversation récupère une conversation spécifique par son ID
func (c *Client) GetConversation(ctx context.Context, threadID string) (*ConversationThread, error) {
	endpoint := fmt.Sprintf("/api/v1/user/threads")
	c.logger.Debugf("Fetching conversation thread: %s", threadID)

	resp, err := c.AuthenticatedRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		c.logger.Errorf("Failed to fetch conversation: %v", err)
		return nil, fmt.Errorf("failed to fetch conversation: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.logger.Errorf("Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	c.logger.Debugf("Raw response: %s", string(body))

	var threadsOutput UserThreadsOutput
	if err := json.Unmarshal(body, &threadsOutput); err != nil {
		c.logger.Errorf("Failed to decode threads response: %v", err)
		return nil, fmt.Errorf("failed to decode threads response: %w", err)
	}

	// Chercher le thread spécifique dans la liste
	for _, thread := range threadsOutput.Threads {
		if thread.ID == threadID {
			c.logger.Infof("Successfully found conversation thread: %s", thread.ID)
			return &thread, nil
		}
	}

	return nil, fmt.Errorf("thread not found: %s", threadID)
}
