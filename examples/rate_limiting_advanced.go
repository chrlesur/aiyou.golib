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
	"sync"
	"sync/atomic"
	"time"

	"github.com/chrlesur/aiyou.golib"
	"github.com/fatih/color"
)

var (
	email        string
	password     string
	assistantID  string
	baseURL      string
	debug        bool
	quietMode    bool
	requestCount int
	rateLimit    float64
	burstSize    int
	requestTypes string

	// Statistiques atomiques
	successCount   uint64
	failureCount   uint64
	rateLimitCount uint64

	// Couleurs pour l'interface
	successColor = color.New(color.FgGreen)
	errorColor   = color.New(color.FgRed)
	infoColor    = color.New(color.FgYellow)
	headerColor  = color.New(color.FgCyan, color.Bold)
)

type RequestStats struct {
	RequestType string
	Duration    time.Duration
	Error       error
}

func init() {
	flag.StringVar(&email, "email", "", "Email pour l'authentification (obligatoire)")
	flag.StringVar(&password, "password", "", "Mot de passe pour l'authentification (obligatoire)")
	flag.StringVar(&assistantID, "assistant", "", "ID de l'assistant pour les tests de chat")
	flag.StringVar(&baseURL, "url", "https://ai.dragonflygroup.fr", "URL de base de l'API")
	flag.BoolVar(&debug, "debug", false, "Active les logs de debug")
	flag.BoolVar(&quietMode, "quiet", false, "Désactive les messages de statut")
	flag.IntVar(&requestCount, "requests", 10, "Nombre de requêtes à effectuer")
	flag.Float64Var(&rateLimit, "rate", 2.0, "Nombre de requêtes par seconde")
	flag.IntVar(&burstSize, "burst", 3, "Taille du burst initial")
	flag.StringVar(&requestTypes, "types", "assistants", "Types de requêtes à tester (assistants, chat)")
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

func makeGenericRequest(ctx context.Context, client *aiyou.Client, reqType string, reqNum int) RequestStats {
	start := time.Now()
	var err error

	switch reqType {
	case "chat":
		if assistantID == "" {
			return RequestStats{
				RequestType: reqType,
				Duration:    0,
				Error:       fmt.Errorf("assistant ID required for chat requests"),
			}
		}
		msg := aiyou.NewTextMessage("user", "Test message for rate limiting")
		_, err = client.CreateChatCompletion(ctx, []aiyou.Message{msg}, assistantID)

	case "assistants":
		response, err := client.GetUserAssistants(ctx)
		if err == nil && !quietMode {
			fmt.Fprintf(os.Stderr, "Récupéré %d assistants\n", len(response.Members))
		}

	default:
		return RequestStats{
			RequestType: reqType,
			Duration:    0,
			Error:       fmt.Errorf("type de requête inconnu: %s", reqType),
		}
	}

	stats := RequestStats{
		RequestType: reqType,
		Duration:    time.Since(start),
		Error:       err,
	}

	if err != nil {
		atomic.AddUint64(&failureCount, 1)
		if _, isRateLimit := err.(*aiyou.RateLimitError); isRateLimit {
			atomic.AddUint64(&rateLimitCount, 1)
		}
	} else {
		atomic.AddUint64(&successCount, 1)
	}

	return stats
}

func processRequestStats(stats RequestStats, reqNum int) {
	if quietMode {
		return
	}

	if stats.Error != nil {
		if _, isRateLimit := stats.Error.(*aiyou.RateLimitError); isRateLimit {
			errorColor.Fprintf(os.Stderr, "Requête %d (%s): Rate limit atteint\n",
				reqNum, stats.RequestType)
		} else {
			errorColor.Fprintf(os.Stderr, "Requête %d (%s): %v\n",
				reqNum, stats.RequestType, stats.Error)
		}
	} else {
		successColor.Fprintf(os.Stderr, "Requête %d (%s): OK (%v)\n",
			reqNum, stats.RequestType, stats.Duration.Round(time.Millisecond))
	}
}

func printSummary(duration time.Duration) {
	if quietMode {
		return
	}

	headerColor.Fprintln(os.Stderr, "\nRésumé des tests:")
	fmt.Fprintf(os.Stderr, "Durée totale: %v\n", duration)
	fmt.Fprintf(os.Stderr, "Configuration:\n")
	fmt.Fprintf(os.Stderr, "- Rate limit: %.1f req/s\n", rateLimit)
	fmt.Fprintf(os.Stderr, "- Burst size: %d\n", burstSize)
	fmt.Fprintf(os.Stderr, "- Requêtes totales: %d\n", requestCount)

	fmt.Fprintf(os.Stderr, "\nRésultats:\n")
	successes := atomic.LoadUint64(&successCount)
	failures := atomic.LoadUint64(&failureCount)
	rateLimits := atomic.LoadUint64(&rateLimitCount)

	successColor.Fprintf(os.Stderr, "✓ Requêtes réussies: %d\n", successes)
	errorColor.Fprintf(os.Stderr, "✗ Requêtes échouées: %d\n", failures)
	fmt.Fprintf(os.Stderr, " - Erreurs de rate limit: %d\n", rateLimits)
	fmt.Fprintf(os.Stderr, " - Autres erreurs: %d\n", failures-rateLimits)

	if successes > 0 {
		fmt.Fprintf(os.Stderr, "\nPerformances:\n")
		fmt.Fprintf(os.Stderr, "- Taux effectif: %.2f req/s\n", float64(successes)/duration.Seconds())
		fmt.Fprintf(os.Stderr, "- Taux de succès: %.1f%%\n", float64(successes)*100/float64(successes+failures))
	}
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
		headerColor.Fprintf(os.Stderr, "\nTest de rate limiting\n")
		fmt.Fprintf(os.Stderr, "Type de requêtes: %s\n", requestTypes)
	}

	var wg sync.WaitGroup
	results := make(chan RequestStats, requestCount)
	start := time.Now()

	// Calcul du délai entre les requêtes
	delay := time.Second / time.Duration(rateLimit)
	if delay < 50*time.Millisecond {
		delay = 50 * time.Millisecond
	}

	// Lancement des requêtes
	for i := 0; i < requestCount; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			stats := makeGenericRequest(context.Background(), client, requestTypes, num)
			processRequestStats(stats, num)
			results <- stats
		}(i + 1)

		time.Sleep(delay)
	}

	// Attendre la fin des requêtes
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collecter les résultats
	for range results {
		// Les statistiques sont déjà traitées dans processRequestStats
	}

	printSummary(time.Since(start))
}
