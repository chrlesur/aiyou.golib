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
	"bytes"
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
   
	// Vérifier l'existence du fichier
	file, err := os.Open(filePath)
	if err != nil {
	c.logger.Errorf("Failed to open audio file: %v", err)
	return nil, fmt.Errorf("failed to open audio file: %w", err)
	}
	defer file.Close()
   
	// Vérifier le format et la taille du fichier
	if err := validateAudioFile(file, filePath); err != nil {
	c.logger.Errorf("Invalid audio file: %v", err)
	return nil, err
	}
   
	// Retourner au début du fichier après validation
	if _, err := file.Seek(0, 0); err != nil {
	return nil, fmt.Errorf("failed to reset file position: %w", err)
	}
   
	// Créer un buffer pour construire la requête multipart
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
   
	// Ajouter les options de transcription
	if opts != nil {
	optionsJson, err := json.Marshal(opts)
	if err == nil {
	if err := writer.WriteField("options", string(optionsJson)); err != nil {
	return nil, fmt.Errorf("failed to write options field: %w", err)
	}
	}
	}
   
	// Ajouter le fichier audio
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
	return nil, fmt.Errorf("failed to create form file: %w", err)
	}
   
	if _, err := io.Copy(part, file); err != nil {
	return nil, fmt.Errorf("failed to copy file content: %w", err)
	}
   
	// Fermer le writer multipart
	if err := writer.Close(); err != nil {
	return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}
   
	// Créer la requête HTTP
	endpoint := "/api/v1/audio/transcriptions"
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+endpoint, &requestBody)
	if err != nil {
	return nil, fmt.Errorf("failed to create request: %w", err)
	}
   
	// Définir le Content-Type avec la boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())
   
	// Ajouter le token d'authentification
	if err := c.auth.Authenticate(ctx); err != nil {
	return nil, fmt.Errorf("failed to authenticate: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.auth.Token())
   
	// Envoyer la requête
	resp, err := c.httpClient.Do(req)
	if err != nil {
	return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
   
	if resp.StatusCode != http.StatusOK {
	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
   
	// Décoder la réponse
	var transcription AudioTranscriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&transcription); err != nil {
	return nil, fmt.Errorf("failed to decode transcription response: %w", err)
	}
   
	c.logger.Infof("Successfully transcribed audio file: %s", filePath)
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
