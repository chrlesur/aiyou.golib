// File: examples/rate_limiting.go

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

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chrlesur/aiyou.golib"
)

func main() {
	// Créer un client avec rate limiting configuré
	client, err := aiyou.NewClient(
		"your-email@example.com",
		"your-password",
		aiyou.WithRateLimiter(aiyou.RateLimiterConfig{
			RequestsPerSecond: 2,               // Limite à 2 requêtes par seconde
			BurstSize:         3,               // Permet un burst initial de 3 requêtes
			WaitTimeout:       time.Second * 5, // Timeout d'attente de 5 secondes
		}),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Démonstration du rate limiting
	fmt.Println("Starting rate limiting demonstration...")

	// Test du burst initial
	fmt.Println("\nTesting burst capacity (3 requests):")
	for i := 0; i < 3; i++ {
		if err := makeRequest(client, i); err != nil {
			log.Printf("Request %d failed: %v", i, err)
		}
	}

	// Test du taux limité
	fmt.Println("\nTesting rate limiting (2 requests per second):")
	for i := 3; i < 8; i++ {
		if err := makeRequest(client, i); err != nil {
			log.Printf("Request %d failed: %v", i, err)
		}
		time.Sleep(time.Millisecond * 100) // Petit délai pour la démonstration
	}
}

func makeRequest(client *aiyou.Client, i int) error {
	ctx := context.Background()

	// Exemple d'utilisation avec GetUserAssistants
	start := time.Now()
	assistants, err := client.GetUserAssistants(ctx)
	duration := time.Since(start)

	if err != nil {
		switch e := err.(type) {
		case *aiyou.RateLimitError:
			fmt.Printf("Request %d: Rate limit exceeded, retry after %d seconds\n",
				i, e.RetryAfter)
		default:
			fmt.Printf("Request %d: Error: %v\n", i, err)
		}
		return err
	}

	fmt.Printf("Request %d: Success (took %v) - Retrieved %d assistants\n",
		i, duration, len(assistants.Assistants))
	return nil
}
