
# Expression de besoin pour l'API d'Intelligence Artificielle de SHIVA

## 1. Objectif

Créer une API d'IA robuste pour SHIVA, s'inspirant de l'API AiYou existante  avec un accent sur la sécurité, la flexibilité et la conformité aux normes SecNumCloud.

## 2. Fonctionnalités principales

### 2.1 Gestion des modèles et des assistants

- Accès unifié aux modèles de différents fournisseurs :
  - OpenAI (o1, GPT-4)
  - Anthropic (Claude)
  - Groq
  - SecNumCloud (Nemotron 70B sur infrastructure vLLM en zone sécurisée); réinternalisation de l'infra NEMO dans la plateforme
- Création et gestion d'assistants personnalisés avec des instructions spécifiques

### 2.2 Chat Completions

- Génération de réponses basées sur un historique de messages
- Support du streaming pour des réponses en temps réel
- Paramètres de contrôle : température, top_p, etc.

### 2.3 Traitement audio

- Transcription audio en texte

### 2.4 Gestion des threads de conversation

- Sauvegarde et récupération des historiques de conversation

## 3. Structure de l'API

### 3.1 Chat Completions

- `POST /api/v1/chat/completions`
  - Corps de la requête :
    - `messages` : Array des messages de la conversation
    - `assistantId` : ID de l'assistant à utiliser
    - `temperature`, `top_p` : Paramètres de contrôle
    - `promptSystem` : Instructions système pour l'assistant
    - `stream` : Booléen pour activer le streaming
    - `threadId` : ID du thread de conversation (optionnel)

### 3.2 Gestion des modèles

- `GET /api/v1/models` : Lister tous les modèles disponibles

### 3.3 Gestion des assistants

- `GET /api/v1/user/assistants` : Récupération des assistants de l'utilisateur
- `POST /api/v1/assistants` : Création d'un nouvel assistant
- `PUT /api/v1/assistants/{assistant_id}` : Modification d'un assistant existant
- `DELETE /api/v1/assistants/{assistant_id}` : Suppression d'un assistant

### 3.4 Gestion des threads

- `POST /api/v1/save` : Sauvegarde d'une conversation
- `GET /api/v1/user/threads` : Récupération des threads de l'utilisateur

### 3.5 Transcription audio

- `POST /api/v1/audio/transcriptions` : Transcription d'un fichier audio en texte ( avec whisper chez openai puis en SNC)

## 4. Sécurité et authentification

### 4.1 Mécanisme d'authentification

- Utilisation du système d'authentification existant de SHIVA
- Authentification basée sur les tokens JWT (JSON Web Tokens)

### 4.2 Endpoints d'authentification

- Réutilisation des endpoints d'authentification de SHIVA :
  - `POST /v2/auth/signin` : Pour l'authentification initiale
  - `POST /v2/auth/refresh` : Pour rafraîchir le token d'authentification
  - `POST /v2/auth/logout` : Pour la déconnexion

### 4.3 Sécurité des requêtes

- Toutes les requêtes à l'API IA doivent inclure un token JWT valide dans l'en-tête d'autorisation
- Utilisation du schéma "Bearer" pour l'authentification

### 4.4 Gestion des permissions

- Intégration avec le système de rôles et permissions de SHIVA
- Définition de permissions spécifiques pour l'accès aux fonctionnalités IA :
  - `ai_read` : Pour l'accès en lecture aux modèles et assistants
  - `ai_write` : Pour la création et modification d'assistants
  - `ai_execute` : Pour l'exécution des modèles et interactions avec les assistants

### 4.5 Sécurité supplémentaire

- Chiffrement de bout en bout pour toutes les communications
- Mise en place de limites de taux (rate limiting) pour prévenir les abus
- Journalisation détaillée de toutes les activités d'authentification et d'autorisation

## 5. Intégration et compatibilité

- Maintien de la structure de l'API existante pour une transition en douceur
- Extension pour supporter les nouveaux fournisseurs de modèles
- SDKs pour langages populaires (Python, JavaScript, Go)

## 6. Performance et scalabilité

- Optimisation pour des réponses rapides, même avec un grand nombre d'utilisateurs (cible a 30 users // pour commencer )
- Support du streaming pour les réponses en temps réel
- Architecture de microservices pour une scalabilité optimale
- Équilibrage de charge intelligent entre les différents architectures de calcul en SNC (LB devant les lames ou les mac mini)

## 7. Monitoring et logging

- Implémentation de métriques détaillées sur l'utilisation de l'API (nb de tokens input / output / total par users / par tenant / par compagnie)
- Logging des erreurs et des événements importants
- Tableaux de bord en temps réel pour l'utilisation et les performances
- Alertes configurables sur des seuils d'utilisation ou de performance

## 8. Documentation et support

- Documentation complète de l'API, incluant des exemples d'utilisation
- Mise en place d'un système de versionnage clair pour l'API
- Environnement sandbox pour tests et développement
- Support technique avec SLAs garantis

## 9. Considérations éthiques et légales

- Implémentation de mécanismes pour prévenir les utilisations abusives de l'IA
- Conformité avec les réglementations sur la protection des données (ex: RGPD)
- Mécanismes de détection et d'atténuation des biais
- Lignes directrices pour une utilisation responsable de l'IA

## 10. Évolutivité

- Conception modulaire pour faciliter l'ajout de nouveaux fournisseurs de modèles
- Préparation pour l'intégration future de nouvelles fonctionnalités d'IA

## 11. Tarification et gestion des coûts

- Modèle de tarification flexible basé sur le nombre de token total
- Outils de suivi et d'optimisation des coûts
- Quotas et limites d'utilisation configurables par utilisateur ou par projet
