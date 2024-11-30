# aiyou.golib

aiyou.golib est un package Go pour interagir avec l'API AI.YOU de Cloud Temple.

## Installation

Pour installer aiyou.golib, utilisez la commande suivante :

go get github.com/chrlesur/aiyou.golib

## Structure du Package

L'organisation des fichiers et répertoires du package :

    aiyou.golib/
    ├── LICENSE.txt # Licence GPL-3.0
    ├── README.md # Documentation principale
    ├── go.mod # Définition du module Go
    ├── aiyou.go # Point d'entrée du package, exports publics
    │
    ├── pkg/ # Code source principal
    │ └── aiyou/
    │ ├── assistants.go # Gestion des assistants IA
    │ ├── audio.go # Transcription audio
    │ ├── auth.go # Authentification et gestion des tokens
    │ ├── chat.go # Fonctionnalités de chat completion
    │ ├── client.go # Client HTTP et configuration
    │ ├── config.go # Structures de configuration
    │ ├── conversation.go # Gestion des conversations
    │ ├── errors.go # Types d'erreurs personnalisés
    │ ├── logging.go # Système de logging
    │ ├── models.go # Définitions des modèles
    │ ├── ratelimit.go # Système de rate limiting
    │ ├── retry.go # Logique de retry
    │ ├── threads.go # Gestion des threads de conversation
    │ ├── types.go # Définitions des types communs
    │ └── utils.go # Fonctions utilitaires
    │
    ├── examples/ # Exemples d'utilisation
    │ ├── assistants.go # Exemple de gestion des assistants
    │ ├── audio.go # Exemple de transcription audio
    │ ├── conversation.go # Exemple de gestion des conversations
    │ ├── message_builder.go # Exemple d'utilisation du MessageBuilder
    │ ├── models.go # Exemple d'utilisation des modèles
    │ ├── rate_limiting.go # Exemple de rate limiting simple
    │ ├── rate_limiting_advanced.go # Exemple de rate limiting avancé
    │ ├── simple_client.go # Client en ligne de commande complet
    │ └── threads.go # Exemple de gestion des threads

### Description des composants principaux

#### Fichiers racine
- `aiyou.go` : Point d'entrée principal du package, expose l'API publique
- `go.mod` : Définition du module et de ses dépendances

#### Package principal (pkg/aiyou)
- **Cœur du client**
- `client.go` : Implémentation du client HTTP principal
- `config.go` : Structures et logique de configuration
- `types.go` : Définitions des types de données communs

- **Fonctionnalités**
- `chat.go` : Implémentation des fonctionnalités de chat
- `audio.go` : Gestion de la transcription audio
- `assistants.go` : Gestion des assistants IA
- `conversation.go` : Gestion des conversations
- `threads.go` : Gestion des threads de discussion

- **Infrastructure**
- `auth.go` : Système d'authentification JWT
- `logging.go` : Système de logging avec protection des données sensibles
- `ratelimit.go` : Implémentation du rate limiting
- `retry.go` : Logique de retry des requêtes
- `errors.go` : Types d'erreurs personnalisés

#### Exemples
Les exemples dans le dossier `examples/` démontrent des cas d'utilisation concrets et servent de documentation interactive. Chaque exemple peut être exécuté indépendamment et inclut des options de configuration complètes.

### Tests

Chaque composant majeur dispose de ses propres tests unitaires (fichiers `*_test.go`). Les tests couvrent :

- Les cas d'utilisation normaux
- La gestion des erreurs
- Les cas limites
- Les performances (benchmarks)

## Utilisation

### Initialisation du client avec options

    import "github.com/chrlesur/aiyou.golib"

    // Client simple
    client, err := aiyou.NewClient("votre-email@exemple.com", "votre-mot-de-passe")
    if err != nil {
    log.Fatalf("Erreur lors de la création du client : %v", err)
    }

    // Client avec options
    client, err := aiyou.NewClient(
    "votre-email@exemple.com",
    "votre-mot-de-passe",
    aiyou.WithBaseURL("https://ai.dragonflygroup.fr"),
    aiyou.WithLogger(customLogger),
    aiyou.WithRateLimiter(aiyou.RateLimiterConfig{
    RequestsPerSecond: 2.0,
    BurstSize: 3,
    WaitTimeout: 5 * time.Second,
    }),
    )

### Authentification

L'authentification est gérée automatiquement par le client. Vous n'avez pas besoin de vous authentifier manuellement avant chaque requête.

### Mode Quiet

Un mode silencieux est disponible pour réduire les logs et les sorties :

    client, err := aiyou.NewClient(
    "email@exemple.com",
    "password",
    aiyou.WithQuietMode(true),
    )

### Chat Completion

aiyou.golib fournit deux méthodes principales pour le chat completion :

#### Chat Completion Standard

    req := aiyou.ChatCompletionRequest{
    Messages: []aiyou.Message{
    {
    Role: "user",
    Content: []aiyou.ContentPart{
    {Type: "text", Text: "Quelle est la capitale de la France ?"},
    },
    },
    },
    AssistantID: "id-de-votre-assistant",
    }

    resp, err := client.ChatCompletion(context.Background(), req)
    if err != nil {
    log.Fatalf("Erreur lors du chat completion : %v", err)
    }

    fmt.Printf("Réponse de l'IA : %s\n", resp.Choices[0].Message.Content[0].Text)

#### Chat Completion en Streaming

    streamReq := aiyou.ChatCompletionRequest{
    Messages: []aiyou.Message{
    {
    Role: "user",
    Content: []aiyou.ContentPart{
    {Type: "text", Text: "Raconte-moi une courte histoire."},
    },
    },
    },
    AssistantID: "id-de-votre-assistant",
    Stream: true,
    }

    stream, err := client.ChatCompletionStream(context.Background(), streamReq)
    if err != nil {
    log.Fatalf("Erreur lors du chat completion en streaming : %v", err)
    }

    for {
    chunk, err := stream.ReadChunk()
    if err == io.EOF {
    break
    }
    if err != nil {
    log.Fatalf("Erreur lors de la lecture du chunk : %v", err)
    }
    fmt.Print(chunk.Choices[0].Delta.Content)
    }

### Transcription Audio

Le package supporte la transcription de fichiers audio :

    // Transcription simple
    opts := &aiyou.AudioTranscriptionRequest{
    Language: "fr",
    Format: "text",
    }

    transcription, err := client.TranscribeAudioFile(context.Background(), "chemin/vers/audio.wav", opts)
    if err != nil {
    log.Fatalf("Erreur de transcription: %v", err)
    }
    fmt.Println(transcription.Transcription)

Utilisation via la ligne de commande :
    go run examples/audio.go --email="user@example.com" --password="pass" --file="audio.wav" --lang="fr" --format="text"

Formats supportés :
- WAV (jusqu'à 25MB)
- MP3 (jusqu'à 25MB)
- M4A (jusqu'à 25MB)

### Gestion des erreurs et Retry

Le package aiyou.golib implémente une gestion avancée des erreurs et un système de retry.

#### Types d'erreurs personnalisés

- `APIError`: Erreurs retournées par l'API AI.YOU
- `AuthenticationError`: Erreurs liées à l'authentification
- `RateLimitError`: Erreurs de dépassement de limite de taux
- `NetworkError`: Erreurs de réseau

#### Système de retry

    client, err := aiyou.NewClient(
    "your-email@example.com",
    "your-password",
    aiyou.WithRetry(3, time.Second),
    )

### Logging

Le package inclut un système de logging flexible qui protège les informations sensibles.

#### Configuration du Logger

    customLogger := aiyou.NewDefaultLogger(os.Stdout)
    client, err := aiyou.NewClient(
    "your-email@example.com",
    "your-password",
    aiyou.WithLogger(customLogger),
    )

#### Niveaux de Log

Le système de logging supporte quatre niveaux :
- `DEBUG` : Informations détaillées de débogage
- `INFO` : Informations opérationnelles générales
- `WARN` : Messages d'avertissement
- `ERROR` : Messages d'erreur

    customLogger.SetLevel(aiyou.DEBUG)  

### Rate Limiting

aiyou.golib inclut un système de rate limiting configurable pour contrôler le débit des requêtes.

#### Configuration du Rate Limiting

    client, err := aiyou.NewClient(
    "your-email@example.com",
    "your-password",
    aiyou.WithRateLimiter(aiyou.RateLimiterConfig{
    RequestsPerSecond: 10, // Limite de requêtes par seconde
    BurstSize: 5, // Nombre de requêtes autorisées en burst
    WaitTimeout: 5 * time.Second, // Timeout d'attente maximum
    }),
    )

#### Gestion des Erreurs de Rate Limiting

    resp, err := client.ChatCompletion(ctx, req)
    if err != nil {
    switch e := err.(type) {
    case *aiyou.RateLimitError:
    if e.IsClientSide {
    fmt.Printf("Rate limit local dépassé. Réessayer dans %d secondes\n", e.RetryAfter)
    } else {
    fmt.Printf("Quota API dépassé. Réessayer dans %d secondes\n", e.RetryAfter)
    }
    }
    return
    }

#### Utilisation avec des Requêtes Concurrentes

    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(i int) {
    defer wg.Done()
    
    ctx := context.Background()
    msg := aiyou.NewTextMessage("user", fmt.Sprintf("Request %d", i))
    
    resp, err := client.CreateChatCompletion(ctx, []aiyou.Message{msg}, "assistant-id")
    if err != nil {
    log.Printf("Request %d failed: %v", i, err)
    return
    }
    log.Printf("Request %d successful", i)
    }(i)
    }
    wg.Wait()

## Exemples

Des exemples complets sont disponibles dans le dossier `examples/` :

# Chat interactif avec historique et commandes
    go run examples/simple_client.go --email="user@example.com" --password="pass" --assistant="asst_123"

# Gestion des assistants
    go run examples/assistants.go --email="user@example.com" --password="pass"

# Test de rate limiting simple
    go run examples/rate_limiting.go --email="user@example.com" --password="pass" --rate=2.0

# Test de rate limiting avancé
    go run examples/rate_limiting_advanced.go --email="user@example.com" --password="pass" --requests=20 --rate=1.0 --burst=2

# Transcription audio
    go run examples/audio.go --email="user@example.com" --password="pass" --file="audio.wav" --lang="fr"

### Options communes des exemples

Tous les exemples supportent les options suivantes :
- `--email` : Email pour l'authentification
- `--password` : Mot de passe
- `--url` : URL de base de l'API (optionnel)
- `--debug` : Active les logs de debug
- `--quiet` : Mode silencieux

## Contribution

Les contributions sont les bienvenues ! N'hésitez pas à ouvrir une issue ou à soumettre une pull request.

## Licence

Ce projet est sous licence GNU General Public License v3.0 (GPL-3.0). Voir le fichier [LICENSE](LICENSE) pour plus de détails.
