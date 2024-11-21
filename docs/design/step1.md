Contexte :
Vous allez commencer le développement d'un package Go nommé "aiyou" qui implémentera l'accès à l'API AI.YOU de Cloud Temple. Ce package permettra une interaction fluide et typée avec tous les endpoints de l'API, y compris l'authentification, les chat completions, la gestion des assistants et des threads utilisateur, et la transcription audio.

Tâche :
Configurez l'environnement de développement initial pour le projet "aiyou".

Instructions détaillées :

1. Création du répertoire et structure :
   - Créez un nouveau répertoire nommé "aiyou".
   - À l'intérieur, créez les sous-répertoires suivants : "cmd", "internal", "pkg".

2. Initialisation du module Go :
   - Ouvrez un terminal dans le répertoire "aiyou".
   - Exécutez la commande : `go mod init github.com/chrlesur/aiyou`

3. Configuration de Git :
   - Initialisez un nouveau dépôt Git : `git init`
   - Créez un fichier .gitignore avec le contenu suivant :
     ```
     # Binaires compilés et fichiers temporaires
     *.exe
     *.exe~
     *.dll
     *.so
     *.dylib
     *.test
     *.out

     # Dépendances
     /vendor/

     # Fichiers de configuration d'IDE
     .idea/
     .vscode/

     # Fichiers de log
     *.log
     ```

4. Création des fichiers de base :
   - Dans le répertoire "pkg/aiyou", créez les fichiers suivants : 
     - client.go
     - auth.go
     - types.go
     - errors.go

5. Documentation et commentaires :
   - Dans chaque fichier, ajoutez un en-tête de licence GPL-3.0.
   - Ajoutez des commentaires de package au début de chaque fichier.

6. Journalisation :
   - Dans client.go, préparez la structure pour intégrer un logger configurable.

7. Commit initial :
   - Ajoutez tous les fichiers créés : `git add .`
   - Faites un commit initial : `git commit -m "Initial project setup"`

Exigences supplémentaires :
- Assurez-vous que tous les commentaires et la documentation sont en anglais.
- Suivez les conventions de nommage standard de Go.
- Préparez la structure pour une future intégration de tests unitaires.

Exemple de contenu pour client.go :

```go
/*
Copyright (C) 2023 Cloud Temple

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

// Package aiyou provides a client for interacting with the AI.YOU API from Cloud Temple.
package aiyou

import (
	"log"
	"net/http"
)

// Client represents an AI.YOU API client.
type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     *log.Logger
	// Add other necessary fields
}

// NewClient creates a new AI.YOU API client.
func NewClient(baseURL string, options ...ClientOption) (*Client, error) {
	// Implementation will be added later
	return nil, nil
}

// ClientOption allows setting custom parameters to the client.
type ClientOption func(*Client) error

// WithLogger sets a custom logger for the client.
func WithLogger(logger *log.Logger) ClientOption {
	return func(c *Client) error {
		c.logger = logger
		return nil
	}
}

// Add other necessary types and functions
```

Résultat attendu :
Une structure de projet initiale propre et bien organisée, prête pour le développement futur du package aiyou.