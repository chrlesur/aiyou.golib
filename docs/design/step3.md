Contexte du projet :
Vous travaillez sur le package Go "aiyou.golib", qui est une interface pour l'API AI.YOU de Cloud Temple. La structure de base du package a été mise en place, et nous devons maintenant implémenter l'authentification.

Important : L'API AI.YOU utilise une authentification basée sur un email et un mot de passe, qui génère un token JWT. Il n'y a pas d'authentification par clé API.

Objectif de cette étape :
Implémenter le mécanisme d'authentification pour le client API, en utilisant l'interface Authenticator déjà définie, mais en l'adaptant pour utiliser email/mot de passe et gérer les tokens JWT.

État actuel du projet :
[
  {
    "name": "simple_client.go",
    "content": "/*\r\nCopyright (C) 2023 Cloud Temple\r\n\r\nThis program is free software: you can redistribute it and/or modify\r\nit under the terms of the GNU General Public License as published by\r\nthe Free Software Foundation, either version 3 of the License, or\r\n(at your option) any later version.\r\n\r\nThis program is distributed in the hope that it will be useful,\r\nbut WITHOUT ANY WARRANTY; without even the implied warranty of\r\nMERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\r\nGNU General Public License for more details.\r\n\r\nYou should have received a copy of the GNU General Public License\r\nalong with this program.  If not, see \u003chttps://www.gnu.org/licenses/\u003e.\r\n*/\r\n\r\n\r\npackage main\r\n\r\nimport (\r\n    \"context\"\r\n    \"fmt\"\r\n    \"log\"\r\n\r\n    \"github.com/chrlesur/aiyou.golib\"\r\n)\r\n\r\nfunc main() {\r\n    client, err := aiyou.NewClient(aiyou.WithAPIKey(\"votre-api-key\"))\r\n    if err != nil {\r\n        log.Fatalf(\"Erreur lors de la création du client: %v\", err)\r\n    }\r\n\r\n    ctx := context.Background()\r\n\r\n    // Exemple d'utilisation d'une méthode authentifiée\r\n    response, err := client.SomeAuthenticatedMethod(ctx, \"paramètre\")\r\n    if err != nil {\r\n        log.Fatalf(\"Erreur lors de l'appel de la méthode authentifiée: %v\", err)\r\n    }\r\n\r\n    fmt.Printf(\"Réponse: %+v\\n\", response)\r\n}",
    "size": 1293,
    "modTime": "2024-11-21T14:53:03.3778979+01:00",
    "path": "examples\\simple_client.go"
  },
  {
    "name": "auth.go",
    "content": "/*\nCopyright (C) 2023 Cloud Temple\n\nThis program is free software: you can redistribute it and/or modify\nit under the terms of the GNU General Public License as published by\nthe Free Software Foundation, either version 3 of the License, or\n(at your option) any later version.\n\nThis program is distributed in the hope that it will be useful,\nbut WITHOUT ANY WARRANTY; without even the implied warranty of\nMERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\nGNU General Public License for more details.\n\nYou should have received a copy of the GNU General Public License\nalong with this program.  If not, see \u003chttps://www.gnu.org/licenses/\u003e.\n*/\n\n// Package aiyou provides authentication functionalities for the AI.YOU API.\npackage aiyou\n\n// Authentication related types and functions will be implemented here\n",
    "size": 816,
    "modTime": "2024-11-21T15:00:40.9349884+01:00",
    "path": "pkg\\aiyou\\auth.go"
  },
  {
    "name": "auth_test.go",
    "content": "",
    "size": 0,
    "modTime": "2024-11-21T14:57:40.5590889+01:00",
    "path": "pkg\\aiyou\\auth_test.go"
  },
  {
    "name": "client.go",
    "content": "/*\nCopyright (C) 2023 Cloud Temple\n\nThis program is free software: you can redistribute it and/or modify\nit under the terms of the GNU General Public License as published by\nthe Free Software Foundation, either version 3 of the License, or\n(at your option) any later version.\n\nThis program is distributed in the hope that it will be useful,\nbut WITHOUT ANY WARRANTY; without even the implied warranty of\nMERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\nGNU General Public License for more details.\n\nYou should have received a copy of the GNU General Public License\nalong with this program.  If not, see \u003chttps://www.gnu.org/licenses/\u003e.\n*/\n\n// Package aiyou provides a client for interacting with the AI.YOU API from Cloud Temple.\npackage aiyou\n\nimport (\n    \"context\"\n    \"net/http\"\n    \"log\"\n)\n\n// Client représente un client pour l'API AI.YOU.\ntype Client struct {\n    baseURL    string\n    httpClient *http.Client\n    logger     *log.Logger\n    auth       Authenticator\n    config     *Config\n}\n\n// NewClient crée une nouvelle instance de Client avec la configuration donnée.\nfunc NewClient(config *Config, options ...ClientOption) (*Client, error) {\n    if err := config.Validate(); err != nil {\n        return nil, err\n    }\n\n    client := \u0026Client{\n        baseURL:    config.BaseURL,\n        httpClient: \u0026http.Client{Timeout: config.Timeout},\n        logger:     log.New(log.Writer(), \"aiyou: \", log.LstdFlags),\n        config:     config,\n    }\n\n    for _, option := range options {\n        option(client)\n    }\n\n    return client, nil\n}\n\n// SetBaseURL définit l'URL de base pour les requêtes API.\nfunc (c *Client) SetBaseURL(url string) {\n    c.baseURL = url\n}\n\n// SetLogger définit le logger pour le client.\nfunc (c *Client) SetLogger(logger *log.Logger) {\n    c.logger = logger\n}\n\n// SetAuthenticator définit l'authentificateur pour le client.\nfunc (c *Client) SetAuthenticator(auth Authenticator) {\n    c.auth = auth\n}\n",
    "size": 1950,
    "modTime": "2024-11-21T15:00:36.2423833+01:00",
    "path": "pkg\\aiyou\\client.go"
  },
  {
    "name": "config.go",
    "content": "/*\r\nCopyright (C) 2023 Cloud Temple\r\n\r\nThis program is free software: you can redistribute it and/or modify\r\nit under the terms of the GNU General Public License as published by\r\nthe Free Software Foundation, either version 3 of the License, or\r\n(at your option) any later version.\r\n\r\nThis program is distributed in the hope that it will be useful,\r\nbut WITHOUT ANY WARRANTY; without even the implied warranty of\r\nMERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\r\nGNU General Public License for more details.\r\n\r\nYou should have received a copy of the GNU General Public License\r\nalong with this program.  If not, see \u003chttps://www.gnu.org/licenses/\u003e.\r\n*/\r\npackage aiyou\r\n\r\nimport (\r\n    \"errors\"\r\n    \"time\"\r\n)\r\n\r\n// Config contient les paramètres de configuration pour le client AI.YOU.\r\ntype Config struct {\r\n    BaseURL    string\r\n    APIKey     string\r\n    Timeout    time.Duration\r\n    RetryCount int\r\n}\r\n\r\n// Validate vérifie que la configuration est valide.\r\nfunc (c *Config) Validate() error {\r\n    if c.BaseURL == \"\" {\r\n        return errors.New(\"BaseURL est requis\")\r\n    }\r\n    if c.APIKey == \"\" {\r\n        return errors.New(\"APIKey est requis\")\r\n    }\r\n    if c.Timeout \u003c= 0 {\r\n        c.Timeout = 30 * time.Second // Valeur par défaut\r\n    }\r\n    if c.RetryCount \u003c 0 {\r\n        c.RetryCount = 0\r\n    }\r\n    return nil\r\n}",
    "size": 1348,
    "modTime": "2024-11-21T13:57:50.9309087+01:00",
    "path": "pkg\\aiyou\\config.go"
  },
  {
    "name": "errors.go",
    "content": "",
    "size": 0,
    "modTime": "2024-11-20T10:44:19.6000093+01:00",
    "path": "pkg\\aiyou\\errors.go"
  },
  {
    "name": "types.go",
    "content": "/*\nCopyright (C) 2023 Cloud Temple\n\nThis program is free software: you can redistribute it and/or modify\nit under the terms of the GNU General Public License as published by\nthe Free Software Foundation, either version 3 of the License, or\n(at your option) any later version.\n\nThis program is distributed in the hope that it will be useful,\nbut WITHOUT ANY WARRANTY; without even the implied warranty of\nMERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\nGNU General Public License for more details.\n\nYou should have received a copy of the GNU General Public License\nalong with this program.  If not, see \u003chttps://www.gnu.org/licenses/\u003e.\n*/\n\n// Package aiyou defines types used across the AI.YOU API client.\npackage aiyou\n\nimport \"context\"\n\n// Authenticator gère l'authentification pour les requêtes API.\ntype Authenticator interface {\n    Authenticate(ctx context.Context) error\n    Token() string\n}\n\n// ChatCompleter gère les requêtes de complétion de chat.\ntype ChatCompleter interface {\n    Complete(ctx context.Context, input ChatCompletionInput) (*ChatCompletionOutput, error)\n}\n\n// ChatCompletionInput représente l'entrée pour une requête de complétion de chat.\ntype ChatCompletionInput struct {\n    Messages []Message\n    // Autres champs à ajouter selon la spécification de l'API\n}\n\n// ChatCompletionOutput représente la sortie d'une requête de complétion de chat.\ntype ChatCompletionOutput struct {\n    Response string\n    // Autres champs à ajouter selon la spécification de l'API\n}\n\n// Message représente un message dans une conversation.\ntype Message struct {\n    Role    string `json:\"role\"`\n    Content string `json:\"content\"`\n}\n",
    "size": 1672,
    "modTime": "2024-11-21T14:57:22.2419642+01:00",
    "path": "pkg\\aiyou\\types.go"
  }
]
1. Dans auth.go :
   - Implémentez une structure JWTAuthenticator qui satisfait l'interface Authenticator :
     ```go
     type JWTAuthenticator struct {
         email    string
         password string
         token    string
         expiry   time.Time
         client   *http.Client
     }
     ```
   - Implémentez les méthodes Authenticate et Token pour JWTAuthenticator.
   - Ajoutez une fonction NewJWTAuthenticator pour créer une nouvelle instance de JWTAuthenticator.

2. Mettez à jour client.go :
   - Modifiez la méthode NewClient pour accepter email et mot de passe :
     ```go
     func NewClient(email, password string, options ...ClientOption) (*Client, error) {
         // Initialisation du client avec l'authentificateur
     }
     ```
   - Ajoutez une méthode au Client pour effectuer des requêtes authentifiées.

3. Dans types.go :
   - Mettez à jour l'interface Authenticator si nécessaire.
   - Ajoutez les structures nécessaires pour la requête et la réponse d'authentification :
     ```go
     type LoginRequest struct {
         Email    string `json:"email"`
         Password string `json:"password"`
     }

     type LoginResponse struct {
         Token     string    `json:"token"`
         ExpiresAt time.Time `json:"expires_at"`
         User      User      `json:"user"`
     }

     type User struct {
         ID           string `json:"id"`
         Email        string `json:"email"`
         ProfileImage string `json:"profileImage"`
         FirstName    string `json:"firstName"`
     }
     ```

4. Créez un fichier auth_test.go :
   - Écrivez des tests unitaires pour JWTAuthenticator, couvrant les cas de succès et d'échec de l'authentification.

5. Mettez à jour le fichier examples/simple_client.go :
   - Ajoutez un exemple d'initialisation du client avec email et mot de passe.

6. Mettez à jour le README.md :
   - Ajoutez une section expliquant comment initialiser le client avec les informations d'authentification.

Directives importantes :
- L'authentification doit être effectuée via l'endpoint /api/login.
- Gérez correctement les erreurs d'authentification et les cas de token expiré.
- Implémentez un mécanisme de rafraîchissement automatique du token avant son expiration.
- Utilisez le logger du client pour enregistrer les événements importants liés à l'authentification.
- Respectez les limites de 500 lignes par fichier et 50 lignes par fonction.
- Suivez les conventions de nommage et de formatage de Go.
- Ajoutez des commentaires détaillés pour toutes les fonctions et méthodes exportées.

Après avoir terminé ces tâches, fournissez un rapport détaillé incluant :
- Un résumé des actions effectuées
- Le contenu des nouveaux fichiers et des modifications majeures
- Les décisions de conception prises, notamment concernant la gestion des tokens et leur rafraîchissement
- Toute difficulté rencontrée
- Des suggestions pour améliorer ou étendre le mécanisme d'authentification

Faisons le ensemble étape par étape. Toutes les textes dans le code doivent être en ANGLAIS.