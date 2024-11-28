Contexte du projet :
Vous travaillez sur un package Go nommé "aiyou.golib" qui sert d'interface pour l'API AI.YOU de Cloud Temple. Ce package permet aux développeurs Go d'interagir facilement avec l'API AI.YOU. Les étapes précédentes ont mis en place la structure de base du projet, l'authentification, l'implémentation de tous les endpoints de l'API, un système avancé de gestion des erreurs et de retry, un système de logging avancé, et des helpers pour la création de messages complexes.

Objectif de cette étape :
Implémenter un système de rate limiting côté client dans le package aiyou.golib. Ce système doit permettre de contrôler le nombre de requêtes envoyées à l'API AI.YOU pour respecter les limites imposées par le service et éviter les erreurs de dépassement de quota.

État actuel du projet :
- Le projet est structuré dans un dossier "pkg/aiyou" contenant les fichiers principaux.
- Un client de base avec authentification et gestion des erreurs est en place.
- Tous les endpoints de l'API sont implémentés.
- Un système de logging avancé est intégré.

Tâches à réaliser :

1. Création du système de rate limiting :
   - Créez un nouveau fichier `ratelimit.go` dans le package aiyou.
   - Implémentez un rate limiter basé sur l'algorithme du seau à jetons (token bucket).
   - Le rate limiter doit être configurable (requêtes par seconde, par minute, capacité du seau).

2. Intégration du rate limiting dans le client :
   - Modifiez la structure `Client` dans `client.go` pour inclure le rate limiter.
   - Ajoutez des options de configuration pour le rate limiting lors de la création du client.

3. Application du rate limiting aux méthodes du client :
   - Modifiez toutes les méthodes du client qui font des appels API pour utiliser le rate limiter avant chaque requête.
   - Gérez les cas où le rate limit est atteint (attente ou erreur selon la configuration).

4. Gestion des erreurs spécifiques au rate limiting :
   - Créez un nouveau type d'erreur `RateLimitError` dans `errors.go`.
   - Implémentez la logique pour retourner cette erreur lorsque le rate limit est atteint et que l'attente n'est pas configurée.

5. Logging des événements liés au rate limiting :
   - Utilisez le système de logging existant pour enregistrer les événements importants liés au rate limiting (attente, erreurs, etc.).

6. Tests unitaires :
   - Créez un fichier `ratelimit_test.go` pour tester le fonctionnement du rate limiter.
   - Ajoutez des tests pour vérifier le comportement du client avec le rate limiting activé.
   - Testez différentes configurations et scénarios (respect des limites, dépassement, attente, erreurs).

7. Documentation :
   - Ajoutez des commentaires détaillés pour toutes les nouvelles fonctions et structures liées au rate limiting.
   - Mettez à jour le README.md avec des informations sur la configuration et l'utilisation du rate limiting.
   - Créez un exemple d'utilisation du rate limiting dans le dossier `examples/`.

8. Configuration flexible :
   - Permettez la configuration de limites différentes pour différents types de requêtes si nécessaire.
   - Ajoutez une option pour désactiver complètement le rate limiting si l'utilisateur le souhaite.

Directives importantes :
- Suivez les conventions de nommage et de formatage standard de Go.
- Assurez-vous que le code est bien commenté et documenté (format godoc).
- Respectez les limites de 500 lignes par fichier et 50 lignes par fonction.
- Le code doit être en anglais, y compris les commentaires et la documentation.
- Assurez-vous que l'implémentation du rate limiting est thread-safe.
- Veillez à ce que le rate limiting n'impacte pas significativement les performances pour les utilisations normales.

Après avoir terminé ces tâches, fournissez un rapport détaillé incluant :
- Un résumé des actions effectuées
- Le contenu des nouveaux fichiers et des modifications majeures
- Les décisions de conception prises pour l'implémentation du rate limiting
- Toute difficulté rencontrée lors de l'implémentation
- Des suggestions pour améliorer ou étendre la fonctionnalité de rate limiting
