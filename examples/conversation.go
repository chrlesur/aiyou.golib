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
	"log"
	"os"
	"time"

	"github.com/chrlesur/aiyou.golib"
)

const (
	timeoutDuration = 30 * time.Second
	retryDelay      = 5 * time.Second
)

func printMessage(msg aiyou.Message, indent string) {
	fmt.Printf("%sRole: %s\n", indent, msg.Role)
	for _, content := range msg.Content {
		fmt.Printf("%sContenu (%s): %s\n", indent, content.Type, content.Text)
	}
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
		aiyou.WithEmailPassword(*email, *password),
		aiyou.WithLogger(logger),
		aiyou.WithBaseURL(*baseURL),
	)
	if err != nil {
		log.Fatalf("Erreur lors de la création du client: %v", err)
	}

	// Préparation des messages de test
	messages := []aiyou.Message{
		{
			Role: "user",
			Content: []aiyou.ContentPart{
				{Type: "text", Text: "Bonjour, pouvez-vous m'aider ?"},
			},
		},
		{
			Role: "assistant",
			Content: []aiyou.ContentPart{
				{Type: "text", Text: "Bien sûr, comment puis-je vous aider ?"},
			},
		},
	}

	messagesJSON, err := json.Marshal(messages)
	if err != nil {
		log.Fatalf("Erreur lors de la sérialisation des messages: %v", err)
	}

	conversation := aiyou.SaveConversationRequest{
		AssistantID:    *assistantID,
		Conversation:   messages[0].Content[0].Text,
		FirstMessage:   messages[0].Content[0].Text,
		ContentJson:    string(messagesJSON),
		ModelName:      "gpt-4",
		IsNewAppThread: true,
	}

	if *debug {
		fmt.Printf("\nRequête de sauvegarde:\n")
		fmt.Printf("AssistantID: %s\n", conversation.AssistantID)
		fmt.Printf("Conversation: %s\n", conversation.Conversation)
		fmt.Printf("FirstMessage: %s\n", conversation.FirstMessage)
		fmt.Printf("ContentJson: %s\n", conversation.ContentJson)
		fmt.Printf("ModelName: %s\n", conversation.ModelName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	fmt.Println("Sauvegarde de la conversation...")
	saveResp, err := client.SaveConversation(ctx, conversation)
	if err != nil {
		switch e := err.(type) {
		case *aiyou.AuthenticationError:
			log.Fatalf("Erreur d'authentification: %v", e)
		case *aiyou.NetworkError:
			log.Fatalf("Erreur réseau: %v", e)
		case *aiyou.APIError:
			log.Fatalf("Erreur API (code %d): %v", e.StatusCode, e)
		default:
			log.Fatalf("Erreur inattendue: %v", err)
		}
	}

	if *debug {
		fmt.Printf("ID de conversation sauvegardée: %s\n", saveResp.ID)
	}

	fmt.Printf("\nConversation sauvegardée, recherche du thread...\n")
	fmt.Printf("Attente de %s pour la propagation...\n", retryDelay)
	time.Sleep(retryDelay)

	threadsOutput, err := client.GetUserThreads(ctx, &aiyou.UserThreadsParams{
		Page:         1,
		ItemsPerPage: 10,
	})
	if err != nil {
		log.Fatalf("Erreur lors de la récupération des threads: %v", err)
	}

	if *debug {
		fmt.Printf("\nThreads trouvés (%d):\n", len(threadsOutput.Threads))
		for _, t := range threadsOutput.Threads {
			fmt.Printf("\n- Thread ID: %s\n", t.ID)
			fmt.Printf("  Content: %q\n", t.Content)
			fmt.Printf("  Assistant: %s\n", t.AssistantIdOpenAi)
			fmt.Printf("  First Message: %q\n", t.FirstMessage)
			if t.AssistantContentJson != "" {
				fmt.Printf("  Content JSON: %s\n", t.AssistantContentJson)
			}
		}
	}

	var thread *aiyou.ConversationThread
	for _, t := range threadsOutput.Threads {
		if t.AssistantIdOpenAi == conversation.AssistantID &&
			t.FirstMessage == conversation.FirstMessage {
			thread = &t
			break
		}
	}

	if thread == nil {
		log.Fatalf("Thread nouvellement créé non trouvé")
	}

	fmt.Printf("\nConversation récupérée :\n")
	fmt.Printf("Thread ID: %s\n", thread.ID)
	fmt.Printf("Thread Param ID: %d\n", thread.ThreadIdParam)
	fmt.Printf("Contenu: %s\n", thread.Content)
	fmt.Printf("Assistant Name: %s\n", thread.AssistantName)
	if thread.AssistantModel != nil {
		fmt.Printf("Assistant Model: %s\n", *thread.AssistantModel)
	}
	fmt.Printf("Assistant ID: %d\n", thread.AssistantId)
	fmt.Printf("Assistant OpenAI ID: %s\n", thread.AssistantIdOpenAi)
	fmt.Printf("Created at: %s\n", thread.CreatedAt.Format(time.RFC3339))
	fmt.Printf("Updated at: %s\n", thread.UpdatedAt.Format(time.RFC3339))
	fmt.Printf("Is New App Thread: %v\n", thread.IsNewAppThread)

	// Affichage des messages
	if thread.AssistantContentJson != "" {
		var messages []aiyou.Message
		if err := json.Unmarshal([]byte(thread.AssistantContentJson), &messages); err != nil {
			fmt.Printf("\nErreur lors du décodage des messages: %v\n", err)
		} else {
			fmt.Printf("\nMessages (%d):\n", len(messages))
			for i, msg := range messages {
				fmt.Printf("\nMessage %d:\n", i+1)
				printMessage(msg, " ")
			}
		}
	} else {
		fmt.Println("\nAucun message dans la conversation")
	}

	fmt.Printf("\nConversation récupérée :\n")
	fmt.Printf("Thread ID: %s\n", thread.ID)
	fmt.Printf("Thread Param ID: %d\n", thread.ThreadIdParam)
	fmt.Printf("Contenu: %s\n", thread.Content)
	fmt.Printf("Assistant Name: %s\n", thread.AssistantName)
	if thread.AssistantModel != nil {
		fmt.Printf("Assistant Model: %s\n", *thread.AssistantModel)
	}
	fmt.Printf("Assistant ID: %d\n", thread.AssistantId)
	fmt.Printf("Assistant OpenAI ID: %s\n", thread.AssistantIdOpenAi)
	fmt.Printf("Created at: %s\n", thread.CreatedAt.Format(time.RFC3339))
	fmt.Printf("Updated at: %s\n", thread.UpdatedAt.Format(time.RFC3339))
	fmt.Printf("Is New App Thread: %v\n", thread.IsNewAppThread)

	// Affichage des messages
	if thread.AssistantContentJson != "" {
		var messages []aiyou.Message
		if err := json.Unmarshal([]byte(thread.AssistantContentJson), &messages); err != nil {
			fmt.Printf("\nErreur lors du décodage des messages: %v\n", err)
		} else {
			fmt.Printf("\nMessages (%d):\n", len(messages))
			for i, msg := range messages {
				fmt.Printf("\nMessage %d:\n", i+1)
				printMessage(msg, " ")
			}
		}
	} else {
		fmt.Println("\nAucun message dans la conversation")
	}
}
