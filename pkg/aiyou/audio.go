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

// File: pkg/aiyou/audio.go
package aiyou

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// Formats audio supportés
var SupportedFormats = []SupportedAudioFormat{
	{Extension: ".mp3", MimeTypes: []string{"audio/mpeg"}, MaxSize: 25 * 1024 * 1024},
	{Extension: ".wav", MimeTypes: []string{"audio/wav", "audio/x-wav"}, MaxSize: 25 * 1024 * 1024},
	{Extension: ".m4a", MimeTypes: []string{"audio/mp4", "audio/x-m4a"}, MaxSize: 25 * 1024 * 1024},
}

// TranscribeAudioFile transcrit un fichier audio en texte
func (c *Client) TranscribeAudioFile(ctx context.Context, filePath string, opts *AudioTranscriptionRequest) (*AudioTranscriptionResponse, error) {
	c.logger.Debugf("Starting audio transcription for file: %s", filePath)

	// Ouvrir et vérifier le fichier
	file, err := os.Open(filePath)
	if err != nil {
		c.logger.Errorf("Failed to open audio file: %v", err)
		return nil, fmt.Errorf("failed to open audio file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		c.logger.Errorf("Failed to get file info: %v", err)
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Créer un pipe pour lire le body
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	// Goroutine pour écrire le multipart form
	go func() {
		defer pw.Close()
		var writeError error

		// Créer la partie fichier
		part, err := writer.CreateFormFile("audioFile", filepath.Base(filePath))
		if err != nil {
			writeError = err
			pw.CloseWithError(err)
			return
		}

		// Copier le fichier
		_, err = io.Copy(part, file)
		if err != nil {
			writeError = err
			pw.CloseWithError(err)
			return
		}

		// Ajouter les champs optionnels
		if opts != nil {
			if opts.Language != "" {
				if err := writer.WriteField("language", opts.Language); err != nil {
					writeError = err
					pw.CloseWithError(err)
					return
				}
			}
			if opts.Format != "" {
				if err := writer.WriteField("format", opts.Format); err != nil {
					writeError = err
					pw.CloseWithError(err)
					return
				}
			}
		}

		// Fermer le writer
		if err := writer.Close(); err != nil {
			writeError = err
			pw.CloseWithError(err)
			return
		}

		if writeError != nil {
			c.logger.Errorf("Error writing multipart form: %v", writeError)
		}
	}()

	// Créer la requête
	endpoint := "/api/v1/audio/transcriptions"
	c.logger.Debugf("Creating request to %s with file size: %d bytes", endpoint, fileInfo.Size())

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+endpoint, pr)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Authentification
	if err := c.auth.Authenticate(ctx); err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.auth.Token())

	// Content-Type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Log des headers
	c.logger.Debugf("Request headers:")
	for name, values := range req.Header {
		c.logger.Debugf(" %s: %v", name, values)
	}

	// Envoyer la requête
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Lire le corps de la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Errorf("Transcription failed with status %d: %s", resp.StatusCode, string(body))
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("transcription failed: %s", string(body)),
		}
	}

	// Décoder la réponse
	var transcription AudioTranscriptionResponse
	if err := json.Unmarshal(body, &transcription); err != nil {
		c.logger.Errorf("Failed to decode response: %v. Body: %s", err, string(body))
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	c.logger.Debugf("Successfully transcribed audio file: %s", filePath)
	return &transcription, nil
}

// validateAudioFile vérifie si le fichier audio est dans un format supporté
func validateAudioFile(file *os.File, filePath string) error {
	// Obtenir les informations du fichier
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Vérifier l'extension
	ext := filepath.Ext(filePath)
	var supportedFormat *SupportedAudioFormat
	for _, format := range SupportedFormats {
		if format.Extension == ext {
			supportedFormat = &format
			break
		}
	}

	if supportedFormat == nil {
		return fmt.Errorf("unsupported audio format: %s", ext)
	}

	// Vérifier la taille
	if fileInfo.Size() > supportedFormat.MaxSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", supportedFormat.MaxSize)
	}

	return nil
}
