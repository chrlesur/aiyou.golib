
# Expression de Besoin : Package Go pour l'API AI.YOU de Cloud Temple (Version 2.0)

## 1. Objectif

Développer un package Go nommé "aiyou.golib" qui implémente l'accès à l'API AI.YOU de Cloud Temple, permettant une interaction fluide et typée avec tous les endpoints décrits dans la spécification OpenAPI fournie.

## 2. Fonctionnalités principales

- Authentification basée sur email/mot de passe avec gestion des tokens JWT
- URL de base de l'API paramétrable
- Implémentation de tous les endpoints décrits dans l'API :
  - Login (`/api/login`) [Implémenté]
  - Chat completions (`/api/v1/chat/completions`) [Implémenté]
  - Création de modèles (`/api/v1/models`)
  - Récupération des assistants de l'utilisateur (`/api/v1/user/assistants`)
  - Sauvegarde d'une conversation (`/api/v1/save`)
  - Récupération des threads de l'utilisateur (`/api/v1/user/threads`)
  - Transcription audio (`/api/v1/audio/transcriptions`)
- Gestion des erreurs et des réponses pour chaque endpoint
- Support pour les opérations en streaming (pour les chat completions)
- Paramétrage complet des modèles (temperature, tokens, etc.) pour chaque appel
- Système de rate limiting côté client
- Système de retry pour les erreurs temporaires de réseau

## 3. Structure du package

- Client principal avec des méthodes pour chaque endpoint
- Types Go correspondant aux schémas définis dans l'API
- Fonctions utilitaires pour la gestion des requêtes et des réponses
- Option de configuration pour l'URL de base de l'API
- Système d'authentification modulaire basé sur une interface
- StreamReader pour gérer les réponses en streaming

## 4. Authentification

- Utilisation d'une interface `Authenticator` pour permettre différentes méthodes d'authentification
- Implémentation d'un `JWTAuthenticator` utilisant email/mot de passe
- Gestion automatique du rafraîchissement des tokens JWT
- Stockage sécurisé des tokens

## 5. Chat Completion

- Implémentation des méthodes `ChatCompletion` et `ChatCompletionStream`
- Support pour les requêtes standard et en streaming
- Structures de données complètes pour les requêtes et les réponses
- Options de configuration avancées (température, top_p, système de prompt, etc.)
- Helpers pour la construction de messages complexes (texte, images, etc.)

## 6. Gestion des erreurs et logging

- Création de types d'erreur spécifiques pour chaque type de réponse d'erreur de l'API
- Implémentation d'un système de retry configurable pour les erreurs temporaires de réseau
- Logging avancé pour faciliter le débogage :
  - Traces détaillées des requêtes et réponses
  - Niveaux de log configurables (debug, info, warn, error)
  - Option pour masquer les informations sensibles dans les logs
- Utilisation des pratiques standard Go pour la gestion des erreurs
- Pas d'utilisation de `panic`, toutes les erreurs doivent être retournées

## 7. Logging

- Utilisation de l'interface standard `log.Logger` de Go
- Option pour que l'utilisateur puisse fournir son propre logger
- Logging des événements importants avec des niveaux de log appropriés
- Logging configurable des requêtes et réponses pour faciliter le débogage

## 8. Tests

- Tests unitaires pour chaque composant du package
- Tests d'intégration simulant les interactions avec l'API
- Utilisation de mock servers pour simuler les réponses de l'API
- Tests spécifiques pour les fonctionnalités de streaming
- Couverture de code visée : au moins 80%

## 9. Documentation

- Limiter chaque fichier de code source Go à un maximum de 500 lignes
- Limiter chaque fonction à un maximum de 50 lignes de code
- Suivre les conventions de nommage et de formatage standard de Go
- Utiliser `gofmt` pour formater le code
- Utiliser `golint` et `go vet` pour la vérification du code
- Utiliser des interfaces pour les composants principaux pour faciliter les tests et l'extensibilité future

## 10. Bonnes pratiques de développement

- Implémentation d'un système de rate limiting côté client
- Configuration flexible des limites (requêtes par seconde, par minute, etc.)
- Options pour personnaliser le comportement en cas de dépassement des limites

## 11. Rate Limiting

- Support du streaming pour les opérations de chat completion
- Option pour activer/désactiver le streaming par appel

## 12. Streaming

- Support robuste du streaming pour les opérations de chat completion
- Implémentation d'un StreamReader pour simplifier la lecture des réponses en streaming
- Gestion efficace des connexions pour le streaming

## 13. Configuration

- Configuration flexible du client (URL de base, timeout, retry, etc.)
- Options pour personnaliser le comportement de l'authentification et du rate limiting
- Configuration avancée pour le chat completion (température, top_p, etc.)

## 14. Priorités de développement révisées

1. Amélioration de la gestion des erreurs et du système de logging
2. Implémentation du système de retry pour les erreurs réseau
3. Développement des helpers pour la construction de messages complexes
4. Implémentation des endpoints restants de l'API
5. Mise en place du système de rate limiting
6. Optimisation des performances et de l'utilisation des ressources
7. Amélioration de la documentation et des exemples

## 15. Considérations futures

- Implémentation d'un système de cache pour certaines requêtes
- Support pour l'annulation de requêtes via `context.Context`
- Métriques et observabilité avancées
- Support pour les opérations concurrentes et la gestion de connexions multiples
- Intégration avec des frameworks de monitoring populaires

## 16. Licence

Ce projet est sous licence GNU General Public License v3.0 (GPL-3.0).
