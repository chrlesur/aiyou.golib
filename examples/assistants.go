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

// examples/assistants.go
package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/chrlesur/aiyou.golib"
)

// truncateText prend un texte et le tronque à environ 30 mots
// en s'assurant de ne pas couper au milieu d'un mot
func truncateText(text string, wordLimit int) string {
	words := strings.Fields(text)
	if len(words) <= wordLimit {
		return text
	}

	// Prendre les premiers mots jusqu'à la limite
	truncated := strings.Join(words[:wordLimit], " ")

	// Ajouter des points de suspension
	return truncated + "..."
}

func main() {
	client, err := aiyou.NewClient("your-email@example.com", "your-password")
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Récupérer les assistants
	assistants, err := client.GetUserAssistants(context.Background())
	if err != nil {
		log.Fatalf("Error getting assistants: %v", err)
	}

	// Afficher les assistants
	fmt.Printf("Total assistants: %d\n", assistants.TotalItems)
	for _, assistant := range assistants.Members {
		fmt.Printf("\nAssistant: %s (%s)\n", assistant.Name, assistant.AssistantID)
		if len(assistant.ThreadHistories) > 0 {
			summary := truncateText(assistant.ThreadHistories[0].FirstMessage, 30)
			fmt.Printf("Last conversation: %s\n", summary)
		}
	}
}
