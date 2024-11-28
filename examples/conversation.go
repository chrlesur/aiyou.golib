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

// File: examples/conversation.go
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

	// Créer une conversation à sauvegarder
	conversation := aiyou.SaveConversationRequest{
		Title:       "Discussion sur l'IA",
		AssistantID: "asst_123",
		Messages: []aiyou.Message{
			{
				Role: "user",
				Content: []aiyou.ContentPart{
					{Type: "text", Text: "Qu'est-ce que l'intelligence artificielle?"},
				},
			},
			{
				Role: "assistant",
				Content: []aiyou.ContentPart{
					{Type: "text", Text: "L'intelligence artificielle est..."},
				},
			},
		},
	}

	// Sauvegarder la conversation
	resp, err := client.SaveConversation(context.Background(), conversation)
	if err != nil {
		log.Fatalf("Error saving conversation: %v", err)
	}

	fmt.Printf("Conversation saved successfully!\n")
	fmt.Printf("Thread ID: %s\n", resp.Thread.ID)
	fmt.Printf("Title: %s\n", resp.Thread.Title)
	fmt.Printf("Created at: %s\n", resp.Thread.CreatedAt)

	// Récupérer la conversation sauvegardée
	thread, err := client.GetConversation(context.Background(), resp.Thread.ID)
	if err != nil {
		log.Fatalf("Error retrieving conversation: %v", err)
	}

	fmt.Printf("\nRetrieved conversation:\n")
	fmt.Printf("Title: %s\n", thread.Title)
	fmt.Printf("Number of messages: %d\n", len(thread.Messages))
}
