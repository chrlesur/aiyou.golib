package aiyou

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func TestClientAuthentication(t *testing.T) {
	logger := NewDefaultLogger(os.Stderr)
	logger.SetLevel(DEBUG)

	testCases := []struct {
		name      string
		email     string
		password  string
		wantError bool
	}{
		{
			name:      "Valid Credentials",
			email:     testConfig.Email,
			password:  testConfig.Password,
			wantError: false,
		},
		{
			name:      "Invalid Credentials",
			email:     "invalid@example.com",
			password:  "wrong_password",
			wantError: true,
		},
		{
			name:      "Empty Credentials",
			email:     "",
			password:  "",
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, err := NewClient(
				WithEmailPassword(tc.email, tc.password),
				WithBaseURL(testConfig.BaseURL),
				WithLogger(logger),
			)

			if err != nil && !tc.wantError {
				t.Fatalf("NewClient failed: %v", err)
			}

			if err == nil && tc.wantError {
				t.Error("Expected error but got none")
				return
			}

			if client != nil {
				// Tester une requête simple pour vérifier l'authentification
				_, err := client.GetUserAssistants(context.Background())
				if err != nil && !tc.wantError {
					t.Errorf("GetUserAssistants failed: %v", err)
				}
				if err == nil && tc.wantError {
					t.Error("Expected error but got none")
				}
			}
		})
	}
}

func TestClientRetry(t *testing.T) {
	logger := NewDefaultLogger(os.Stderr)
	logger.SetLevel(DEBUG)

	client, err := NewClient(
		WithEmailPassword(testConfig.Email, testConfig.Password),
		WithBaseURL(testConfig.BaseURL),
		WithLogger(logger),
		WithRetry(3, time.Second),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test avec une requête réelle
	resp, err := client.GetUserAssistants(context.Background())
	if err != nil {
		t.Errorf("GetUserAssistants failed: %v", err)
	}

	t.Logf("Successfully retrieved %d assistants with retry configuration", len(resp.Members))
}

func TestClientRateLimit(t *testing.T) {
	logger := NewDefaultLogger(os.Stderr)
	logger.SetLevel(DEBUG)

	client, err := NewClient(
		WithEmailPassword(testConfig.Email, testConfig.Password),
		WithBaseURL(testConfig.BaseURL),
		WithLogger(logger),
		WithRateLimiter(RateLimiterConfig{
			RequestsPerSecond: 2,
			BurstSize:         1,
			WaitTimeout:       time.Second,
		}),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Faire plusieurs requêtes rapides pour tester le rate limiting
	for i := 0; i < 5; i++ {
		t.Logf("Making request %d/5", i+1)
		_, err := client.GetUserAssistants(context.Background())
		if err != nil {
			if strings.Contains(err.Error(), "rate limit") {
				t.Logf("Got expected rate limit error: %v", err)
			} else {
				t.Errorf("Unexpected error: %v", err)
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func TestClientWithBearerToken(t *testing.T) {
	logger := NewDefaultLogger(os.Stderr)
	logger.SetLevel(DEBUG)

	// D'abord, obtenir un token valide via l'authentification normale
	client, err := NewClient(
		WithEmailPassword(testConfig.Email, testConfig.Password),
		WithBaseURL(testConfig.BaseURL),
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("Failed to create initial client: %v", err)
	}

	// Faire une requête pour s'assurer que nous avons un token valide
	_, err = client.GetUserAssistants(context.Background())
	if err != nil {
		t.Fatalf("Failed to authenticate: %v", err)
	}

	// Créer un nouveau client avec le token bearer
	tokenClient, err := NewClient(
		WithBearerToken("test_token"), // Note: normalement on utiliserait le vrai token
		WithBaseURL(testConfig.BaseURL),
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("Failed to create token client: %v", err)
	}

	// Tester le client avec token
	_, err = tokenClient.GetUserAssistants(context.Background())
	// On s'attend à une erreur car c'est un token de test
	if err == nil {
		t.Error("Expected error with test token")
	} else {
		t.Logf("Got expected error with test token: %v", err)
	}
}

func TestClientTimeouts(t *testing.T) {
	logger := NewDefaultLogger(os.Stderr)
	logger.SetLevel(DEBUG)

	client, err := NewClient(
		WithEmailPassword(testConfig.Email, testConfig.Password),
		WithBaseURL(testConfig.BaseURL),
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Créer un contexte avec timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err = client.GetUserAssistants(ctx)
	if err != nil {
		t.Logf("Got expected timeout or error: %v", err)
	}
}
