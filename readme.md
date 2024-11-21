
# aiyou.golib

aiyou.golib est un package Go pour interagir avec l'API AI.YOU de Cloud Temple.

## Installation

Pour installer aiyou.golib, utilisez la commande suivante :

```
go get github.com/chrlesur/aiyou.golib
```

## Utilisation

### Initialisation du client

Pour commencer à utiliser aiyou.golib, vous devez d'abord initialiser un client :

```go
import "github.com/chrlesur/aiyou.golib"

client, err := aiyou.NewClient("votre-email@exemple.com", "votre-mot-de-passe")
if err != nil {
    log.Fatalf("Erreur lors de la création du client : %v", err)
}
```

### Authentification

L'authentification est gérée automatiquement par le client. Vous n'avez pas besoin de vous authentifier manuellement avant chaque requête.

## Chat Completion

aiyou.golib fournit deux méthodes principales pour le chat completion :

### Chat Completion Standard

Utilisez la méthode `ChatCompletion` pour une requête de chat completion standard :

```go
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
```

### Chat Completion en Streaming

Pour les réponses en streaming, utilisez la méthode `ChatCompletionStream` :

```go
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
    Stream:      true,
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
    fmt.Print(chunk.Choices[0].Message.Content[0].Text)
}
```

Assurez-vous de gérer les erreurs de manière appropriée dans votre application.

Bien sûr, voici le bout de README au format Markdown :

## Gestion des erreurs et Retry

Le package aiyou.golib implémente une gestion avancée des erreurs et un système de retry pour améliorer la robustesse des interactions avec l'API AI.YOU.

### Types d'erreurs personnalisés

- `APIError`: Erreurs retournées par l'API AI.YOU
- `AuthenticationError`: Erreurs liées à l'authentification
- `RateLimitError`: Erreurs de dépassement de limite de taux
- `NetworkError`: Erreurs de réseau

### Système de retry

Le client peut être configuré pour réessayer automatiquement les opérations en cas d'erreurs temporaires. Voici comment configurer le retry lors de la création du client :

```go
client, err := aiyou.NewClient(
    "your-email@example.com",
    "your-password",
    aiyou.WithRetry(3, time.Second),
)
```

Cet exemple configure le client pour effectuer jusqu'à 3 tentatives, avec un délai initial d'une seconde entre chaque tentative. Le délai augmente de façon exponentielle après chaque échec.

### Exemple d'utilisation

```go
resp, err := client.ChatCompletion(ctx, req)
if err != nil {
    switch e := err.(type) {
    case *aiyou.APIError:
        fmt.Printf("API error: %d - %s\n", e.StatusCode, e.Message)
    case *aiyou.AuthenticationError:
        fmt.Printf("Authentication error: %s\n", e.Message)
    case *aiyou.RateLimitError:
        fmt.Printf("Rate limit exceeded. Retry after %d seconds\n", e.RetryAfter)
    case *aiyou.NetworkError:
        fmt.Printf("Network error: %v\n", e.Err)
    default:
        fmt.Printf("Unexpected error: %v\n", err)
    }
    return
}
```

Cette gestion des erreurs et le système de retry améliorent la fiabilité des applications utilisant aiyou.golib, en gérant automatiquement les erreurs temporaires et en fournissant des informations détaillées sur les erreurs rencontrées.

Bien sûr, voici le README complet en format Markdown :

```markdown
# aiyou.golib

aiyou.golib est un package Go qui fournit une interface pour interagir avec l'API AI.YOU de Cloud Temple.

## Installation

Pour installer aiyou.golib, utilisez la commande go get :

```bash
go get github.com/chrlesur/aiyou.golib
```

## Utilisation

Voici un exemple simple d'utilisation du client :

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/chrlesur/aiyou.golib"
)

func main() {
    client, err := aiyou.NewClient("your-email@example.com", "your-password")
    if err != nil {
        log.Fatalf("Error creating client: %v", err)
    }

    ctx := context.Background()

    req := aiyou.ChatCompletionRequest{
        Messages: []aiyou.Message{
            {
                Role: "user",
                Content: []aiyou.ContentPart{
                    {Type: "text", Text: "What is the capital of France?"},
                },
            },
        },
        AssistantID: "your-assistant-id",
    }

    resp, err := client.ChatCompletion(ctx, req)
    if err != nil {
        log.Fatalf("Error in ChatCompletion: %v", err)
    }

    fmt.Printf("AI response: %s\n", resp.Choices[0].Message.Content[0].Text)
}
```

## Configuration

Le client peut être configuré avec plusieurs options :

```go
client, err := aiyou.NewClient(
    "your-email@example.com",
    "your-password",
    aiyou.WithBaseURL("https://custom-api-url.com"),
    aiyou.WithRetry(5, time.Second),
    aiyou.WithLogger(customLogger),
)
```

## Logging

Le package aiyou.golib inclut un système de logging flexible qui vous permet de contrôler la verbosité des logs et de protéger les informations sensibles.

### Configuration du Logger

Par défaut, le client utilise un logger de base qui écrit dans `os.Stderr`. Vous pouvez fournir votre propre logger personnalisé lors de la création d'un nouveau client :

```go
customLogger := aiyou.NewDefaultLogger(os.Stdout, "myapp: ", log.LstdFlags)
client, err := aiyou.NewClient(
    "your-email@example.com",
    "your-password",
    aiyou.WithLogger(customLogger),
)
```

### Niveaux de Log

Le système de logging supporte quatre niveaux de log :

- DEBUG : Informations détaillées de débogage
- INFO : Informations opérationnelles générales
- WARN : Messages d'avertissement
- ERROR : Messages d'erreur

Vous pouvez définir le niveau de log sur le logger par défaut :

```go
customLogger.SetLevel(aiyou.DEBUG)
```

### Logging Sécurisé

Le client utilise automatiquement un logging sécurisé pour masquer les informations sensibles telles que les adresses email, les tokens JWT et les mots de passe. Cela garantit que les données sensibles ne sont pas exposées dans vos logs.

### Logging Personnalisé

Si vous avez besoin de logger des informations supplémentaires dans votre application tout en utilisant le client, vous pouvez utiliser la fonction `SafeLog` pour vous assurer que les informations sensibles sont masquées :

```go
safeLog := aiyou.SafeLog(customLogger)
safeLog(aiyou.INFO, "Traitement de la requête pour l'utilisateur : %s", emailUtilisateur)
```

Cela logguera le message avec l'adresse email masquée.

### Logging dans les Implémentations Personnalisées

Si vous étendez la fonctionnalité du client ou créez des implémentations personnalisées, assurez-vous d'utiliser la méthode `safeLog` fournie par le client au lieu d'utiliser directement le logger. Cela garantit que tous les logs passent par le même processus de masquage des informations sensibles :

```go
func (c *CustomClient) SomeMethod() error {
    c.safeLog(aiyou.DEBUG, "Exécution d'une opération avec le token : %s", someToken)
    // Reste de l'implémentation de la méthode
}
```

En suivant ces directives, vous pouvez vous assurer que votre application log des informations utiles pour le débogage et la surveillance, tout en protégeant les données sensibles.

## Exemples

Vous pouvez trouver des exemples d'utilisation plus détaillés dans le dossier `examples/` de ce dépôt.

## Contribution

Les contributions sont les bienvenues ! N'hésitez pas à ouvrir une issue ou à soumettre une pull request.

## Licence

Ce projet est sous licence GNU General Public License v3.0 (GPL-3.0). Voir le fichier [LICENSE](LICENSE) pour plus de détails.
