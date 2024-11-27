/*
Copyright (C) 2024 Cloud Temple

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
// File: pkg/aiyou/auth_test.go

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

func TestJWTAuthenticator_Authenticate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/login" {
			t.Errorf("Expected to request '/api/login', got: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got: %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(LoginResponse{
			Token:     "test_token",
			ExpiresAt: time.Now().Add(time.Hour),
		})
	}))
	defer server.Close()

	// Créer un client HTTP avec un timeout
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	auth := NewJWTAuthenticator("test@example.com", "password", server.URL, httpClient, NewDefaultLogger(io.Discard))

	err := auth.Authenticate(context.Background())
	if err != nil {
		t.Errorf("Authenticate failed: %v", err)
	}

	if auth.Token() == "" {
		t.Error("Expected token to be set after authentication")
	}
}

func TestJWTAuthenticator_TokenExpired(t *testing.T) {
	auth := &JWTAuthenticator{
		token:  "expired_token",
		expiry: time.Now().Add(-1 * time.Hour),
		logger: NewDefaultLogger(io.Discard),
	}

	if !auth.tokenExpired() {
		t.Error("Expected token to be expired")
	}

	auth.expiry = time.Now().Add(1 * time.Hour)
	if auth.tokenExpired() {
		t.Error("Expected token to not be expired")
	}
}
