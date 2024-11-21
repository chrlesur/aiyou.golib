// Fichier: pkg/aiyou/chat_test.go

package aiyou

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChatCompletion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/login" {
			// Simuler une réponse d'authentification réussie
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"token":"test_token","expires_at":"2099-01-01T00:00:00Z"}`))
			return
		}
		if r.URL.Path != "/api/v1/chat/completions" {
			t.Errorf("Expected to request '/api/v1/chat/completions', got: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got: %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"chatcmpl-123","object":"chat.completion","created":1677652288,"model":"gpt-3.5-turbo-0613","choices":[{"index":0,"message":{"role":"assistant","content":[{"type":"text","text":"Hello, how can I assist you today?"}]},"finish_reason":"stop"}],"usage":{"prompt_tokens":9,"completion_tokens":12,"total_tokens":21}}`))
	}))
	defer server.Close()

	// Create a client using the mock server URL
	client, err := NewClient("test@example.com", "password", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Test ChatCompletion
	req := ChatCompletionRequest{
		Messages: []Message{
			{Role: "user", Content: []ContentPart{{Type: "text", Text: "Hello!"}}},
		},
		AssistantID: "test-assistant",
	}
	resp, err := client.ChatCompletion(context.Background(), req)
	if err != nil {
		t.Fatalf("Error in ChatCompletion: %v", err)
	}

	if resp.Choices[0].Message.Content[0].Text != "Hello, how can I assist you today?" {
		t.Errorf("Unexpected response content: %s", resp.Choices[0].Message.Content[0].Text)
	}
}

func TestChatCompletionStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/login" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"token":"test_token","expires_at":"2099-01-01T00:00:00Z"}`))
			return
		}
		if r.URL.Path != "/api/v1/chat/completions" {
			t.Errorf("Expected to request '/api/v1/chat/completions', got: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		// Simulate streaming response
		for i := 0; i < 3; i++ {
			chunk := ChatCompletionResponse{
				ID:      "chatcmpl-123",
				Object:  "chat.completion.chunk",
				Created: 1677652288,
				Model:   "gpt-3.5-turbo-0613",
				Choices: []Choice{
					{
						Message: Message{
							Role:    "assistant",
							Content: []ContentPart{{Type: "text", Text: "Hello"}},
						},
						FinishReason: "",
					},
				},
			}
			jsonData, _ := json.Marshal(chunk)
			_, err := fmt.Fprintf(w, "data: %s\n\n", jsonData)
			if err != nil {
				t.Errorf("Error writing chunk: %v", err)
			}
			w.(http.Flusher).Flush()
		}
	}))
	defer server.Close()

	// Create a client using the mock server URL
	client, err := NewClient("test@example.com", "password", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Test ChatCompletionStream
	req := ChatCompletionRequest{
		Messages: []Message{
			{Role: "user", Content: []ContentPart{{Type: "text", Text: "Hello!"}}},
		},
		AssistantID: "test-assistant",
		Stream:      true,
	}
	stream, err := client.ChatCompletionStream(context.Background(), req)
	if err != nil {
		t.Fatalf("Error in ChatCompletionStream: %v", err)
	}

	chunkCount := 0
	for {
		chunk, err := stream.ReadChunk()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Error reading chunk: %v", err)
		}
		chunkCount++
		if chunk.Choices[0].Message.Content[0].Text != "Hello" {
			t.Errorf("Unexpected chunk content: %s", chunk.Choices[0].Message.Content[0].Text)
		}
	}

	if chunkCount != 3 {
		t.Errorf("Expected 3 chunks, got %d", chunkCount)
	}
}
