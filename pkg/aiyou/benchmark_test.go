package aiyou

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// setupTestServer crée un serveur de test pour les benchmarks
func setupTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/login":
			json.NewEncoder(w).Encode(LoginResponse{
				Token:     "test_token",
				ExpiresAt: time.Now().Add(time.Hour),
			})
		case "/api/v1/chat/completions":
			response := ChatCompletionResponse{
				ID:      "test_id",
				Object:  "chat.completion",
				Created: time.Now().Unix(),
				Model:   "test_model",
				Choices: []Choice{
					{
						Message: Message{
							Role: "assistant",
							Content: []ContentPart{
								{Type: "text", Text: "Test response"},
							},
						},
					},
				},
			}
			json.NewEncoder(w).Encode(response)
		}
	}))
}

func BenchmarkChatCompletion(b *testing.B) {
	// Configuration du serveur de test
	server := setupTestServer()
	defer server.Close()

	// Créer le client avec un logger silencieux
	client, err := NewClient(
		WithEmailPassword("test@example.com", "password"),
		WithBaseURL(server.URL),
		WithLogger(NewDefaultLogger(io.Discard)),
	)
	if err != nil {
		b.Fatalf("Failed to create client: %v", err)
	}

	// Préparer la requête de test
	req := ChatCompletionRequest{
		Messages: []Message{
			{
				Role: "user",
				Content: []ContentPart{
					{Type: "text", Text: "Hello"},
				},
			},
		},
		AssistantID: "test-assistant",
	}

	b.ResetTimer()

	// Exécuter le benchmark
	for i := 0; i < b.N; i++ {
		_, err := client.ChatCompletion(context.Background(), req)
		if err != nil {
			b.Fatalf("ChatCompletion failed: %v", err)
		}
	}
}
