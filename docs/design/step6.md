Contexte du projet :
Vous travaillez sur un package Go nommé "aiyou.golib" qui sert d'interface pour l'API AI.YOU de Cloud Temple. Ce package permet aux développeurs Go d'interagir facilement avec l'API AI.YOU. Les étapes précédentes ont mis en place la structure de base du projet, l'authentification, l'implémentation de l'endpoint de chat completion, et récemment, un système avancé de gestion des erreurs et de retry.

Objectif de cette étape :
Implémenter un système de logging avancé dans le package aiyou.golib. Ce système doit permettre un suivi détaillé des opérations, être configurable, et s'intégrer harmonieusement avec les fonctionnalités existantes, notamment le système de retry.

État actuel du projet :
- Le projet est structuré dans un dossier "pkg/aiyou" contenant les fichiers principaux.
- L'authentification, le chat completion, et un système de gestion des erreurs et de retry sont en place.
- Un système de logging de base existe, mais nécessite des améliorations significatives.

Tâches à réaliser :

1. Création d'un package de logging :
   - Créez un nouveau fichier `logging.go` dans le package aiyou.
   - Implémentez une interface Logger qui étend l'interface standard `log.Logger` de Go :
     ```go
     type Logger interface {
         log.Logger
         Debugf(format string, args ...interface{})
         Infof(format string, args ...interface{})
         Warnf(format string, args ...interface{})
         Errorf(format string, args ...interface{})
     }
     ```

2. Implémentation d'un logger par défaut :
   - Créez une structure `defaultLogger` qui implémente l'interface `Logger`.
   - Ajoutez des options pour configurer le niveau de log (DEBUG, INFO, WARN, ERROR).

3. Intégration du logger dans le client :
   - Modifiez la structure `Client` dans `client.go` pour inclure le nouveau logger :
     ```go
     type Client struct {
         // ... autres champs existants
         logger Logger
     }
     ```
   - Ajoutez une option pour configurer le logger lors de la création du client :
     ```go
     func WithLogger(logger Logger) ClientOption {
         // Implémentez cette fonction
     }
     ```

4. Utilisation du logger dans les méthodes existantes :
   - Mettez à jour toutes les méthodes du client pour utiliser le nouveau système de logging.
   - Assurez-vous de logger les événements importants, comme les débuts et fins d'appels API, les retries, et les erreurs.

5. Intégration avec le système de retry :
   - Modifiez la fonction de retry dans `retry.go` pour utiliser le nouveau logger :
     ```go
     func retryOperation(ctx context.Context, logger Logger, maxRetries int, operation func() error) error {
         // Ajoutez des logs détaillés pour chaque tentative et résultat
     }
     ```

6. Masquage des informations sensibles :
   - Implémentez une fonction utilitaire pour masquer les informations sensibles dans les logs (tokens, mots de passe, etc.).

7. Tests unitaires :
   - Créez des tests dans `logging_test.go` pour vérifier le bon fonctionnement du nouveau système de logging.
   - Mettez à jour les tests existants pour utiliser et vérifier le logging approprié.

8. Documentation :
   - Mettez à jour la documentation dans les fichiers Go pour refléter les nouvelles fonctionnalités de logging.
   - Ajoutez une section dans le README.md expliquant comment configurer et utiliser le système de logging.

Directives importantes :
- Suivez les conventions de nommage et de formatage standard de Go.
- Assurez-vous que le code est bien commenté et documenté (format godoc).
- Respectez les limites de 500 lignes par fichier et 50 lignes par fonction.
- Le code doit être en anglais, y compris les commentaires et la documentation.
- Veillez à ce que le système de logging soit performant et n'impacte pas significativement les performances globales du package.

Après avoir terminé ces tâches, fournissez un rapport détaillé incluant :
- Un résumé des actions effectuées
- Le contenu des nouveaux fichiers et des modifications majeures
- Les décisions de conception prises, notamment concernant la structure du logger et son intégration avec les fonctionnalités existantes
- Toute difficulté rencontrée
- Des suggestions pour améliorer davantage le système de logging ou son utilisation dans le package

---
