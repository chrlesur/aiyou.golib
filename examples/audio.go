/*
Copyright (C) 2024 Cloud Temple

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chrlesur/aiyou.golib"
	"github.com/fatih/color"
)

var (
	email     string
	password  string
	baseURL   string
	debug     bool
	quietMode bool
	audioFile string
	language  string
	format    string

	// Couleurs pour l'interface
	successColor = color.New(color.FgGreen)
	errorColor   = color.New(color.FgRed)
	infoColor    = color.New(color.FgYellow)
	headerColor  = color.New(color.FgCyan, color.Bold)
)

func init() {
	flag.StringVar(&email, "email", "", "Email pour l'authentification (obligatoire)")
	flag.StringVar(&password, "password", "", "Mot de passe pour l'authentification (obligatoire)")
	flag.StringVar(&baseURL, "url", "https://ai.dragonflygroup.fr", "URL de base de l'API")
	flag.BoolVar(&debug, "debug", false, "Active les logs de debug")
	flag.BoolVar(&quietMode, "quiet", false, "Désactive les messages de statut")
	flag.StringVar(&audioFile, "file", "", "Fichier audio à transcrire (obligatoire)")
	flag.StringVar(&language, "lang", "fr", "Langue source (fr, en, ...)")
	flag.StringVar(&format, "format", "text", "Format de sortie (text, json, ...)")
}

func createClient() (*aiyou.Client, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("email et mot de passe requis")
	}

	logger := aiyou.NewDefaultLogger(os.Stderr)
	if debug {
		logger.SetLevel(aiyou.DEBUG)
	} else if quietMode {
		logger.SetLevel(aiyou.ERROR)
	} else {
		logger.SetLevel(aiyou.INFO)
	}

	// Utilisation de WithEmailPassword au lieu des paramètres directs
	return aiyou.NewClient(
		aiyou.WithEmailPassword(email, password),
		aiyou.WithLogger(logger),
		aiyou.WithBaseURL(baseURL),
	)
}

func validateAudioFile(filePath string) (os.FileInfo, error) {
	if filePath == "" {
		return nil, fmt.Errorf("chemin du fichier audio requis")
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("le fichier n'existe pas: %s", filePath)
		}
		return nil, fmt.Errorf("erreur lors de l'accès au fichier: %v", err)
	}

	if fileInfo.IsDir() {
		return nil, fmt.Errorf("le chemin spécifié est un dossier: %s", filePath)
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	supportedExts := map[string]bool{
		".mp3": true,
		".wav": true,
		".m4a": true,
	}

	if !supportedExts[ext] {
		return nil, fmt.Errorf("format de fichier non supporté: %s (formats acceptés: mp3, wav, m4a)", ext)
	}

	// Limite de taille (25MB)
	const maxSize = 25 * 1024 * 1024
	if fileInfo.Size() > maxSize {
		return nil, fmt.Errorf("fichier trop volumineux: %.2fMB (maximum: 25MB)", float64(fileInfo.Size())/(1024*1024))
	}

	return fileInfo, nil
}

func showSpinner(done chan bool) {
	if quietMode {
		return
	}

	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for {
		select {
		case <-done:
			fmt.Fprintf(os.Stderr, "\r\033[K") // Efface la ligne
			return
		default:
			fmt.Fprintf(os.Stderr, "\r%s Transcription en cours...", spinner[i])
			i = (i + 1) % len(spinner)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func extractErrorDetails(err error) string {
	if apiErr, ok := err.(*aiyou.APIError); ok {
		return fmt.Sprintf("Erreur API (%d): %s", apiErr.StatusCode, apiErr.Message)
	}
	return err.Error()
}

func main() {
	flag.Parse()

	if email == "" || password == "" || audioFile == "" {
		fmt.Fprintln(os.Stderr, "Les paramètres email, password et file sont obligatoires")
		flag.Usage()
		os.Exit(1)
	}

	// Validation du fichier audio
	fileInfo, err := validateAudioFile(audioFile)
	if err != nil {
		errorColor.Fprintf(os.Stderr, "Erreur de validation: %v\n", err)
		os.Exit(1)
	}

	absPath, err := filepath.Abs(audioFile)
	if err != nil {
		errorColor.Fprintf(os.Stderr, "Erreur lors de la résolution du chemin: %v\n", err)
		os.Exit(1)
	}

	client, err := createClient()
	if err != nil {
		errorColor.Fprintf(os.Stderr, "Erreur lors de la création du client: %v\n", err)
		os.Exit(1)
	}

	if !quietMode {
		headerColor.Fprintf(os.Stderr, "\nTranscription audio\n")
		infoColor.Fprintf(os.Stderr, "Fichier: %s\n", absPath)
		infoColor.Fprintf(os.Stderr, "Taille: %.2f MB\n", float64(fileInfo.Size())/(1024*1024))
		infoColor.Fprintf(os.Stderr, "Format: %s\n", filepath.Ext(audioFile))
		infoColor.Fprintf(os.Stderr, "Langue: %s\n", language)
		infoColor.Fprintf(os.Stderr, "Format de sortie: %s\n\n", format)
	}
	// Options de transcription
	opts := &aiyou.AudioTranscriptionRequest{
		Language: language,
		Format:   format,
		FileName: filepath.Base(audioFile),
	}

	if debug {
		infoColor.Fprintf(os.Stderr, "Envoi de la requête avec les options:\n")
		infoColor.Fprintf(os.Stderr, "- Nom du fichier: %s\n", opts.FileName)
		infoColor.Fprintf(os.Stderr, "- Langue: %s\n", opts.Language)
		infoColor.Fprintf(os.Stderr, "- Format: %s\n", opts.Format)
	}

	// Démarrer le spinner
	spinnerDone := make(chan bool)
	go showSpinner(spinnerDone)

	// Transcription avec timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	response, err := client.TranscribeAudioFile(ctx, absPath, opts)
	close(spinnerDone) // Arrêter le spinner

	if err != nil {
		fmt.Fprintln(os.Stderr) // Nouvelle ligne après le spinner
		errorColor.Fprintf(os.Stderr, "Erreur lors de la transcription: %s\n", extractErrorDetails(err))
		if debug {
			// Afficher plus d'informations sur le fichier en cas d'erreur
			errorColor.Fprintf(os.Stderr, "\nDétails du fichier:\n")
			errorColor.Fprintf(os.Stderr, "- Chemin complet: %s\n", absPath)
			errorColor.Fprintf(os.Stderr, "- Taille: %d bytes\n", fileInfo.Size())
			errorColor.Fprintf(os.Stderr, "- Mode: %v\n", fileInfo.Mode())
			errorColor.Fprintf(os.Stderr, "- Modifié le: %v\n", fileInfo.ModTime())

			// Vérifier si le fichier est lisible
			if file, err := os.Open(absPath); err != nil {
				errorColor.Fprintf(os.Stderr, "Impossible d'ouvrir le fichier: %v\n", err)
			} else {
				file.Close()
				errorColor.Fprintf(os.Stderr, "Le fichier est accessible en lecture\n")
			}
		}
		os.Exit(1)
	}

	if !quietMode {
		successColor.Fprintln(os.Stderr, "\nTranscription réussie!")
		fmt.Fprintln(os.Stderr, strings.Repeat("-", 50))
	}

	// Afficher la transcription
	fmt.Println(response.Transcription)

	if !quietMode {
		fmt.Fprintln(os.Stderr, strings.Repeat("-", 50))
	}
}
