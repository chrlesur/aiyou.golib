// File: pkg/aiyou/conversation_test.go
package aiyou

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSaveConversation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/login" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(LoginResponse{
				Token:     "test_token",
				ExpiresAt: time.Now().Add(time.Hour),
			})
			return
		}

		if r.URL.Path == "/api/v1/save" && r.Method == "POST" {
			var req SaveConversationRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("Failed to decode request body: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			response := SaveConversationResponse{
				Thread: ConversationThread{
					ID:          "thread_123",
					Title:       req.Title,
					Messages:    req.Messages,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					AssistantID: req.AssistantID,
				},
				Status: "success",
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

	req := SaveConversationRequest{
		Title:       "Test Conversation",
		AssistantID: "asst_123",
		Messages: []Message{
			{
				Role: "user",
				Content: []ContentPart{
					{Type: "text", Text: "Hello!"},
				},
			},
		},
	}

	resp, err := client.SaveConversation(context.Background(), req)
	if err != nil {
		t.Fatalf("SaveConversation failed: %v", err)
	}

	if resp.Thread.Title != req.Title {
		t.Errorf("Expected thread title %s, got %s", req.Title, resp.Thread.Title)
	}
}
