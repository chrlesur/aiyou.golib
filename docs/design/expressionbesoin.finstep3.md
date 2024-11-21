# Expression de Besoin : Package Go pour l'API AI.YOU de Cloud Temple

## 1. Objectif

Développer un package Go nommé "aiyou.golib" qui implémente l'accès à l'API AI.YOU de Cloud Temple, permettant une interaction fluide et typée avec tous les endpoints décrits dans la spécification OpenAPI fournie.

## 2. Fonctionnalités principales

- Authentification basée sur email/mot de passe avec gestion des tokens JWT
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
- Système de rate limiting côté client

## 3. Structure du package

- Client principal avec des méthodes pour chaque endpoint
- Types Go correspondant aux schémas définis dans l'API
- Fonctions utilitaires pour la gestion des requêtes et des réponses
- Option de configuration pour l'URL de base de l'API
- Système d'authentification modulaire basé sur une interface

## 4. Authentification

- Utilisation d'une interface `Authenticator` pour permettre différentes méthodes d'authentification
- Implémentation d'un `JWTAuthenticator` utilisant email/mot de passe
- Gestion automatique du rafraîchissement des tokens JWT
- Stockage sécurisé des tokens

## 5. Gestion des erreurs

- Création de types d'erreur spécifiques pour chaque type de réponse d'erreur de l'API
- Gestion appropriée des codes de statut HTTP
- Utilisation des pratiques standard Go pour la gestion des erreurs
- Pas d'utilisation de `panic`, toutes les erreurs doivent être retournées

## 6. Logging

- Utilisation de l'interface standard `log.Logger` de Go
- Option pour que l'utilisateur puisse fournir son propre logger
- Logging des événements importants (début/fin des appels API, erreurs) avec des niveaux de log appropriés
- Éviter de logger des informations sensibles (tokens, mots de passe)

## 7. Tests

- Tests unitaires pour chaque composant du package
- Tests d'intégration simulant les interactions avec l'API
- Couverture de code visée : au moins 80%
- Utilisation de mocks pour simuler les réponses de l'API dans les tests

## 8. Documentation

- Documentation GoDoc complète pour chaque fonction exportée
- README détaillé avec guide de démarrage rapide, exemples d'utilisation et configuration
- Exemples de code pour chaque fonctionnalité majeure dans un dossier `examples/`
- Guide d'utilisation détaillé dans le dossier `docs/`

## 9. Bonnes pratiques de développement

- Limiter chaque fichier de code source Go à un maximum de 500 lignes
- Limiter chaque fonction à un maximum de 50 lignes de code
- Suivre les conventions de nommage et de formatage standard de Go
- Utiliser `gofmt` pour formater le code
- Utiliser `golint` et `go vet` pour la vérification du code
- Utiliser des interfaces pour les composants principaux pour faciliter les tests et l'extensibilité future

## 10. Rate Limiting

- Implémentation d'un système de rate limiting côté client
- Configuration flexible des limites (requêtes par seconde, par minute, etc.)
- Options pour personnaliser le comportement en cas de dépassement des limites

## 11. Streaming

- Support du streaming pour les opérations de chat completion
- Option pour activer/désactiver le streaming par appel

## 12. Configuration

- Configuration flexible du client (URL de base, timeout, retry, etc.)
- Options pour personnaliser le comportement de l'authentification et du rate limiting

## 13. Priorités de développement

1. Finalisation de l'authentification et gestion des tokens JWT
2. Implémentation de l'endpoint de chat completion (avec et sans streaming)
3. Développement des autres endpoints de l'API
4. Mise en place du système de rate limiting
5. Amélioration de la gestion des erreurs et du logging
6. Développement des tests unitaires et d'intégration
7. Documentation complète et exemples

## 14. Considérations futures

- Implémentation d'un système de cache pour certaines requêtes
- Support pour l'annulation de requêtes via `context.Context`
- Métriques et observabilité avancées
- Support pour les opérations concurrentes et la gestion de connexions multiples

## 15. Licence

Ce projet est sous licence GNU General Public License v3.0 (GPL-3.0).
