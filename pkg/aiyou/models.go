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

// File: pkg/aiyou/models.go
package aiyou

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// CreateModel crée un nouveau modèle dans le système AI.YOU
func (c *Client) CreateModel(ctx context.Context, req ModelRequest) (*ModelResponse, error) {
	endpoint := "/api/v1/models"
	c.logger.Debugf("Creating new model with name: %s", req.Name)

	jsonData, err := json.Marshal(req)
	if err != nil {
		c.logger.Errorf("Failed to marshal model request: %v", err)
		return nil, fmt.Errorf("failed to marshal model request: %w", err)
	}

	resp, err := c.AuthenticatedRequest(ctx, "POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Errorf("Failed to create model: %v", err)
		return nil, fmt.Errorf("failed to create model: %w", err)
	}
	defer resp.Body.Close()

	var modelResp ModelResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelResp); err != nil {
		c.logger.Errorf("Failed to decode model response: %v", err)
		return nil, fmt.Errorf("failed to decode model response: %w", err)
	}

	c.logger.Infof("Successfully created model with ID: %s", modelResp.Model.ID)
	return &modelResp, nil
}

// GetModels récupère la liste des modèles disponibles
func (c *Client) GetModels(ctx context.Context) (*ModelsResponse, error) {
	endpoint := "/api/v1/models"
	c.logger.Debugf("Fetching models from %s", endpoint)

	resp, err := c.AuthenticatedRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		c.logger.Errorf("Failed to fetch models: %v", err)
		return nil, fmt.Errorf("failed to fetch models: %w", err)
	}
	defer resp.Body.Close()

	var modelsResp ModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		c.logger.Errorf("Failed to decode models response: %v", err)
		return nil, fmt.Errorf("failed to decode models response: %w", err)
	}

	c.logger.Infof("Successfully retrieved %d models", len(modelsResp.Models))
	return &modelsResp, nil
}
