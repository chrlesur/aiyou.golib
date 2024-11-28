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
	"os"
	"time"

	"github.com/chrlesur/aiyou.golib"
)

const maxRetries = 3
const retryDelay = 5 * time.Second

func main() {
	// Créer un client avec plus de logs
	logger := aiyou.NewDefaultLogger(os.Stdout)
	logger.SetLevel(aiyou.DEBUG)

	client, err := aiyou.NewClient(
		"christophe.lesur@cloud-temple.com",
		"jitkin-Vobjav-jythi3",
		aiyou.WithLogger(logger),
	)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Créer une conversation à sauvegarder
	conversation := aiyou.SaveConversationRequest{
		Title:       "Test de conversation",
		AssistantID: "asst_VZAhLX1aPVnVQPXCtvsdAgg4",
		Messages: []aiyou.Message{
			{
				Role: "user",
				Content: []aiyou.ContentPart{
					{Type: "text", Text: "Bonjour, pouvez-vous m'aider ?"},
				},
			},
			{
				Role: "assistant",
				Content: []aiyou.ContentPart{
					{Type: "text", Text: "Bien sûr, comment puis-je vous aider ?"},
				},
			},
		},
	}

	// Sauvegarder la conversation
	fmt.Println("Sauvegarde de la conversation...")
	resp, err := client.SaveConversation(context.Background(), conversation)
	if err != nil {
		log.Printf("Erreur lors de la sauvegarde: %v", err)
		return
	}

	fmt.Printf("\nConversation sauvegardée avec succès !\n")
	fmt.Printf("Thread ID: %s\n", resp.ID)
	fmt.Printf("Created at: %s\n", time.Unix(resp.CreatedAt, 0))
	fmt.Printf("Messages count: %d\n", len(resp.Messages))

	// Attente initiale
	fmt.Printf("\nAttente de 5 secondes pour la propagation...\n")
	time.Sleep(5 * time.Second)

	// Boucle de tentatives pour récupérer la conversation
	var thread *aiyou.ConversationThread
	for i := 0; i < maxRetries; i++ {
		fmt.Printf("\nTentative %d/%d de récupération de la conversation...\n", i+1, maxRetries)

		thread, err = client.GetConversation(context.Background(), resp.ID)
		if err == nil {
			break
		}

		if i < maxRetries-1 {
			fmt.Printf("Attente de %v avant nouvelle tentative...\n", retryDelay)
			time.Sleep(retryDelay)
		} else {
			log.Printf("Échec de la récupération après %d tentatives: %v", maxRetries, err)
			return
		}
	}
	fmt.Printf("\nConversation récupérée :\n")
	fmt.Printf("Thread ID: %s\n", thread.ID)
	fmt.Printf("Title: %s\n", thread.Title)
	fmt.Printf("Assistant ID: %s\n", thread.AssistantID)
	fmt.Printf("User ID: %s\n", thread.UserID)
	fmt.Printf("Created at: %s\n", thread.CreatedAt.Format(time.RFC3339))
	fmt.Printf("Last message at: %s\n", thread.LastMessageAt.Format(time.RFC3339))

	for i, msg := range thread.Messages {
		fmt.Printf("\nMessage %d:\n", i+1)
		fmt.Printf("Role: %s\n", msg.Role)
		for _, content := range msg.Content {
			fmt.Printf("Content: %s\n", content.Text)
		}
	}
}
