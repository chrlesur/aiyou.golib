# Expression de Besoin : Package GO pour l'API AI.YOU de Cloud Temple

## 1. Objectif

Développer un package GO nommé "aiyou" qui implémente l'accès à l'API AI.YOU de Cloud Temple, permettant une interaction fluide et typée avec tous les endpoints décrits dans la spécification OpenAPI fournie.

## 2. Fonctionnalités principales

- Authentification et gestion des tokens JWT
- URL de base de l'API paramétrable
- Implémentation de tous les endpoints décrits dans l'API :
  - Login (`/api/login`)
  - Chat completions (`/api/v1/chat/completions`)
  - Création de modèles (`/api/v1/models`)
  - Récupération des assistants de l'utilisateur (`/api/v1/user/assistants`)
  - Sauvegarde d'une conversation (`/api/v1/save`)
  - Récupération des threads de l'utilisateur (`/api/v1/user/threads`)
  - Transcription audio (`/api/v1/audio/transcriptions`)
- Gestion des erreurs et des réponses pour chaque endpoint
- Support pour les opérations en streaming (pour les chat completions) comme option configurable
- Paramétrage complet des modèles (temperature, tokens, etc.) pour chaque appel

## 3. Structure du package

- Client principal avec des méthodes pour chaque endpoint
- Types Go correspondant aux schémas définis dans l'API
- Fonctions utilitaires pour la gestion des requêtes et des réponses
- Option de configuration pour l'URL de base de l'API

## 4. Implémentation détaillée

- Création d'un client avec options de configuration, incluant l'URL de base de l'API
- Méthodes pour chaque endpoint, respectant les signatures suivantes :
  - `Login(email, password string) (*LoginResponse, error)`
  - `ChatCompletions(ctx context.Context, input ChatCompletionInput) (*ChatCompletionOutput, error)`
  - `ChatCompletionsStream(ctx context.Context, input ChatCompletionInput) (<StreamReader>, error)`
  - `CreateModel(input ModelInput) error`
  - `GetUserAssistants() ([]UserAssistant, error)`
  - `SaveThread(input SaveThreadInput) (*SaveThreadOutput, error)`
  - `GetUserThreads(page, itemsPerPage int, search string) (*UserThreadsOutput, error)`
  - `TranscribeAudio(audioFile io.Reader) (string, error)`
- Implémentation du streaming pour les chat completions comme une option
- Structure `ChatCompletionInput` permettant un paramétrage complet du modèle :

```go
type ChatCompletionInput struct {
    Messages     []Message
    Model        string
    Temperature  float32
    MaxTokens    int
    TopP         float32
    FrequencyPenalty float32
    PresencePenalty  float32
    Stop         []string
    Stream       bool
    AssistantID  string
    // Autres paramètres pertinents...
}
```

## 5. Gestion des erreurs

- Création de types d'erreur spécifiques pour chaque type de réponse d'erreur de l'API
- Gestion appropriée des codes de statut HTTP
- Utilisation des pratiques standard Go pour la gestion des erreurs

## 6. Tests

- Tests unitaires pour chaque méthode du client
- Tests d'intégration avec des mocks pour simuler les réponses de l'API
- Tests spécifiques pour les options de streaming et les différents paramétrages de modèles
- Couverture de code visée : au moins 80%

## 7. Documentation

- Documentation GoDoc complète pour chaque fonction exportée
- Exemples d'utilisation pour chaque endpoint
- README détaillé avec guide de démarrage rapide
- Documentation au format Markdown compatible GitHub
- Documentation détaillée sur les options de paramétrage des modèles et l'utilisation du streaming
- Documentation sur la configuration de l'URL de base de l'API

## 8. Considérations de sécurité

- Gestion sécurisée des tokens JWT
- Utilisation de HTTPS pour toutes les communications

## 9. Bonnes pratiques de développement

- Limiter chaque fichier de code source Go à un maximum de 500 lignes
- Limiter chaque fonction à un maximum de 50 lignes de code
- Suivre les conventions de nommage et de formatage standard de Go
- Utiliser `gofmt` pour formater le code
- Utiliser `golint` et `go vet` pour la vérification du code
- Gérer toutes les erreurs de manière appropriée, sans utiliser de `panic`
- Utiliser des interfaces pour les composants principaux pour faciliter les tests et l'extensibilité future

## 10. Facilité d'utilisation

- API fluide et intuitive pour l'utilisateur final
- Gestion automatique de l'authentification et du rafraîchissement des tokens
- Options de configuration flexibles (URL de base, timeout, retries, paramètres de modèle par défaut, etc.)
- Méthodes distinctes pour les appels en streaming et non-streaming
- Paramètres de modèle facilement configurables pour chaque appel

## 11. Journalisation et gestion des erreurs

- Utiliser l'interface standard `log.Logger` de Go pour la journalisation
- Permettre à l'utilisateur de la librairie de fournir son propre logger
- Journaliser les événements importants (début/fin des appels API, erreurs) avec des niveaux de log appropriés
- Éviter de logger des informations sensibles (tokens, mots de passe)
- Suivre le modèle d'erreur de Go : "errors are values"
- Créer des types d'erreur personnalisés pour chaque catégorie d'erreur de l'API

## 12. Conformité aux normes de développement des librairies Go

- Structure de projet standard avec `/internal` et `/pkg`
- Utilisation de `go.mod` pour la gestion des dépendances
- Conventions de nommage Go standard
- Documentation complète et exemples de code
- Tests unitaires, sous-tests, et benchmarks
- Versionnage sémantique

## 13. Documentation au format Markdown

- README.md principal avec badges, description, installation, utilisation rapide, etc.
- Dossier `docs/` avec documentation détaillée
- Guides d'utilisation pour chaque fonctionnalité majeure
- CONTRIBUTING.md, CHANGELOG.md, et CODE_OF_CONDUCT.md
- Dossier `examples/` avec des exemples de code complets
- Conformité avec la syntaxe Markdown de GitHub
- Guide spécifique sur l'utilisation des options de streaming et le paramétrage des modèles

## 14. Informations spécifiques et licence

- Mainteneur principal : GitHub login "chrlesur"
- Email de contact : christophe.lesur@cloud-temple.com
- Licence : GNU General Public License v3.0 (GPL-3.0)
- Inclure le texte complet de la licence GPL-3.0 dans un fichier LICENSE
- Ajouter un en-tête de licence approprié à chaque fichier source
- Inclure une directive de licence dans le fichier go.mod

## 15. Exemple d'utilisation

```go
import (
    "context"
    "log"
    "os"

    "github.com/chrlesur/aiyou"
)

func main() {
    logger := log.New(os.Stdout, "aiyou: ", log.LstdFlags)

    client, err := aiyou.NewClient(
        "email@example.com",
        "password",
        aiyou.WithLogger(logger),
        aiyou.WithBaseURL("https://api.aiyou.example.com"), // URL paramétrable
    )
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    input := aiyou.ChatCompletionInput{
        Messages: []aiyou.Message{
            {Role: "user", Content: []aiyou.Content{{Type: "text", Text: "Bonjour"}}},
        },
        AssistantId: "602",
        Temperature: 0.7,
        MaxTokens: 150,
        Stream: false,
    }

    ctx := context.Background()
    output, err := client.ChatCompletions(ctx, input)
    if err != nil {
        var apiErr *aiyou.APIError
        if errors.As(err, &apiErr) {
            log.Printf("API error occurred: %v", apiErr)
        } else {
            log.Printf("Unexpected error: %v", err)
        }
        return
    }

    // Traitement de la réponse non-streaming
    log.Printf("Response: %s", output.Response)

    // Exemple d'utilisation du streaming
    streamInput := input
    streamInput.Stream = true
    stream, err := client.ChatCompletionsStream(ctx, streamInput)
    if err != nil {
        log.Fatalf("Failed to start stream: %v", err)
    }
    defer stream.Close()

    for {
        response, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Printf("Error receiving stream: %v", err)
            break
        }
        log.Printf("Streamed chunk: %s", response.Content)
    }
}
```

## 16. Gestion du Rate Limiting

- Implémentation d'un système de rate limiting côté client pour respecter les limites de l'API AI.YOU :
  - Utilisation d'un algorithme de rate limiting (ex: token bucket)
  - Configuration flexible des limites (requêtes par seconde, par minute, etc.)
  - Gestion automatique des retards entre les requêtes si nécessaire
  - Options pour personnaliser le comportement en cas de dépassement des limites :
    - Attente automatique
    - Retour d'erreur
    - Callback personnalisé

- Exposition des informations de rate limiting à l'utilisateur :
  - Méthodes pour vérifier le nombre de requêtes restantes
  - Événements ou callbacks pour notifier l'approche des limites

- Gestion des en-têtes de rate limiting renvoyés par l'API :
  - Parsing et stockage des informations de limite
  - Ajustement dynamique du rate limiting côté client en fonction des réponses de l'API

- Documentation claire sur l'utilisation et la configuration du rate limiting

Exemple d'utilisation :

```go
client, _ := aiyou.NewClient(
    "email@example.com",
    "password",
    aiyou.WithRateLimit(10, time.Second), // 10 requêtes par seconde
    aiyou.WithRateLimitBehavior(aiyou.RateLimitWait),
)

// Le client gérera automatiquement le rate limiting
```

## 17. Système de Cache

- Implémentation d'un système de cache optionnel pour certaines requêtes :
  - Cache en mémoire par défaut avec option pour des backends externes (ex: Redis)
  - Configuration de la durée de vie du cache par type de requête
  - Possibilité de désactiver le cache pour certaines requêtes spécifiques

- Types de requêtes à considérer pour le caching :
  - Récupération des assistants de l'utilisateur
  - Récupération des threads de l'utilisateur (avec invalidation appropriée)
  - Autres réponses relativement statiques de l'API

- Méthodes pour gérer manuellement le cache :
  - Invalidation de certaines entrées du cache
  - Préchargement de données dans le cache
  - Récupération de l'état actuel du cache

- Stratégies de rafraîchissement du cache :
  - Rafraîchissement en arrière-plan
  - Invalidation basée sur le temps

- Gestion des erreurs liées au cache :
  - Fallback vers l'API en cas d'échec du cache
  - Logging des erreurs de cache

- Documentation détaillée sur la configuration et l'utilisation du système de cache

Exemple d'utilisation :

```go
client, _ := aiyou.NewClient(
    "email@example.com",
    "password",
    aiyou.WithCache(aiyou.NewInMemoryCache()),
    aiyou.WithCacheTTL(aiyou.CacheKeyUserAssistants, 5*time.Minute),
)

// Les appels à GetUserAssistants utiliseront automatiquement le cache
assistants, err := client.GetUserAssistants()

// Invalidation manuelle du cache si nécessaire
client.Cache().Invalidate(aiyou.CacheKeyUserAssistants)
```
