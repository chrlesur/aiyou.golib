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
	"log"
	"os"
	"strings"
	"time"

	"github.com/chrlesur/aiyou.golib"
)

// truncateText prend un texte et le tronque à environ 30 mots
func truncateText(text string, wordLimit int) string {
	words := strings.Fields(text)
	if len(words) <= wordLimit {
		return text
	}
	return strings.Join(words[:wordLimit], " ") + "..."
}

func main() {
	// Définition des flags pour l'authentification
	email := flag.String("email", "", "Email pour l'authentification (obligatoire)")
	password := flag.String("password", "", "Mot de passe pour l'authentification (obligatoire)")
	debug := flag.Bool("debug", false, "Active les logs de debug")
	baseURL := flag.String("url", "https://ai.dragonflygroup.fr", "URL de base de l'API")
	flag.Parse()

	// Vérification des paramètres obligatoires
	if *email == "" || *password == "" {
		fmt.Println("Les paramètres email et password sont obligatoires")
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

	// Création du client avec gestion des erreurs détaillée
	client, err := aiyou.NewClient(
		*email,
		*password,
		aiyou.WithLogger(logger),
		aiyou.WithBaseURL(*baseURL),
	)
	if err != nil {
		log.Fatalf("Erreur lors de la création du client: %v", err)
	}

	// Context avec timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Récupération des assistants avec gestion des erreurs
	response, err := client.GetUserAssistants(ctx)
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

	// Affichage des assistants
	fmt.Printf("Nombre total d'assistants: %d\n", response.TotalItems)
	for _, assistant := range response.Members {
		fmt.Printf("\n=== Assistant ===\n")
		fmt.Printf("ID: %s\n", assistant.AssistantID)
		fmt.Printf("Nom: %s\n", assistant.Name)
		if assistant.Model != "" {
			fmt.Printf("Modèle: %s\n", assistant.Model)
		}
		if assistant.ModelAi != "" {
			fmt.Printf("Modèle AI: %s\n", assistant.ModelAi)
		}
		if assistant.Instructions != "" {
			fmt.Printf("Instructions: %s\n", truncateText(assistant.Instructions, 30))
		}

		if len(assistant.ThreadHistories) > 0 {
			fmt.Printf("\nDernières conversations:\n")
			for _, history := range assistant.ThreadHistories {
				fmt.Printf("- Thread ID: %s\n", history.ThreadID)
				if history.FirstMessage != "" {
					fmt.Printf(" Premier message: %s\n", truncateText(history.FirstMessage, 20))
				}
			}
		} else {
			fmt.Printf("\nAucune conversation enregistrée\n")
		}
		fmt.Println(strings.Repeat("-", 50))
	}

	// Afficher les métadonnées de la réponse
	if response.Context != "" {
		fmt.Printf("\nContexte: %s\n", response.Context)
	}
	if response.Type != "" {
		fmt.Printf("Type: %s\n", response.Type)
	}
}
