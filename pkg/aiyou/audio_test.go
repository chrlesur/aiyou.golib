// File: pkg/aiyou/audio_test.go
package aiyou

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestTranscribeAudioFile(t *testing.T) {
	// Créer un fichier audio de test temporaire
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.mp3")
	if err := os.WriteFile(testFile, []byte("fake audio content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/login" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(LoginResponse{
				Token:     "test_token",
				ExpiresAt: time.Now().Add(time.Hour),
			})
			return
		}

		if r.URL.Path == "/api/v1/audio/transcriptions" {
			// Vérifier que le Content-Type contient multipart/form-data
			contentType := r.Header.Get("Content-Type")
			if !strings.Contains(contentType, "multipart/form-data") {
				t.Errorf("Expected Content-Type to contain multipart/form-data, got %s", contentType)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Vérifier l'authentification
			authHeader := r.Header.Get("Authorization")
			if authHeader != "Bearer test_token" {
				t.Errorf("Expected Authorization header 'Bearer test_token', got %s", authHeader)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Parser le formulaire multipart
			if err := r.ParseMultipartForm(32 << 20); err != nil {
				t.Errorf("Failed to parse multipart form: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Vérifier la présence du fichier
			file, _, err := r.FormFile("file")
			if err != nil {
				t.Errorf("Failed to get file from form: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			defer file.Close()

			// Renvoyer une réponse de succès
			response := AudioTranscriptionResponse{
				Text:      "This is a test transcription",
				Language:  "en",
				Duration:  10.5,
				CreatedAt: time.Now(),
				Status:    "completed",
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

	opts := &AudioTranscriptionRequest{
		Language: "en",
		Format:   "text",
	}

	resp, err := client.TranscribeAudioFile(context.Background(), testFile, opts)
	if err != nil {
		t.Fatalf("TranscribeAudioFile failed: %v", err)
	}

	if resp.Text != "This is a test transcription" {
		t.Errorf("Expected transcription text 'This is a test transcription', got '%s'", resp.Text)
	}

	if resp.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", resp.Status)
	}
}

func TestValidateAudioFile(t *testing.T) {
	tests := []struct {
		name          string
		fileContent   []byte
		fileName      string
		expectedError bool
	}{
		{
			name:          "Valid MP3 file",
			fileContent:   make([]byte, 1000),
			fileName:      "test.mp3",
			expectedError: false,
		},
		{
			name:          "Unsupported format",
			fileContent:   make([]byte, 1000),
			fileName:      "test.xyz",
			expectedError: true,
		},
		{
			name:          "File too large",
			fileContent:   make([]byte, 26*1024*1024),
			fileName:      "test.mp3",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, tt.fileName)

			if err := os.WriteFile(testFile, tt.fileContent, 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			file, err := os.Open(testFile)
			if err != nil {
				t.Fatalf("Failed to open test file: %v", err)
			}
			defer file.Close()

			err = validateAudioFile(file, testFile)
			if (err != nil) != tt.expectedError {
				t.Errorf("validateAudioFile() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}
