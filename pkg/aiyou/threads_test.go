// File: pkg/aiyou/threads_test.go
package aiyou

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// ... (garder TestGetUserThreads inchangé)

func TestGetUserThreadsError(t *testing.T) {
	testCases := []struct {
		name      string
		setup     func(w http.ResponseWriter)
		wantErr   bool
		errSubstr string
	}{
		{
			name: "Unauthorized",
			setup: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Unauthorized",
				})
			},
			wantErr:   true,
			errSubstr: "status code: 401",
		},
		{
			name: "Server Error",
			setup: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Internal Server Error",
				})
			},
			wantErr:   true,
			errSubstr: "status code: 500",
		},
		{
			name: "Invalid JSON",
			setup: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("invalid json"))
			},
			wantErr:   true,
			errSubstr: "failed to decode",
		},
		{
			name: "Empty Response",
			setup: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(""))
			},
			wantErr:   true,
			errSubstr: "EOF",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/api/login":
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(LoginResponse{
						Token:     "test_token",
						ExpiresAt: time.Now().Add(time.Hour),
					})
				case "/api/v1/user/threads":
					tc.setup(w)
				}
			}))
			defer server.Close()

			client, err := NewClient(
				"test@example.com",
				"password",
				WithBaseURL(server.URL),
			)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			params := &UserThreadsParams{
				Page:         1,
				ItemsPerPage: 10,
			}

			_, err = client.GetUserThreads(context.Background(), params)

			// Vérification des erreurs
			if tc.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}
				if tc.errSubstr != "" && !strings.Contains(err.Error(), tc.errSubstr) {
					t.Errorf("Expected error containing %q, got %q", tc.errSubstr, err.Error())
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestGetUserThreads(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/login":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(LoginResponse{
				Token:     "test_token",
				ExpiresAt: time.Now().Add(time.Hour),
			})

		case "/api/v1/user/threads":
			// Vérifier les paramètres de requête
			query := r.URL.Query()
			page := query.Get("page")
			itemsPerPage := query.Get("itemsPerPage")
			if page == "" || itemsPerPage == "" {
				t.Error("Missing pagination parameters")
			}

			response := UserThreadsOutput{
				Threads: []ConversationThread{
					{
						ID:          "thread_123",
						Title:       "Test Thread",
						AssistantID: "asst_123",
						Messages: []Message{
							{
								Role: "user",
								Content: []ContentPart{
									{
										Type: "text",
										Text: "Hello!",
									},
								},
							},
						},
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
						LastMessageAt: time.Now(),
					},
				},
				TotalItems:   1,
				ItemsPerPage: 10,
				CurrentPage:  1,
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

	params := &UserThreadsParams{
		Page:         1,
		ItemsPerPage: 10,
	}

	response, err := client.GetUserThreads(context.Background(), params)
	if err != nil {
		t.Fatalf("GetUserThreads failed: %v", err)
	}

	// Vérifications
	if response.CurrentPage != 1 {
		t.Errorf("Expected current page 1, got %d", response.CurrentPage)
	}

	if response.ItemsPerPage != 10 {
		t.Errorf("Expected items per page 10, got %d", response.ItemsPerPage)
	}

	if len(response.Threads) != 1 {
		t.Errorf("Expected 1 thread, got %d", len(response.Threads))
	}

	// Vérifier le contenu du thread
	if len(response.Threads) > 0 {
		thread := response.Threads[0]
		if thread.ID != "thread_123" {
			t.Errorf("Expected thread ID 'thread_123', got '%s'", thread.ID)
		}
		if thread.Title != "Test Thread" {
			t.Errorf("Expected thread title 'Test Thread', got '%s'", thread.Title)
		}
		if len(thread.Messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(thread.Messages))
		}
	}
}
