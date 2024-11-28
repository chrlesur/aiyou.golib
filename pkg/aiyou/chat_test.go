// File: pkg/aiyou/chat_test.go
package aiyou_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/chrlesur/aiyou.golib/pkg/aiyou"
)

// Variables de test
var (
	testLoginResponse = aiyou.LoginResponse{
		Token:     "test_token",
		ExpiresAt: time.Now().Add(time.Hour),
		User: aiyou.User{
			ID:    "1",
			Email: "test@example.com",
		},
	}

	testMessage = aiyou.Message{
		Role: "user",
		Content: []aiyou.ContentPart{
			{Type: "text", Text: "Hello!"},
		},
	}

	testChatRequest = aiyou.ChatCompletionRequest{
		Messages:    []aiyou.Message{testMessage},
		AssistantID: "asst_123",
		Temperature: 0.7,
		TopP:        1.0,
		Stream:      false,
	}
)

func createTestServer(t *testing.T, handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

func createTestClient(t *testing.T, serverURL string) *aiyou.Client {
	client, err := aiyou.NewClient(
		"test@example.com",
		"password",
		aiyou.WithBaseURL(serverURL),
	)
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}
	return client
}

func TestChatCompletion(t *testing.T) {
	server := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/login":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(testLoginResponse)

		case "/api/v1/chat/completions":
			if r.Method != "POST" {
				t.Errorf("Expected POST request, got: %s", r.Method)
			}

			if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
				t.Error("Missing or invalid Authorization header")
			}

			var reqBody aiyou.ChatCompletionRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Errorf("Failed to decode request body: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			response := aiyou.ChatCompletionResponse{
				ID:      "resp_123",
				Object:  "chat.completion",
				Created: time.Now().Unix(),
				Model:   "gpt-4",
				Choices: []aiyou.Choice{
					{
						Message: aiyou.Message{
							Role: "assistant",
							Content: []aiyou.ContentPart{
								{Type: "text", Text: "Hello, how can I help you today?"},
							},
						},
						FinishReason: "stop",
					},
				},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}
	})
	defer server.Close()

	client := createTestClient(t, server.URL)
	resp, err := client.ChatCompletion(context.Background(), testChatRequest)
	if err != nil {
		t.Fatalf("ChatCompletion failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected non-nil response")
	}

	if len(resp.Choices) == 0 {
		t.Fatal("Expected at least one choice in response")
	}

	if resp.Choices[0].Message.Content[0].Text != "Hello, how can I help you today?" {
		t.Errorf("Unexpected response text: %s", resp.Choices[0].Message.Content[0].Text)
	}
}

func TestChatCompletionStream(t *testing.T) {
	server := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/login":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(testLoginResponse)

		case "/api/v1/chat/completions":
			// Vérifier que c'est une requête de streaming
			var reqBody aiyou.ChatCompletionRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Errorf("Failed to decode request body: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if !reqBody.Stream {
				t.Error("Expected stream to be true")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			flusher, ok := w.(http.Flusher)
			if !ok {
				t.Fatal("Streaming not supported")
				return
			}

			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			w.WriteHeader(http.StatusOK)
			flusher.Flush()

			messages := []string{"Hello", ", how", " can I", " help you", " today?"}
			for _, msg := range messages {
				chunk := aiyou.ChatCompletionResponse{
					ID: "resp_123",
					Object: "chat.completion.chunk",
					Created: time.Now().Unix(),
					Model: "gpt-4",
					Choices: []aiyou.Choice{
						{
							Message: aiyou.Message{
								Role: "assistant",
								Content: []aiyou.ContentPart{
									{Type: "text", Text: msg},
								},
							},
						},
					},
				}

				data, err := json.Marshal(chunk)
				if err != nil {
					t.Errorf("Failed to marshal chunk: %v", err)
					return
				}

				// Format SSE
				_, err = fmt.Fprintf(w, "%s\n\n", data)
				if err != nil {
					t.Errorf("Failed to write chunk: %v", err)
					return
				}
				flusher.Flush()
				time.Sleep(50 * time.Millisecond)
			}
		}
	})
	defer server.Close()

	client := createTestClient(t, server.URL)
	
	streamReq := testChatRequest
	streamReq.Stream = true

	stream, err := client.ChatCompletionStream(context.Background(), streamReq)
	if err != nil {
		t.Fatalf("ChatCompletionStream failed: %v", err)
	}
	if stream == nil {
		t.Fatal("Expected non-nil stream")
	}
	defer stream.Close()

	var fullResponse strings.Builder
	for {
		chunk, err := stream.ReadChunk()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Error reading chunk: %v", err)
		}
		if chunk == nil {
			continue // Skip empty chunks
		}

		if len(chunk.Choices) > 0 {
			message := chunk.Choices[0].Message
			if len(message.Content) > 0 {
				fullResponse.WriteString(message.Content[0].Text)
			}
		}
	}

	expected := "Hello, how can I help you today?"
	if fullResponse.String() != expected {
		t.Errorf("Expected response '%s', got '%s'", expected, fullResponse.String())
	}
}

func TestChatCompletionErrors(t *testing.T) {
	testCases := []struct {
		name           string
		serverResponse func(w http.ResponseWriter)
		expectedError  bool
	}{
		{
			name: "Unauthorized",
			serverResponse: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
			},
			expectedError: true,
		},
		{
			name: "Bad Request",
			serverResponse: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
			},
			expectedError: true,
		},
		{
			name: "Rate Limit",
			serverResponse: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]string{"error": "Rate limit exceeded"})
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/api/login":
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(testLoginResponse)
				case "/api/v1/chat/completions":
					tc.serverResponse(w)
				}
			})
			defer server.Close()

			client := createTestClient(t, server.URL)
			_, err := client.ChatCompletion(context.Background(), testChatRequest)
			if (err != nil) != tc.expectedError {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}
