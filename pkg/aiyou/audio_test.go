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

func TestAudioFileValidation(t *testing.T) {
	testCases := []struct {
		name    string
		content []byte
		ext     string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid MP3",
			content: make([]byte, 1000),
			ext:     ".mp3",
			wantErr: false,
		},
		{
			name:    "Valid WAV",
			content: make([]byte, 1000),
			ext:     ".wav",
			wantErr: false,
		},
		{
			name:    "Invalid extension",
			content: make([]byte, 1000),
			ext:     ".xyz",
			wantErr: true,
			errMsg:  "unsupported audio format: .xyz",
		},
		{
			name:    "File too large",
			content: make([]byte, 26*1024*1024), // Plus grand que la limite
			ext:     ".mp3",
			wantErr: true,
			errMsg:  "file size exceeds maximum allowed size",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Créer un fichier temporaire
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "test"+tc.ext)
			if err := os.WriteFile(testFile, tc.content, 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Ouvrir le fichier pour la validation
			file, err := os.Open(testFile)
			if err != nil {
				t.Fatalf("Failed to open test file: %v", err)
			}
			defer file.Close()

			// Tester la validation
			err = validateAudioFile(file, testFile)

			// Vérifier les résultats
			if tc.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				} else if tc.errMsg != "" && !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("Expected error containing %q, got %q", tc.errMsg, err.Error())
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestSupportedAudioFormats(t *testing.T) {
	// Vérifier que les formats supportés sont correctement définis
	for _, format := range SupportedFormats {
		t.Run(format.Extension, func(t *testing.T) {
			if format.MaxSize <= 0 {
				t.Errorf("Format %s has invalid max size: %d", format.Extension, format.MaxSize)
			}
			if len(format.MimeTypes) == 0 {
				t.Errorf("Format %s has no mime types", format.Extension)
			}
			// Vérifier que l'extension commence par un point
			if format.Extension[0] != '.' {
				t.Errorf("Format extension %s should start with a dot", format.Extension)
			}
		})
	}
}

func TestTranscribeAudioFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/login":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(LoginResponse{
				Token:     "test_token",
				ExpiresAt: time.Now().Add(time.Hour),
			})

		case "/api/v1/audio/transcriptions":
			// Vérifier la méthode
			if r.Method != "POST" {
				t.Errorf("Expected POST request, got: %s", r.Method)
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			// Vérifier le Content-Type
			contentType := r.Header.Get("Content-Type")
			if !strings.HasPrefix(contentType, "multipart/form-data") {
				t.Errorf("Expected multipart/form-data, got: %s", contentType)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Vérifier le fichier
			if err := r.ParseMultipartForm(32 << 20); err != nil {
				t.Errorf("Failed to parse multipart form: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Simuler une réponse réussie
			response := AudioTranscriptionResponse{
				Transcription: "This is a test transcription",
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

	// Créer un fichier de test
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.mp3")
	if err := os.WriteFile(testFile, []byte("test audio content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test avec différentes options
	opts := &AudioTranscriptionRequest{
		Language: "fr",
		Format:   "text",
	}

	resp, err := client.TranscribeAudioFile(context.Background(), testFile, opts)
	if err != nil {
		t.Fatalf("TranscribeAudioFile failed: %v", err)
	}

	if resp.Transcription != "This is a test transcription" {
		t.Errorf("Unexpected transcription: %s", resp.Transcription)
	}
}
