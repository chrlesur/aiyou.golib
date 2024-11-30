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

// assistants.go
package aiyou

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetUserAssistants récupère la liste des assistants disponibles pour l'utilisateur
func (c *Client) GetUserAssistants(ctx context.Context) (*AssistantsResponse, error) {
	endpoint := "/api/v1/user/assistants"
	c.logger.Debugf("Fetching user assistants from %s", endpoint)

	resp, err := c.AuthenticatedRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		c.logger.Errorf("Failed to fetch assistants: %v", err)
		return nil, fmt.Errorf("failed to fetch assistants: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var assistantsResp AssistantsResponse
	if err := json.NewDecoder(resp.Body).Decode(&assistantsResp); err != nil {
		c.logger.Errorf("Failed to decode assistants response: %v", err)
		return nil, fmt.Errorf("failed to decode assistants response: %w", err)
	}

	c.logger.Infof("Successfully retrieved %d assistants", len(assistantsResp.Members))
	return &assistantsResp, nil
}
