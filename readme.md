
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

## Exemples

Vous pouvez trouver des exemples d'utilisation plus détaillés dans le dossier `examples/` de ce dépôt.

## Contribution

Les contributions sont les bienvenues ! N'hésitez pas à ouvrir une issue ou à soumettre une pull request.

## Licence

Ce projet est sous licence GNU General Public License v3.0 (GPL-3.0). Voir le fichier [LICENSE](LICENSE) pour plus de détails.
