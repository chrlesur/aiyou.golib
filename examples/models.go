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

// printModels affiche les détails des modèles disponibles de façon formatée
func printModels(response *aiyou.AssistantsResponse) {
	modelMap := make(map[string]bool)
	fmt.Printf("\nModèles disponibles (sur %d assistants) :\n", response.TotalItems)

	for _, assistant := range response.Members {
		if assistant.Model != "" && !modelMap[assistant.Model] {
			modelMap[assistant.Model] = true
			fmt.Printf("\nModèle: %s\n", assistant.Model)
		}
		if assistant.ModelAi != "" && !modelMap[assistant.ModelAi] {
			modelMap[assistant.ModelAi] = true
			fmt.Printf("Modèle AI: %s\n", assistant.ModelAi)
		}
	}

	fmt.Printf("\nNombre total de modèles uniques: %d\n", len(modelMap))
	fmt.Println(strings.Repeat("-", 50))
}

func main() {
	// Définition des flags
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

	// Context avec timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	// Récupération des assistants pour obtenir les modèles disponibles
	fmt.Println("Récupération des modèles via les assistants...")
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

	// Affichage des modèles
	if response.TotalItems > 0 {
		printModels(response)
	} else {
		fmt.Println("Aucun assistant trouvé. Impossible de déterminer les modèles disponibles.")
	}
}
