Contexte du projet :
Vous travaillez sur un package Go nommé "aiyou.golib" qui sert d'interface pour l'API AI.YOU de Cloud Temple. Ce package permet aux développeurs Go d'interagir facilement avec l'API AI.YOU. Les étapes précédentes ont mis en place la structure de base du projet et l'authentification.

Objectif de cette étape :
Implémenter l'endpoint de chat completion dans le package aiyou.golib. Cet endpoint est crucial car il permet aux utilisateurs d'interagir avec le modèle de langage AI.YOU pour générer des réponses basées sur des prompts donnés.

État actuel du projet :
- Le projet est structuré avec un dossier principal "pkg/aiyou" contenant les fichiers principaux.
- L'authentification par email/mot de passe est déjà implémentée.
- Le client de base est configuré pour faire des requêtes authentifiées.

Tâches à réaliser :

1. Dans types.go, ajoutez les structures nécessaires pour les requêtes et réponses de chat completion :
   ```go
   type ChatCompletionRequest struct {
       Messages []Message `json:"messages"`
       AssistantID string `json:"assistantId"`
       Temperature float32 `json:"temperature"`
       TopP float32 `json:"top_p"`
       PromptSystem string `json:"promptSystem"`
       Stream bool `json:"stream"`
       Stop []string `json:"stop,omitempty"`
       ThreadID string `json:"threadId,omitempty"`
   }

   type Message struct {
       Role string `json:"role"`
       Content []ContentPart `json:"content"`
   }

   type ContentPart struct {
       Type string `json:"type"`
       Text string `json:"text"`
   }

   type ChatCompletionResponse struct {
       // Ajoutez les champs nécessaires selon la spécification de l'API
   }
   ```

2. Créez un nouveau fichier chat.go dans le dossier pkg/aiyou et implémentez les fonctions suivantes :
   - `(c *Client) ChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error)`
   - `(c *Client) ChatCompletionStream(ctx context.Context, req ChatCompletionRequest) (<StreamReader>, error)`

3. Dans client.go, ajoutez des méthodes pour faciliter l'utilisation de ces fonctions :
   ```go
   func (c *Client) CreateChatCompletion(ctx context.Context, messages []Message, assistantID string) (*ChatCompletionResponse, error)
   func (c *Client) CreateChatCompletionStream(ctx context.Context, messages []Message, assistantID string) (<StreamReader>, error)
   ```

4. Implémentez la gestion du streaming dans chat.go. Utilisez `io.Reader` pour gérer les réponses en streaming.

5. Ajoutez des tests unitaires dans un fichier chat_test.go pour toutes les nouvelles fonctions.

6. Mettez à jour le fichier examples/simple_client.go avec un exemple d'utilisation de la fonction de chat completion.

7. Mettez à jour le README.md avec des informations sur comment utiliser les nouvelles fonctions de chat completion.

Directives importantes :
- Suivez les conventions de nommage et de formatage standard de Go.
- Gérez les erreurs de manière appropriée, en créant des types d'erreur personnalisés si nécessaire.
- Utilisez le système de logging existant pour enregistrer les événements importants.
- Assurez-vous que le code est bien commenté et documenté (format godoc).
- Respectez les limites de 500 lignes par fichier et 50 lignes par fonction.
- Le code doit être en anglais, y compris les commentaires et la documentation.

Après avoir terminé ces tâches, fournissez un rapport détaillé incluant :
- Un résumé des actions effectuées
- Le contenu des nouveaux fichiers et des modifications majeures
- Les décisions de conception prises, notamment concernant la gestion du streaming
- Toute difficulté rencontrée
- Des suggestions pour améliorer ou étendre la fonctionnalité de chat completion
