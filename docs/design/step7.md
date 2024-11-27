
Contexte du projet :
Vous travaillez sur un package Go nommé "aiyou.golib" qui sert d'interface pour l'API AI.YOU de Cloud Temple. Ce package permet aux développeurs Go d'interagir facilement avec l'API AI.YOU. Les étapes précédentes ont mis en place la structure de base du projet, l'authentification, l'implémentation de l'endpoint de chat completion, un système avancé de gestion des erreurs et de retry, ainsi qu'un système de logging avancé.

Objectif de cette étape :
Développer des fonctions helper pour faciliter la construction de messages complexes pour le chat completion. Ces helpers doivent permettre aux utilisateurs de créer facilement des messages contenant du texte, des images, ou d'autres types de contenu supportés par l'API AI.YOU.

État actuel du projet :
- Le projet est structuré dans un dossier "pkg/aiyou" contenant les fichiers principaux.
- L'authentification, le chat completion, la gestion des erreurs, le retry et le logging avancé sont en place.
- Les structures de base pour les messages existent déjà dans types.go.

Tâches à réaliser :

1. Création d'un nouveau fichier pour les helpers :
   - Créez un fichier `message_helpers.go` dans le package aiyou.

2. Implémentation des fonctions helper pour les différents types de contenu :
   - Créez des fonctions pour construire facilement des messages texte, des images, et d'autres types de contenu supportés par l'API. Par exemple :
     ```go
     func NewTextMessage(role, content string) Message {
         // Implémentez cette fonction
     }

     func NewImageMessage(role, imageURL string) Message {
         // Implémentez cette fonction
     }

     // Ajoutez d'autres fonctions helper selon les besoins
     ```

3. Implémentation d'un builder pour les messages complexes :
   - Créez une structure `MessageBuilder` qui permet de construire des messages avec plusieurs types de contenu. Par exemple :
     ```go
     type MessageBuilder struct {
         message Message
     }

     func NewMessageBuilder(role string) *MessageBuilder {
         // Implémentez cette fonction
     }

     func (mb *MessageBuilder) AddText(text string) *MessageBuilder {
         // Implémentez cette méthode
     }

     func (mb *MessageBuilder) AddImage(imageURL string) *MessageBuilder {
         // Implémentez cette méthode
     }

     func (mb *MessageBuilder) Build() Message {
         // Implémentez cette méthode
     }
     ```

4. Mise à jour des structures existantes si nécessaire :
   - Si besoin, mettez à jour les structures de Message et ContentPart dans `types.go` pour supporter tous les types de contenu.

5. Intégration avec le client existant :
   - Mettez à jour la méthode ChatCompletion dans `chat.go` pour utiliser facilement ces nouveaux helpers.

6. Tests unitaires :
   - Créez des tests dans `message_helpers_test.go` pour vérifier le bon fonctionnement de tous les helpers et du builder.

7. Documentation :
   - Ajoutez des commentaires détaillés pour chaque fonction helper et méthode du builder.
   - Mettez à jour le README.md avec des exemples d'utilisation des nouveaux helpers.

8. Exemple d'utilisation :
   - Créez un nouvel exemple dans le dossier `examples/` montrant comment utiliser les helpers pour construire des messages complexes.

Directives importantes :
- Suivez les conventions de nommage et de formatage standard de Go.
- Assurez-vous que le code est bien commenté et documenté (format godoc).
- Respectez les limites de 500 lignes par fichier et 50 lignes par fonction.
- Le code doit être en anglais, y compris les commentaires et la documentation.
- Utilisez le système de logging existant pour enregistrer les événements importants liés à la création des messages.

Après avoir terminé ces tâches, fournissez un rapport détaillé incluant :
- Un résumé des actions effectuées
- Le contenu des nouveaux fichiers et des modifications majeures
- Les décisions de conception prises, notamment concernant la structure du builder et des helpers
- Toute difficulté rencontrée
- Des suggestions pour améliorer davantage la création et la gestion des messages complexes
