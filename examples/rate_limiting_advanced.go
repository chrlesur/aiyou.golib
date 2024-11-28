// File: examples/rate_limiting_advanced.go

/*
Copyright (C) 2024 Cloud Temple

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.
*/

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/chrlesur/aiyou.golib"
)

func main() {
	// Configuration du logger
	logger := aiyou.NewDefaultLogger(os.Stdout)
	logger.SetLevel(aiyou.DEBUG)

	// Configuration du client avec rate limiting
	client, err := aiyou.NewClient(
		"your-email@example.com",
		"your-password",
		aiyou.WithLogger(logger),
		aiyou.WithRateLimiter(aiyou.RateLimiterConfig{
			RequestsPerSecond: 2, // Limite à 2 requêtes par seconde
			BurstSize:         3, // Permet un burst initial de 3 requêtes
			WaitTimeout:       5 * time.Second,
		}),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Démonstration de requêtes concurrentes avec rate limiting
	demonstrateConcurrentRequests(client)
}

func demonstrateConcurrentRequests(client *aiyou.Client) {
	fmt.Println("Starting concurrent requests demonstration...")

	var wg sync.WaitGroup
	requestCount := 10
	results := make(chan string, requestCount)

	// Lancer plusieurs requêtes concurrentes
	for i := 0; i < requestCount; i++ {
		wg.Add(1)
		go func(reqNum int) {
			defer wg.Done()

			ctx := context.Background()
			start := time.Now()

			// Simuler différents types de requêtes
			var result string
			switch reqNum % 3 {
			case 0:
				// Requête de chat
				msg := aiyou.NewTextMessage("user", "Hello, how are you?")
				resp, err := client.CreateChatCompletion(ctx, []aiyou.Message{msg}, "your-assistant-id")
				if err != nil {
					handleError(err, reqNum, results)
					return
				}
				result = fmt.Sprintf("Chat request %d completed in %v", reqNum, time.Since(start))

			case 1:
				// Requête d'assistants
				assistants, err := client.GetUserAssistants(ctx)
				if err != nil {
					handleError(err, reqNum, results)
					return
				}
				result = fmt.Sprintf("Assistants request %d completed in %v, found %d assistants",
					reqNum, time.Since(start), len(assistants.Assistants))

			case 2:
				// Requête de modèles
				models, err := client.GetModels(ctx)
				if err != nil {
					handleError(err, reqNum, results)
					return
				}
				result = fmt.Sprintf("Models request %d completed in %v, found %d models",
					reqNum, time.Since(start), len(models.Models))
			}

			results <- result
		}(i)
	}

	// Goroutine pour collecter et afficher les résultats
	go func() {
		for result := range results {
			fmt.Println(result)
		}
	}()

	// Attendre que toutes les requêtes soient terminées
	wg.Wait()
	close(results)

	fmt.Println("All requests completed!")
}

func handleError(err error, reqNum int, results chan<- string) {
	switch e := err.(type) {
	case *aiyou.RateLimitError:
		results <- fmt.Sprintf("Request %d: %v", reqNum, e)
	case *aiyou.NetworkError:
		results <- fmt.Sprintf("Request %d: Network error: %v", reqNum, e)
	case *aiyou.AuthenticationError:
		results <- fmt.Sprintf("Request %d: Authentication error: %v", reqNum, e)
	default:
		results <- fmt.Sprintf("Request %d: Unexpected error: %v", reqNum, err)
	}
}
