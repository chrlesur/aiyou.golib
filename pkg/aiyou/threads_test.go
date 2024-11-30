// File: pkg/aiyou/threads_test.go
package aiyou

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

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
						User: User{
							ID:    1,
							Email: "test@example.com",
						},
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
				User: User{
					ID:    1,
					Email: "test@example.com",
				},
			})

		case "/api/v1/user/threads":
			// Vérifier les paramètres de requête
			query := r.URL.Query()
			page := query.Get("page")
			itemsPerPage := query.Get("itemsPerPage")
			if page == "" || itemsPerPage == "" {
				t.Error("Missing pagination parameters")
			}

			now := time.Now()
			response := UserThreadsOutput{
				Threads: []ConversationThread{
					{
						ID:                   "thread_123",
						ThreadIdParam:        1,
						Content:              "Test Thread",
						AssistantName:        "Test Assistant",
						AssistantModel:       stringPtr("gpt-4"),
						AssistantId:          1,
						AssistantIdOpenAi:    "asst_123",
						FirstMessage:         "Hello!",
						CreatedAt:            now,
						UpdatedAt:            now,
						IsNewAppThread:       true,
						AssistantContentJson: `{"messages":[{"role":"user","content":[{"type":"text","text":"Hello!"}]}]}`,
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
		if thread.ThreadIdParam != 1 {
			t.Errorf("Expected thread param ID 1, got %d", thread.ThreadIdParam)
		}
		if thread.Content != "Test Thread" {
			t.Errorf("Expected content 'Test Thread', got '%s'", thread.Content)
		}
		if thread.AssistantName != "Test Assistant" {
			t.Errorf("Expected assistant name 'Test Assistant', got '%s'", thread.AssistantName)
		}
		if *thread.AssistantModel != "gpt-4" {
			t.Errorf("Expected assistant model 'gpt-4', got '%s'", *thread.AssistantModel)
		}
		if thread.FirstMessage != "Hello!" {
			t.Errorf("Expected first message 'Hello!', got '%s'", thread.FirstMessage)
		}
	}
}

func TestDeleteThread(t *testing.T) {
	testCases := []struct {
		name      string
		threadID  string
		status    int
		wantError bool
	}{
		{
			name:      "Successful deletion",
			threadID:  "thread_123",
			status:    http.StatusOK,
			wantError: false,
		},
		{
			name:      "Alternative success status",
			threadID:  "thread_123",
			status:    http.StatusNoContent,
			wantError: false,
		},
		{
			name:      "Not Found",
			threadID:  "invalid_thread",
			status:    http.StatusNotFound,
			wantError: true,
		},
		{
			name:      "Server Error",
			threadID:  "thread_123",
			status:    http.StatusInternalServerError,
			wantError: true,
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
						User: User{
							ID:    1,
							Email: "test@example.com",
						},
					})
				case fmt.Sprintf("/api/v1/threads/%s", tc.threadID):
					if r.Method != "DELETE" {
						t.Errorf("Expected DELETE method, got %s", r.Method)
					}
					w.WriteHeader(tc.status)
				}
			}))
			defer server.Close()

			client, err := NewClient("test@example.com", "password", WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			err = client.DeleteThread(context.Background(), tc.threadID)
			if (err != nil) != tc.wantError {
				t.Errorf("DeleteThread() error = %v, wantError %v", err, tc.wantError)
			}
		})
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
