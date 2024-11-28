# aiyou.golib

aiyou.golib est un package Go pour interagir avec l'API AI.YOU de Cloud Temple.

## Installation

Pour installer aiyou.golib, utilisez la commande suivante :

go get github.com/chrlesur/aiyou.golib

## Utilisation

### Initialisation du client

Pour commencer à utiliser aiyou.golib, vous devez d'abord initialiser un client :

import "github.com/chrlesur/aiyou.golib"

client, err := aiyou.NewClient("votre-email@exemple.com", "votre-mot-de-passe")
if err != nil {
 log.Fatalf("Erreur lors de la création du client : %v", err)
}

### Authentification

L'authentification est gérée automatiquement par le client. Vous n'avez pas besoin de vous authentifier manuellement avant chaque requête.

## Chat Completion

aiyou.golib fournit deux méthodes principales pour le chat completion :

### Chat Completion Standard

Utilisez la méthode `ChatCompletion` pour une requête de chat completion standard :

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

### Chat Completion en Streaming

Pour les réponses en streaming, utilisez la méthode `ChatCompletionStream` :

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
 fmt.Print(chunk.Choices[0].Message.Content[0].Text)
}

## Gestion des erreurs et Retry

Le package aiyou.golib implémente une gestion avancée des erreurs et un système de retry pour améliorer la robustesse des interactions avec l'API AI.YOU.

### Types d'erreurs personnalisés

- `APIError`: Erreurs retournées par l'API AI.YOU
- `AuthenticationError`: Erreurs liées à l'authentification
- `RateLimitError`: Erreurs de dépassement de limite de taux
- `NetworkError`: Erreurs de réseau

### Système de retry

Le client peut être configuré pour réessayer automatiquement les opérations en cas d'erreurs temporaires :

client, err := aiyou.NewClient(
 "your-email@example.com",
 "your-password",
 aiyou.WithRetry(3, time.Second),
)

## Logging

Le package aiyou.golib inclut un système de logging flexible qui vous permet de contrôler la verbosité des logs et de protéger les informations sensibles.

### Configuration du Logger

customLogger := aiyou.NewDefaultLogger(os.Stdout)
client, err := aiyou.NewClient(
 "your-email@example.com",
 "your-password",
 aiyou.WithLogger(customLogger),
)

### Niveaux de Log

Le système de logging supporte quatre niveaux :
- DEBUG : Informations détaillées de débogage
- INFO : Informations opérationnelles générales
- WARN : Messages d'avertissement
- ERROR : Messages d'erreur

customLogger.SetLevel(aiyou.DEBUG)

## Rate Limiting

aiyou.golib inclut un système de rate limiting configurable pour contrôler le débit des requêtes vers l'API AI.YOU et éviter les erreurs de dépassement de quota.

### Configuration du Rate Limiting

client, err := aiyou.NewClient(
 "your-email@example.com",
 "your-password",
 aiyou.WithRateLimiter(aiyou.RateLimiterConfig{
 RequestsPerSecond: 10, // Limite de requêtes par seconde
 BurstSize: 5, // Nombre de requêtes autorisées en burst
 WaitTimeout: time.Second * 5, // Timeout d'attente maximum
 }),
)

### Options de Configuration

- `RequestsPerSecond` : Définit le nombre maximum de requêtes autorisées par seconde
- `BurstSize` : Permet un certain nombre de requêtes à exécuter immédiatement
- `WaitTimeout` : Durée maximale d'attente avant de retourner une erreur de rate limit

### Gestion des Erreurs de Rate Limiting

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

### Utilisation avec des Requêtes Concurrentes

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

### Performance et Impact

Le rate limiting est implémenté avec un algorithme de token bucket efficace qui ajoute une overhead négligeable (< 1ms) aux requêtes. Le système est conçu pour :
- Minimiser l'impact sur les performances en utilisation normale
- Gérer efficacement les pics de charge
- Éviter les contentions en cas d'accès concurrent

### Désactivation du Rate Limiting

client, err := aiyou.NewClient(
 "your-email@example.com",
 "your-password",
 // Pas d'option WithRateLimiter = pas de rate limiting
)

## Optimisation des Performances

Le package a été optimisé pour les performances, particulièrement en termes d'utilisation mémoire et de vitesse d'exécution. Voici les métriques clés :

### Mémoire et Vitesse
- Opérations HTTP : ~120μs par opération
- Utilisation mémoire : ~13KB par opération
- Allocations mémoire : 238 allocations par opération

### Tests de Performance
Pour exécuter les benchmarks de performance :
# Lancer les benchmarks
go test ./pkg/aiyou -run=^$ -bench=BenchmarkChatCompletion -benchmem

# Générer et analyser le profil mémoire
go test ./pkg/aiyou -run=^$ -bench=BenchmarkChatCompletion -benchmem -memprofile=mem.prof
go tool pprof mem.prof

### Optimisations
Le package inclut plusieurs optimisations de performance :
- Pool de buffers pour les opérations de logging
- Gestion efficace de la mémoire pour les opérations HTTP
- Traitement optimisé des chaînes de caractères pour le masquage des informations sensibles

Ces optimisations ont permis :
- Une exécution 50% plus rapide
- Une réduction de 90% de l'utilisation mémoire
- Une réduction de 79% des allocations mémoire

Pour les environnements de production, il est recommandé de surveiller l'utilisation de la mémoire et d'ajuster les niveaux de logging en conséquence.

## Exemples

Vous pouvez trouver des exemples d'utilisation plus détaillés dans le dossier `examples/` de ce dépôt.

## Contribution

Les contributions sont les bienvenues ! N'hésitez pas à ouvrir une issue ou à soumettre une pull request.

## Licence

Ce projet est sous licence GNU General Public License v3.0 (GPL-3.0). Voir le fichier [LICENSE](LICENSE) pour plus de détails.
C'est un README.md complet qui couvre toutes les fonctionnalités principales du package, y compris la nouvelle fonctionnalité de rate limiting, tout en maintenant une structure cohérente et claire.