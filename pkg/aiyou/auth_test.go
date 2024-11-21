/*
Copyright (C) 2023 Cloud Temple

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
package aiyou

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
)

func TestJWTAuthenticator_Authenticate(t *testing.T) {
    // Créer un serveur de test
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/api/login" {
            t.Errorf("Expected to request '/api/login', got: %s", r.URL.Path)
        }
        if r.Header.Get("Content-Type") != "application/json" {
            t.Errorf("Expected Content-Type: application/json, got: %s", r.Header.Get("Content-Type"))
        }
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"token":"test_token","expires_at":"2023-01-01T00:00:00Z","user":{"id":"1","email":"test@example.com"}}`))
    }))
    defer server.Close()

    auth := NewJWTAuthenticator("test@example.com", "password", server.URL, server.Client())

    ctx := context.Background()
    err := auth.Authenticate(ctx)

    if err != nil {
        t.Errorf("Authenticate returned an error: %v", err)
    }

    if auth.Token() != "test_token" {
        t.Errorf("Expected token to be 'test_token', got: %s", auth.Token())
    }
}

func TestJWTAuthenticator_TokenExpired(t *testing.T) {
    auth := &JWTAuthenticator{
        token:  "expired_token",
        expiry: time.Now().Add(-1 * time.Hour),
    }

    if !auth.tokenExpired() {
        t.Error("Expected token to be expired")
    }

    auth.expiry = time.Now().Add(1 * time.Hour)
    if auth.tokenExpired() {
        t.Error("Expected token to not be expired")
    }
}