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

// File: examples/audio.go
package main

import (
 "context"
 "fmt"
 "log"

 "github.com/chrlesur/aiyou.golib"
)

func main() {
 client, err := aiyou.NewClient("your-email@example.com", "your-password")
 if err != nil {
 log.Fatalf("Error creating client: %v", err)
 }

 // Options de transcription
 opts := &aiyou.AudioTranscriptionRequest{
 Language: "fr", // Langue source (optionnel)
 Format: "text", // Format de sortie
 }

 // Transcrire un fichier audio
 filePath := "path/to/your/audio/file.mp3"
 transcription, err := client.TranscribeAudioFile(context.Background(), filePath, opts)
 if err != nil {
 log.Fatalf("Error transcribing audio: %v", err)
 }

 // Afficher les résultats
 fmt.Printf("Transcription completed!\n")
 fmt.Printf("Text: %s\n", transcription.Text)
 fmt.Printf("Language detected: %s\n", transcription.Language)
 fmt.Printf("Duration: %.2f seconds\n", transcription.Duration)
 fmt.Printf("Created at: %s\n", transcription.CreatedAt)
 fmt.Printf("Status: %s\n", transcription.Status)
}