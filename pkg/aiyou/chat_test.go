// File: pkg/aiyou/chat_test.go
package aiyou

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
)

// Variables de test
var (
	testLoginResponse = LoginResponse{
		Token:     "test_token",
		ExpiresAt: time.Now().Add(time.Hour),
		User: User{
			ID:    1,
			Email: "test@example.com",
		},
	}

	testMessage = Message{
		Role: "user",
		Content: []ContentPart{
			{Type: "text", Text: "Hello!"},
		},
	}

	testChatRequest = ChatCompletionRequest{
		Messages:    []Message{testMessage},
		AssistantID: "asst_123",
		Temperature: 7,
		TopP:        1.0,
		Stream:      false,
	}
)

func createTestServer(t *testing.T, handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

func createTestClient(t *testing.T, serverURL string) *Client {
	client, err := NewClient(
		"test@example.com",
		"password",
		WithBaseURL(serverURL),
	)
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}
	return client
}

func TestChatCompletion(t *testing.T) {
	testCases := []struct {
		name           string
		nonStreamResp  func(w http.ResponseWriter)
		streamResp     func(w http.ResponseWriter, flusher http.Flusher)
		expectedText   string
		expectFallback bool
	}{
		{
			name: "Direct non-streaming success",
			nonStreamResp: func(w http.ResponseWriter) {
				response := ChatCompletionResponse{
					ID:      "resp_123",
					Object:  "chat.completion",
					Created: time.Now().Unix(),
					Model:   "gpt-4",
					Choices: []Choice{
						{
							Message: Message{
								Role: "assistant",
								Content: []ContentPart{
									{Type: "text", Text: "Hello, how can I help you today?"},
								},
							},
							FinishReason: "stop",
						},
					},
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(response)
			},
			streamResp:     nil,
			expectedText:   "Hello, how can I help you today?",
			expectFallback: false,
		},
		{
			name: "Fallback to streaming",
			nonStreamResp: func(w http.ResponseWriter) {
				errorResp := map[string]interface{}{
					"object":  "error",
					"message": "Stream options can only be defined when `stream=True`",
					"type":    "BadRequestError",
					"code":    400,
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(errorResp)
			},
			streamResp: func(w http.ResponseWriter, flusher http.Flusher) {
				w.Header().Set("Content-Type", "text/event-stream")
				w.Header().Set("Cache-Control", "no-cache")
				w.Header().Set("Connection", "keep-alive")
				w.WriteHeader(http.StatusOK)

				messages := []string{"Hello", ", how", " can I", " help you", " today?"}
				for _, msg := range messages {
					chunk := ChatCompletionResponse{
						ID:      "resp_123",
						Object:  "chat.completion.chunk",
						Created: time.Now().Unix(),
						Model:   "gpt-4",
						Choices: []Choice{
							{
								Delta: &Delta{
									Content: msg,
								},
							},
						},
					}
					data, _ := json.Marshal(chunk)
					fmt.Fprintf(w, "%s\n\n", data)
					flusher.Flush()
					time.Sleep(10 * time.Millisecond)
				}
			},
			expectedText:   "Hello, how can I help you today?",
			expectFallback: true,
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
					if r.Method != "POST" {
						t.Errorf("Expected POST request, got: %s", r.Method)
					}

					var reqBody map[string]interface{}
					if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
						t.Errorf("Failed to decode request body: %v", err)
						return
					}

					isStream, _ := reqBody["stream"].(bool)
					if isStream {
						if tc.streamResp != nil {
							flusher, ok := w.(http.Flusher)
							if !ok {
								t.Fatal("Streaming not supported")
								return
							}
							tc.streamResp(w, flusher)
						}
					} else {
						tc.nonStreamResp(w)
					}
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

			gotText := resp.Choices[0].Message.Content[0].Text
			if gotText != tc.expectedText {
				t.Errorf("Expected text %q, got %q", tc.expectedText, gotText)
			}
		})
	}
}

func TestChatCompletionStream(t *testing.T) {
	server := createTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/login":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(testLoginResponse)

		case "/api/v1/chat/completions":
			var reqBody ChatCompletionRequest
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
				chunk := ChatCompletionResponse{
					ID:      "resp_123",
					Object:  "chat.completion.chunk",
					Created: time.Now().Unix(),
					Model:   "gpt-4",
					Choices: []Choice{
						{
							Message: Message{
								Role: "assistant",
								Content: []ContentPart{
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

				fmt.Fprintf(w, "%s\n\n", data)
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
			continue
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
		t.Errorf("Expected response %q, got %q", expected, fullResponse.String())
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
				json.NewEncoder(w).Encode(map[string]interface{}{
					"object":  "error",
					"message": "Unauthorized",
					"type":    "AuthenticationError",
					"code":    401,
				})
			},
			expectedError: true,
		},
		{
			name: "Bad Request",
			serverResponse: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"object":  "error",
					"message": "Invalid request",
					"type":    "BadRequestError",
					"code":    400,
				})
			},
			expectedError: true,
		},
		{
			name: "Rate Limit",
			serverResponse: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"object":  "error",
					"message": "Rate limit exceeded",
					"type":    "RateLimitError",
					"code":    429,
				})
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

			if tc.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else {
					// Vérifier le type d'erreur
					switch tc.name {
					case "Unauthorized":
						if _, ok := err.(*APIError); !ok {
							t.Errorf("Expected APIError, got %T", err)
						}
					case "Bad Request":
						if _, ok := err.(*APIError); !ok {
							t.Errorf("Expected APIError, got %T", err)
						}
					case "Rate Limit":
						if _, ok := err.(*RateLimitError); !ok {
							t.Errorf("Expected RateLimitError, got %T", err)
						}
					}
				}
			} else if err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
		})
	}
}
