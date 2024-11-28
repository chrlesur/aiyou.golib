// File: pkg/aiyou/threads_test.go
package aiyou

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetUserThreads(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/login" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(LoginResponse{
				Token:     "test_token",
				ExpiresAt: time.Now().Add(time.Hour),
			})
			return
		}

		if r.URL.Path == "/api/v1/user/threads" {
			// Vérifier les paramètres de requête
			query := r.URL.Query()
			if assistantID := query.Get("assistantId"); assistantID != "" {
				t.Logf("Received assistantId filter: %s", assistantID)
			}

			response := ThreadsResponse{
				Threads: []ConversationThread{
					{
						ID:          "thread_123",
						Title:       "Test Thread",
						AssistantID: "asst_123",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				},
				Total: 1,
				Page:  1,
				Limit: 10,
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	client, err := NewClient("test@example.com", "password", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	filter := &ThreadFilter{
		AssistantID: "asst_123",
		Page:        1,
		Limit:       10,
	}

	resp, err := client.GetUserThreads(context.Background(), filter)
	if err != nil {
		t.Fatalf("GetUserThreads failed: %v", err)
	}

	if len(resp.Threads) != 1 {
		t.Errorf("Expected 1 thread, got %d", len(resp.Threads))
	}

	if resp.Threads[0].ID != "thread_123" {
		t.Errorf("Expected thread ID 'thread_123', got '%s'", resp.Threads[0].ID)
	}
}

func TestDeleteThread(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/login" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(LoginResponse{
				Token:     "test_token",
				ExpiresAt: time.Now().Add(time.Hour),
			})
			return
		}

		if r.URL.Path == "/api/v1/threads/thread_123" && r.Method == "DELETE" {
			w.WriteHeader(http.StatusOK)
			return
		}
	}))
	defer server.Close()

	client, err := NewClient("test@example.com", "password", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.DeleteThread(context.Background(), "thread_123")
	if err != nil {
		t.Errorf("DeleteThread failed: %v", err)
	}
}
