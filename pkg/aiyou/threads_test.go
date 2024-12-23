package aiyou

import (
	"context"
	"os"
	"testing"
)

func TestGetUserThreads(t *testing.T) {
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

	// Test de la récupération des threads avec pagination
	params := &UserThreadsParams{
		Page:         1,
		ItemsPerPage: 10,
	}

	resp, err := client.GetUserThreads(context.Background(), params)
	if err != nil {
		t.Fatalf("GetUserThreads failed: %v", err)
	}

	// Vérification et log des résultats
	if resp == nil {
		t.Fatal("Expected non-nil response")
	}

	t.Logf("Total threads: %d", resp.TotalItems)
	t.Logf("Current page: %d", resp.CurrentPage)
	t.Logf("Items per page: %d", resp.ItemsPerPage)
	t.Logf("Number of threads in response: %d", len(resp.Threads))

	// Log détaillé des threads reçus
	for i, thread := range resp.Threads {
		t.Logf("Thread %d:", i+1)
		t.Logf("  ID: %s", thread.ID)
		t.Logf("  Content: %s", thread.Content)
		t.Logf("  Assistant Name: %s", thread.AssistantName)
		t.Logf("  First Message: %s", thread.FirstMessage)
		t.Logf("  Created At: %v", thread.CreatedAt)
	}
}

func TestGetUserThreadsUnauthorized(t *testing.T) {
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

	params := &UserThreadsParams{
		Page:         1,
		ItemsPerPage: 10,
	}

	_, err = client.GetUserThreads(context.Background(), params)
	if err == nil {
		t.Error("Expected error for unauthorized request")
	}

	t.Logf("Got expected error: %v", err)
}

func TestDeleteThread(t *testing.T) {
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

	// D'abord, récupérer un thread existant
	params := &UserThreadsParams{
		Page:         1,
		ItemsPerPage: 1, // On ne prend que le premier thread
	}

	threads, err := client.GetUserThreads(context.Background(), params)
	if err != nil {
		t.Fatalf("Failed to get threads: %v", err)
	}

	if len(threads.Threads) == 0 {
		t.Skip("No threads available for deletion test")
		return
	}

	// Tenter de supprimer le premier thread
	threadID := threads.Threads[0].ID
	t.Logf("Attempting to delete thread: %s", threadID)

	err = client.DeleteThread(context.Background(), threadID)
	if err != nil {
		t.Errorf("DeleteThread failed: %v", err)
	} else {
		t.Logf("Successfully deleted thread: %s", threadID)
	}
}

func TestDeleteNonExistentThread(t *testing.T) {
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

	nonExistentID := "non_existent_thread_id"
	err = client.DeleteThread(context.Background(), nonExistentID)
	if err == nil {
		t.Error("Expected error when deleting non-existent thread")
	}

	t.Logf("Got expected error: %v", err)
}
