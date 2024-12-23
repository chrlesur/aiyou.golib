package aiyou

import (
	"context"
	"os"
	"testing"
)

func TestGetUserAssistants(t *testing.T) {
	// Créer un logger pour les tests
	logger := NewDefaultLogger(os.Stderr)
	logger.SetLevel(DEBUG) // Active les logs pour voir ce qui se passe

	// Utilisation des credentials depuis test_config.go
	client, err := NewClient(
		WithEmailPassword(testConfig.Email, testConfig.Password),
		WithBaseURL(testConfig.BaseURL),
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	resp, err := client.GetUserAssistants(context.Background())
	if err != nil {
		t.Fatalf("GetUserAssistants failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected non-nil response")
	}

	// Log des résultats pour le debug
	t.Logf("Found %d assistants", len(resp.Members))
	for i, assistant := range resp.Members {
		t.Logf("Assistant %d: ID=%s, Name=%s", i+1, assistant.ID, assistant.Name)
	}
}

func TestGetUserAssistantsUnauthorized(t *testing.T) {
	logger := NewDefaultLogger(os.Stderr)
	logger.SetLevel(DEBUG)

	// Test avec des credentials invalides
	client, err := NewClient(
		WithEmailPassword("invalid@example.com", "wrong_password"),
		WithBaseURL(testConfig.BaseURL), // On garde la bonne URL
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.GetUserAssistants(context.Background())
	if err == nil {
		t.Error("Expected error for unauthorized request")
	}

	t.Logf("Got expected error: %v", err)
}
