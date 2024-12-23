package aiyou

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestSaveConversation(t *testing.T) {
	// Configuration du logger pour voir ce qui se passe
	logger := NewDefaultLogger(os.Stderr)
	logger.SetLevel(DEBUG)

	// Création du client avec les credentials de test
	client, err := NewClient(
		WithEmailPassword(testConfig.Email, testConfig.Password),
		WithBaseURL(testConfig.BaseURL),
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Créer une requête de test pour sauvegarder une conversation
	req := SaveConversationRequest{
		AssistantID:  "287", // ID d'un assistant existant dans le système
		Conversation: "Test conversation created at " + time.Now().Format(time.RFC3339),
		FirstMessage: "Hello from test!",
		ContentJson:  `{"messages":[{"role":"user","content":[{"type":"text","text":"Hello from test"}]}]}`,
		ModelName:    "gpt-4",
	}

	t.Logf("Attempting to save conversation with assistant ID: %s", req.AssistantID)

	resp, err := client.SaveConversation(context.Background(), req)
	if err != nil {
		t.Fatalf("SaveConversation failed: %v", err)
	}

	// Vérification et log des résultats
	if resp == nil {
		t.Fatal("Expected non-nil response")
	}

	t.Logf("Successfully saved conversation:")
	t.Logf("  Thread ID: %s", resp.ID)
	t.Logf("  Object: %s", resp.Object)
	t.Logf("  Created At: %d", resp.CreatedAt)

	// Test de récupération de la conversation
	thread, err := client.GetConversation(context.Background(), resp.ID)
	if err != nil {
		t.Fatalf("GetConversation failed: %v", err)
	}

	t.Logf("Retrieved conversation:")
	t.Logf("  Thread ID: %s", thread.ID)
	t.Logf("  Content: %s", thread.Content)
	t.Logf("  Assistant Name: %s", thread.AssistantName)
	t.Logf("  First Message: %s", thread.FirstMessage)
	t.Logf("  Created At: %v", thread.CreatedAt)
}

func TestSaveConversationUnauthorized(t *testing.T) {
	logger := NewDefaultLogger(os.Stderr)
	logger.SetLevel(DEBUG)

	// Test avec des credentials invalides
	client, err := NewClient(
		WithEmailPassword("invalid@example.com", "wrong_password"),
		WithBaseURL(testConfig.BaseURL),
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req := SaveConversationRequest{
		AssistantID:  "287",
		Conversation: "Test unauthorized conversation",
		FirstMessage: "Hello!",
		ContentJson:  "{}",
		ModelName:    "gpt-4",
	}

	_, err = client.SaveConversation(context.Background(), req)
	if err == nil {
		t.Error("Expected error for unauthorized request")
	}

	t.Logf("Got expected error: %v", err)
}

func TestGetConversation(t *testing.T) {
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

	// D'abord créer une nouvelle conversation pour le test
	saveReq := SaveConversationRequest{
		AssistantID:  "287",
		Conversation: "Test conversation for retrieval " + time.Now().Format(time.RFC3339),
		FirstMessage: "Hello from retrieval test!",
		ContentJson:  `{"messages":[{"role":"user","content":[{"type":"text","text":"Hello from retrieval test"}]}]}`,
		ModelName:    "gpt-4",
	}

	savedResp, err := client.SaveConversation(context.Background(), saveReq)
	if err != nil {
		t.Fatalf("Failed to create test conversation: %v", err)
	}

	t.Logf("Created test conversation with ID: %s", savedResp.ID)

	// Maintenant récupérer la conversation
	thread, err := client.GetConversation(context.Background(), savedResp.ID)
	if err != nil {
		t.Fatalf("GetConversation failed: %v", err)
	}

	t.Logf("Retrieved conversation details:")
	t.Logf("  Thread ID: %s", thread.ID)
	t.Logf("  Content: %s", thread.Content)
	t.Logf("  Assistant Name: %s", thread.AssistantName)
	t.Logf("  Model: %s", *thread.AssistantModel)
	t.Logf("  First Message: %s", thread.FirstMessage)
	t.Logf("  Created At: %v", thread.CreatedAt)

	// Vérifier que les données correspondent
	if thread.FirstMessage != saveReq.FirstMessage {
		t.Errorf("First message mismatch. Expected: %s, Got: %s", saveReq.FirstMessage, thread.FirstMessage)
	}
}

func TestGetNonExistentConversation(t *testing.T) {
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

	_, err = client.GetConversation(context.Background(), "non_existent_thread_id")
	if err == nil {
		t.Error("Expected error when getting non-existent conversation")
	}

	t.Logf("Got expected error: %v", err)
}
