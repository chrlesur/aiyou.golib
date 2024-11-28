// File: examples/threads.go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chrlesur/aiyou.golib"
)

func main() {
	client, err := aiyou.NewClient("your-email@example.com", "your-password")
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Définir les filtres de recherche
	filter := &aiyou.ThreadFilter{
		AssistantID: "asst_123",
		StartDate:   time.Now().AddDate(0, -1, 0), // Dernier mois
		Page:        1,
		Limit:       10,
	}

	// Récupérer les threads
	threads, err := client.GetUserThreads(context.Background(), filter)
	if err != nil {
		log.Fatalf("Error getting threads: %v", err)
	}

	// Afficher les résultats
	fmt.Printf("Total threads: %d\n", threads.Total)
	fmt.Printf("Page: %d/%d\n\n", threads.Page, (threads.Total+threads.Limit-1)/threads.Limit)

	for _, thread := range threads.Threads {
		fmt.Printf("Thread: %s\n", thread.Title)
		fmt.Printf("ID: %s\n", thread.ID)
		fmt.Printf("Assistant: %s\n", thread.AssistantID)
		fmt.Printf("Created: %s\n", thread.CreatedAt.Format(time.RFC3339))
		fmt.Printf("Last updated: %s\n", thread.UpdatedAt.Format(time.RFC3339))
		fmt.Printf("Messages: %d\n\n", len(thread.Messages))
	}

	// Exemple de suppression d'un thread
	if len(threads.Threads) > 0 {
		threadToDelete := threads.Threads[0].ID
		if err := client.DeleteThread(context.Background(), threadToDelete); err != nil {
			log.Printf("Error deleting thread %s: %v", threadToDelete, err)
		} else {
			fmt.Printf("Successfully deleted thread: %s\n", threadToDelete)
		}
	}
}
