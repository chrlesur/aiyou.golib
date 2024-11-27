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

// File: pkg/aiyou/client_test.go

package aiyou

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClientWithRetry(t *testing.T) {
	var logBuf bytes.Buffer
	customLogger := NewDefaultLogger(&logBuf)
	customLogger.SetLevel(DEBUG)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/login" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"token":"test_token","expires_at":"2099-01-01T00:00:00Z"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient("test@example.com", "password",
		WithBaseURL(server.URL),
		WithRetry(3, time.Millisecond),
		WithLogger(customLogger))

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Vérifier que SafeLog est correctement initialisé
	if client.safeLog == nil {
		t.Fatalf("SafeLog is not initialized")
	}

	// Test direct de SafeLog
	client.safeLog(INFO, "Test SafeLog with email: %s", "test@example.com")

	ctx := context.Background()
	_, err = client.AuthenticatedRequest(ctx, "GET", "/test", nil)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	logOutput := logBuf.String()
	t.Logf("Full log output:\n%s", logOutput)

	// Vérifier que les entrées de log contiennent le nom du fichier et le numéro de ligne
	logEntries := strings.Split(logOutput, "\n")
	for _, entry := range logEntries {
		if entry == "" {
			continue
		}
		if !strings.Contains(entry, ".go:") {
			t.Errorf("Log entry does not contain file name and line number: %s", entry)
		}
	}

	// Vérifier que les adresses e-mail sont masquées
	if strings.Contains(logOutput, "test@example.com") {
		t.Errorf("Log output contains unmasked email address")
	}

	if !strings.Contains(logOutput, "[EMAIL REDACTED]") {
		t.Errorf("Log output does not contain '[EMAIL REDACTED]'")
	}

	// Vérifier que le statut de la réponse est correct
	if !strings.Contains(logOutput, "Request completed with status: 200") {
		t.Errorf("Log output does not contain correct status code")
	}

	// Vérifier la présence de messages de log spécifiques
	expectedLogEntries := []string{
		"Preparing authenticated request: GET /test",
		"Sending request to",
		"Request completed with status: 200",
	}

	for _, entry := range expectedLogEntries {
		if !strings.Contains(logOutput, entry) {
			t.Errorf("Expected log to contain %q, but it didn't.", entry)
		}
	}

	// Vérifier l'absence d'informations sensibles
	sensitiveInfo := []string{
		"password",
		"Bearer test_token",
	}

	for _, info := range sensitiveInfo {
		if strings.Contains(logOutput, info) {
			t.Errorf("Log should not contain sensitive info %q, but it did.", info)
		}
	}
}
