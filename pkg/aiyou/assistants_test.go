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

func TestGetUserAssistantsErrors(t *testing.T) {
	testCases := []struct {
		name        string
		setup       func(w http.ResponseWriter)
		wantErr     bool
		errContains string
	}{
		{
			name: "Unauthorized",
			setup: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Unauthorized",
				})
			},
			wantErr:     true,
			errContains: "status code: 401",
		},
		{
			name: "Bad Request",
			setup: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Bad Request",
				})
			},
			wantErr:     true,
			errContains: "status code: 400",
		},
		{
			name: "Server Error",
			setup: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Internal Server Error",
				})
			},
			wantErr:     true,
			errContains: "status code: 500",
		},
		{
			name: "Invalid Response",
			setup: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("invalid json"))
			},
			wantErr:     true,
			errContains: "failed to decode",
		},
		{
			name: "Empty Response",
			setup: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(""))
			},
			wantErr:     true,
			errContains: "EOF",
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
				case "/api/v1/user/assistants":
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

			_, err = client.GetUserAssistants(context.Background())

			if tc.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tc.errContains) {
					t.Errorf("Expected error containing %q, got %q", tc.errContains, err.Error())
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
