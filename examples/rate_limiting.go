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
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/chrlesur/aiyou.golib"
	"github.com/fatih/color"
)

var (
	email        string
	password     string
	baseURL      string
	debug        bool
	quietMode    bool
	requestCount int
	rateLimit    float64
	burstSize    int

	// Couleurs pour l'interface
	successColor = color.New(color.FgGreen)
	errorColor   = color.New(color.FgRed)
	infoColor    = color.New(color.FgYellow)
)

func init() {
	flag.StringVar(&email, "email", "", "Email pour l'authentification (obligatoire)")
	flag.StringVar(&password, "password", "", "Mot de passe pour l'authentification (obligatoire)")
	flag.StringVar(&baseURL, "url", "https://ai.dragonflygroup.fr", "URL de base de l'API")
	flag.BoolVar(&debug, "debug", false, "Active les logs de debug")
	flag.BoolVar(&quietMode, "quiet", false, "Désactive les messages de statut")
	flag.IntVar(&requestCount, "requests", 10, "Nombre de requêtes à effectuer")
	flag.Float64Var(&rateLimit, "rate", 2.0, "Nombre de requêtes par seconde")
	flag.IntVar(&burstSize, "burst", 3, "Taille du burst initial")
}

func createClient() (*aiyou.Client, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("email et mot de passe requis")
	}

	logger := aiyou.NewDefaultLogger(os.Stderr)
	if debug {
		logger.SetLevel(aiyou.DEBUG)
	} else if quietMode {
		logger.SetLevel(aiyou.ERROR)
	} else {
		logger.SetLevel(aiyou.INFO)
	}

	return aiyou.NewClient(
		email,
		password,
		aiyou.WithLogger(logger),
		aiyou.WithBaseURL(baseURL),
		aiyou.WithRateLimiter(aiyou.RateLimiterConfig{
			RequestsPerSecond: rateLimit,
			BurstSize:         burstSize,
			WaitTimeout:       5 * time.Second,
		}),
	)
}

func makeRequest(client *aiyou.Client, reqNum int, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	ctx := context.Background()
	start := time.Now()

	assistants, err := client.GetUserAssistants(ctx)
	duration := time.Since(start)

	if err != nil {
		switch e := err.(type) {
		case *aiyou.RateLimitError:
			if !quietMode {
				errorColor.Fprintf(os.Stderr, "Requête %d: Limite de taux atteinte, réessayer dans %d secondes\n",
					reqNum, e.RetryAfter)
			}
		default:
			if !quietMode {
				errorColor.Fprintf(os.Stderr, "Requête %d: Erreur: %v\n", reqNum, err)
			}
		}
		results <- fmt.Sprintf("Requête %d: Échec - %v", reqNum, err)
		return
	}

	if !quietMode {
		successColor.Fprintf(os.Stderr, "Requête %d: Succès (durée: %v) - %d assistants récupérés\n",
			reqNum, duration, len(assistants.Members))
	}
	results <- fmt.Sprintf("Requête %d: Succès - %v", reqNum, duration)
}

func main() {
	flag.Parse()

	if email == "" || password == "" {
		fmt.Fprintln(os.Stderr, "Les paramètres email et password sont obligatoires")
		flag.Usage()
		os.Exit(1)
	}

	client, err := createClient()
	if err != nil {
		errorColor.Fprintf(os.Stderr, "Erreur lors de la création du client: %v\n", err)
		os.Exit(1)
	}

	if !quietMode {
		infoColor.Fprintf(os.Stderr, "\nDémonstration du rate limiting\n")
		infoColor.Fprintf(os.Stderr, "Requêtes par seconde: %.1f\n", rateLimit)
		infoColor.Fprintf(os.Stderr, "Taille du burst: %d\n", burstSize)
		infoColor.Fprintf(os.Stderr, "Nombre total de requêtes: %d\n\n", requestCount)
	}

	start := time.Now()
	var wg sync.WaitGroup
	results := make(chan string, requestCount)

	// Lancement des requêtes concurrentes
	for i := 0; i < requestCount; i++ {
		wg.Add(1)
		go makeRequest(client, i+1, &wg, results)
		time.Sleep(50 * time.Millisecond) // Petit délai entre les lancements
	}

	// Attendre que toutes les requêtes soient terminées
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collecter les résultats
	var successes, failures int
	for result := range results {
		if !quietMode {
			if strings.Contains(result, "Succès") {
				successes++
			} else {
				failures++
			}
		}
	}

	if !quietMode {
		duration := time.Since(start)
		fmt.Fprintln(os.Stderr, "\nRésumé:")
		fmt.Fprintf(os.Stderr, "Durée totale: %v\n", duration)
		fmt.Fprintf(os.Stderr, "Requêtes réussies: %d\n", successes)
		fmt.Fprintf(os.Stderr, "Requêtes échouées: %d\n", failures)
		fmt.Fprintf(os.Stderr, "Taux effectif: %.2f req/s\n", float64(successes)/duration.Seconds())
	}
}
