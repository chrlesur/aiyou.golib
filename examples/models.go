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

// File: examples/models.go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/chrlesur/aiyou.golib"
)

func main() {
	client, err := aiyou.NewClient("your-email@example.com", "your-password")
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Créer un nouveau modèle
	modelReq := aiyou.ModelRequest{
		Name:        "Custom GPT Model",
		Description: "A custom GPT model for specific tasks",
		Properties: aiyou.ModelProperties{
			MaxTokens:    4096,
			Temperature:  0.7,
			Provider:     "OpenAI",
			Capabilities: []string{"chat", "completion"},
		},
	}

	model, err := client.CreateModel(context.Background(), modelReq)
	if err != nil {
		log.Fatalf("Error creating model: %v", err)
	}
	fmt.Printf("Created model: %s (ID: %s)\n", model.Model.Name, model.Model.ID)

	// Lister tous les modèles
	models, err := client.GetModels(context.Background())
	if err != nil {
		log.Fatalf("Error getting models: %v", err)
	}

	fmt.Printf("\nTotal models: %d\n", models.Total)
	for _, m := range models.Models {
		fmt.Printf("\nModel: %s (%s)\n", m.Name, m.ID)
		fmt.Printf("Description: %s\n", m.Description)
		fmt.Printf("Version: %s\n", m.Version)
		fmt.Printf("Capabilities: %v\n", m.Properties.Capabilities)
	}
}
