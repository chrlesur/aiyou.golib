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
		// Test authentication endpoint
		if r.URL.Path == "/api/login" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(LoginResponse{
				Token:     "test_token",
				ExpiresAt: time.Now().Add(time.Hour),
			})
			return
		}

		// Test save conversation endpoint
		if r.URL.Path == "/api/v1/save" && r.Method == "POST" {
			// Décoder la requête pour vérification
			var req SaveConversationRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("Failed to decode request body: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Vérifier les champs requis
			if req.AssistantID == "" {
				t.Error("AssistantID is required")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Créer une réponse de test
			response := SaveConversationResponse{
				ID:        "thread_123",
				Object:    "thread",
				CreatedAt: time.Now().Unix(),
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	client, err := NewClient("test@example.com", "password", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Créer une requête de test
	req := SaveConversationRequest{
		AssistantID:    "asst_123",
		Conversation:   "Test conversation",
		FirstMessage:   "Hello, how can I help you?",
		ContentJson:    `{"messages":[{"role":"user","content":[{"type":"text","text":"Hello"}]}]}`,
		ModelName:      "gpt-4",
		IsNewAppThread: true,
	}

	resp, err := client.SaveConversation(context.Background(), req)
	if err != nil {
		t.Fatalf("SaveConversation failed: %v", err)
	}

	// Vérifier la réponse
	if resp.ID == "" {
		t.Error("Expected thread ID in response")
	}

	if resp.Object != "thread" {
		t.Errorf("Expected object type 'thread', got %s", resp.Object)
	}

	if resp.CreatedAt == 0 {
		t.Error("Expected non-zero CreatedAt timestamp")
	}
}

func TestSaveConversationErrors(t *testing.T) {
	testCases := []struct {
		name           string
		serverResponse func(w http.ResponseWriter)
		expectedError  bool
	}{
		{
			name: "Unauthorized",
			serverResponse: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Unauthorized",
				})
			},
			expectedError: true,
		},
		{
			name: "Bad Request",
			serverResponse: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Invalid request",
				})
			},
			expectedError: true,
		},
		{
			name: "Server Error",
			serverResponse: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/api/login" {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(LoginResponse{
						Token:     "test_token",
						ExpiresAt: time.Now().Add(time.Hour),
					})
					return
				}

				if r.URL.Path == "/api/v1/save" {
					tc.serverResponse(w)
				}
			}))
			defer server.Close()

			client, err := NewClient("test@example.com", "password", WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			req := SaveConversationRequest{
				AssistantID:  "asst_123",
				Conversation: "Test conversation",
				FirstMessage: "Hello",
				ContentJson:  "{}",
				ModelName:    "gpt-4",
			}

			_, err = client.SaveConversation(context.Background(), req)
			if (err != nil) != tc.expectedError {
				t.Errorf("Expected error: %v, got error: %v", tc.expectedError, err)
			}
		})
	}
}

func TestSaveConversationValidation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/login" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(LoginResponse{
				Token: "test_token",
				ExpiresAt: time.Now().Add(time.Hour),
			})
			return
		}

		if r.URL.Path == "/api/v1/save" {
			var req SaveConversationRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Validation des champs requis
			if req.AssistantID == "" || req.Conversation == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Réponse de succès
			w.WriteHeader(http.StatusCreated) // Changé de 200 à 201
			response := SaveConversationResponse{
				ID: "thread_123",
				Object: "thread",
				CreatedAt: time.Now().Unix(),
			}
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	testCases := []struct {
		name string
		request SaveConversationRequest
		wantErr bool
	}{
		{
			name: "Empty AssistantID",
			request: SaveConversationRequest{
				AssistantID: "",
				Conversation: "Test",
				FirstMessage: "Hello",
			},
			wantErr: true,
		},
		{
			name: "Empty Conversation",
			request: SaveConversationRequest{
				AssistantID: "asst_123",
				Conversation: "",
				FirstMessage: "Hello",
			},
			wantErr: true,
		},
		{
			name: "Valid Request",
			request: SaveConversationRequest{
				AssistantID: "asst_123",
				Conversation: "Test conversation",
				FirstMessage: "Hello",
				ContentJson: "{}",
				ModelName: "gpt-4",
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, err := NewClient("test@example.com", "password", WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			_, err = client.SaveConversation(context.Background(), tc.request)
			if (err != nil) != tc.wantErr {
				t.Errorf("Expected error: %v, got error: %v", tc.wantErr, err)
			}
		})
	}
}