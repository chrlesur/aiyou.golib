// assistants_test.go
package aiyou

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetUserAssistants(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test authentication endpoint
		if r.URL.Path == "/api/login" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(LoginResponse{
				Token:     "test_token",
				ExpiresAt: time.Now().Add(time.Hour),
			})
			return
		}

		// Test assistants endpoint
		if r.URL.Path == "/api/v1/user/assistants" {
			response := AssistantsResponse{
				Assistants: []Assistant{
					{
						ID:          "asst_123",
						Name:        "Test Assistant",
						Description: "A test assistant",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
						ModelID:     "model_123",
						IsPublic:    true,
					},
				},
				Total: 1,
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	// Create client
	client, err := NewClient("test@example.com", "password", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test GetUserAssistants
	assistantsResp, err := client.GetUserAssistants(context.Background())
	if err != nil {
		t.Fatalf("GetUserAssistants failed: %v", err)
	}

	if len(assistantsResp.Assistants) != 1 {
		t.Errorf("Expected 1 assistant, got %d", len(assistantsResp.Assistants))
	}

	if assistantsResp.Assistants[0].Name != "Test Assistant" {
		t.Errorf("Expected assistant name 'Test Assistant', got '%s'", assistantsResp.Assistants[0].Name)
	}
}
