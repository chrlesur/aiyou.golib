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

const (
	timeoutDuration = 30 * time.Second
)

func main() {
	// Définition des flags
	email := flag.String("email", "", "Email pour l'authentification (obligatoire)")
	password := flag.String("password", "", "Mot de passe pour l'authentification (obligatoire)")
	debug := flag.Bool("debug", false, "Active les logs de debug")
	baseURL := flag.String("url", "https://ai.dragonflygroup.fr", "URL de base de l'API")
	page := flag.Int("page", 1, "Numéro de page pour la pagination")
	itemsPerPage := flag.Int("items", 10, "Nombre d'éléments par page")
	search := flag.String("search", "", "Terme de recherche optionnel")
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

	// Création du client
	client, err := aiyou.NewClient(
		aiyou.WithEmailPassword(*email, *password),
		aiyou.WithLogger(logger),
		aiyou.WithBaseURL(*baseURL),
	)
	if err != nil {
		log.Fatalf("Erreur lors de la création du client: %v", err)
	}

	// Context avec timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	// Paramètres de recherche des threads
	params := &aiyou.UserThreadsParams{
		Page:         *page,
		ItemsPerPage: *itemsPerPage,
		Search:       *search,
	}

	// Récupération des threads
	fmt.Printf("Récupération des threads (page %d, %d éléments par page)...\n", params.Page, params.ItemsPerPage)
	threadsOutput, err := client.GetUserThreads(ctx, params)
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

	// Affichage des résultats
	fmt.Printf("\nTotal des threads: %d\n", threadsOutput.TotalItems)
	fmt.Printf("Page actuelle: %d/%d\n\n",
		threadsOutput.CurrentPage,
		(threadsOutput.TotalItems+threadsOutput.ItemsPerPage-1)/threadsOutput.ItemsPerPage,
	)

	for _, thread := range threadsOutput.Threads {
		fmt.Printf("=== Thread ===\n")
		fmt.Printf("ID: %s\n", thread.ID)
		fmt.Printf("Contenu: %s\n", thread.Content)
		fmt.Printf("Assistant: %s\n", thread.AssistantName)
		if thread.AssistantModel != nil {
			fmt.Printf("Modèle: %s\n", *thread.AssistantModel)
		}
		fmt.Printf("Assistant OpenAI ID: %s\n", thread.AssistantIdOpenAi)
		fmt.Printf("Premier message: %s\n", thread.FirstMessage)
		fmt.Printf("Créé le: %s\n", thread.CreatedAt.Format(time.RFC3339))
		fmt.Printf("Mis à jour le: %s\n", thread.UpdatedAt.Format(time.RFC3339))

		if thread.AssistantContentJson != "" {
			fmt.Printf("\nContenu JSON:\n%s\n", thread.AssistantContentJson)
		}
		fmt.Println(strings.Repeat("-", 50))
	}

	// Afficher les méta-informations de pagination
	fmt.Printf("\nInformations de pagination:\n")
	fmt.Printf("Page actuelle: %d\n", threadsOutput.CurrentPage)
	fmt.Printf("Éléments par page: %d\n", threadsOutput.ItemsPerPage)
	fmt.Printf("Total des éléments: %d\n", threadsOutput.TotalItems)
}
