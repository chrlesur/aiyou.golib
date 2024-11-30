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
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chrlesur/aiyou.golib"
)

const (
	timeoutDuration = 30 * time.Second
)

func printMessage(msg aiyou.Message, indent string) {
	fmt.Printf("%sRole: %s\n", indent, msg.Role)
	for _, content := range msg.Content {
		fmt.Printf("%sContenu (%s): %s\n", indent, content.Type, content.Text)
	}
	fmt.Println(strings.Repeat("-", 50))
}

func demonstrateNonStreamingMode(client *aiyou.Client, assistantID string, logger aiyou.Logger) error {
	ctx := context.Background()

	builder := aiyou.NewMessageBuilder("user", logger)
	message := builder.
		AddText("Qu'est-ce que le SecNumCloud ?").
		Build()

	fmt.Println("Message à envoyer :")
	printMessage(message, " ")

	req := aiyou.ChatCompletionRequest{
		Messages:     []aiyou.Message{message},
		AssistantID:  assistantID,
		Temperature:  1,
		TopP:         1.0,
		Stream:       false,
		PromptSystem: "", // Explicitement vide
	}

	// Log la requête pour débugger
	reqJSON, _ := json.MarshalIndent(req, "", " ")
	fmt.Printf("\nRequête JSON:\n%s\n", string(reqJSON))

	fmt.Println("\nEnvoi du message (mode non-streaming)...")
	resp, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return fmt.Errorf("erreur mode non-streaming: %v", err)
	}

	fmt.Println("\nRéponse de l'assistant (non-streaming):")
	if len(resp.Choices) > 0 {
		printMessage(resp.Choices[0].Message, " ")
	} else {
		fmt.Println("Aucune réponse reçue")
	}

	return nil
}

func demonstrateStreamingMode(client *aiyou.Client, assistantID string, logger aiyou.Logger) error {
	ctx := context.Background()

	builder := aiyou.NewMessageBuilder("user", logger)
	message := builder.
		AddText("Comment fonctionne l'authentification dans le cloud ?").
		Build()

	fmt.Println("Message à envoyer :")
	printMessage(message, " ")

	req := aiyou.ChatCompletionRequest{
		Messages:    []aiyou.Message{message},
		AssistantID: assistantID,
		Temperature: 1,
		TopP:        1.0,
		Stream:      true,
	}

	stream, err := client.ChatCompletionStream(ctx, req)
	if err != nil {
		return fmt.Errorf("erreur création stream: %v", err)
	}
	defer stream.Close()

	fmt.Println("\nRéponse de l'assistant (streaming):")
	for {
		chunk, err := stream.ReadChunk()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Errorf("Error reading chunk: %v", err)
			return fmt.Errorf("error reading stream: %w", err)
		}

		if chunk == nil || len(chunk.Choices) == 0 {
			continue
		}

		choice := chunk.Choices[0]
		if choice.Delta != nil && choice.Delta.Content != "" {
			fmt.Print(choice.Delta.Content)
			os.Stdout.Sync()
		}
	}
	fmt.Println("\n" + strings.Repeat("-", 50))
	return nil
}

func main() {
	// Définition des flags
	email := flag.String("email", "", "Email pour l'authentification (obligatoire)")
	password := flag.String("password", "", "Mot de passe pour l'authentification (obligatoire)")
	assistantID := flag.String("assistant", "", "ID de l'assistant (obligatoire)")
	debug := flag.Bool("debug", false, "Active les logs de debug")
	baseURL := flag.String("url", "https://ai.dragonflygroup.fr", "URL de base de l'API")
	flag.Parse()

	// Vérification des paramètres obligatoires
	if *email == "" || *password == "" || *assistantID == "" {
		fmt.Println("Les paramètres email, password et assistant sont obligatoires")
		flag.Usage()
		os.Exit(1)
	}

	// Configuration du logger
	logger := aiyou.NewDefaultLogger(os.Stdout)
	if *debug {
		logger.SetLevel(aiyou.DEBUG)
	} else {
		logger.SetLevel(aiyou.INFO)
	}

	// Création du client
	client, err := aiyou.NewClient(
		*email,
		*password,
		aiyou.WithLogger(logger),
		aiyou.WithBaseURL(*baseURL),
	)
	if err != nil {
		log.Fatalf("Erreur lors de la création du client: %v", err)
	}

	// Démonstration des deux modes
	fmt.Println("Démonstration des modes streaming et non-streaming")

	// Test du mode non-streaming
	if err := demonstrateNonStreamingMode(client, *assistantID, logger); err != nil {
		log.Printf("Erreur en mode non-streaming: %v", err)
	}

	// Test du mode streaming
	if err := demonstrateStreamingMode(client, *assistantID, logger); err != nil {
		log.Printf("Erreur en mode streaming: %v", err)
	}
}
