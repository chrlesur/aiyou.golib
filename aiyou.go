// Package aiyou provides a client for interacting with the AI.YOU API from Cloud Temple.
// It supports both email/password and bearer token authentication methods.
package aiyou

import (
	"context"
	"io"
	"time"

	internal "github.com/chrlesur/aiyou.golib/pkg/aiyou"
)

// Réexportation des types publics
type (
	// Client et options de configuration
	Client            = internal.Client
	ClientOption      = internal.ClientOption
	MessageBuilder    = internal.MessageBuilder
	StreamReader      = internal.StreamReader
	RateLimiter       = internal.RateLimiter
	RateLimiterConfig = internal.RateLimiterConfig

	// Interfaces fondamentales
	Authenticator = internal.Authenticator // Interface pour l'authentification (JWT ou Bearer)
	Logger        = internal.Logger        // Interface pour le logging personnalisé

	// Structures de messages et contenus
	Message       = internal.Message     // Représente un message dans la conversation
	ContentPart   = internal.ContentPart // Partie de contenu d'un message (texte, image, etc.)
	StreamOptions = internal.StreamOptions

	// Structures d'authentification
	LoginRequest  = internal.LoginRequest  // Requête d'authentification par email/mot de passe
	LoginResponse = internal.LoginResponse // Réponse contenant le token JWT
	User          = internal.User

	// Structures de chat completion
	ChatCompletionRequest  = internal.ChatCompletionRequest  // Requête de completion de chat
	ChatCompletionResponse = internal.ChatCompletionResponse // Réponse de completion
	Usage                  = internal.Usage                  // Statistiques d'utilisation
	Choice                 = internal.Choice                 // Choix de réponse
	Delta                  = internal.Delta                  // Pour le streaming

	// Structures des assistants et threads
	Assistant          = internal.Assistant
	ThreadHistory      = internal.ThreadHistory
	AssistantsResponse = internal.AssistantsResponse

	// Structures des modèles
	Model           = internal.Model
	ModelProperties = internal.ModelProperties
	ModelRequest    = internal.ModelRequest
	ModelResponse   = internal.ModelResponse
	ModelsResponse  = internal.ModelsResponse

	// Structures des conversations
	ConversationThread       = internal.ConversationThread
	SaveConversationRequest  = internal.SaveConversationRequest
	SaveConversationResponse = internal.SaveConversationResponse
	UserThreadsOutput        = internal.UserThreadsOutput
	ThreadFilter             = internal.ThreadFilter
	UserThreadsParams        = internal.UserThreadsParams

	// Structures audio
	AudioTranscriptionRequest  = internal.AudioTranscriptionRequest
	AudioTranscriptionResponse = internal.AudioTranscriptionResponse
	SupportedAudioFormat       = internal.SupportedAudioFormat

	// Types d'erreurs personnalisées
	APIError            = internal.APIError            // Erreurs API générales
	AuthenticationError = internal.AuthenticationError // Erreurs d'authentification
	RateLimitError      = internal.RateLimitError      // Erreurs de limitation de débit
	NetworkError        = internal.NetworkError        // Erreurs réseau

	// Types de log
	LogLevel = internal.LogLevel
)

// Constantes de niveau de log
const (
	DEBUG = internal.DEBUG // Niveau de log pour le développement
	INFO  = internal.INFO  // Niveau de log pour les informations générales
	WARN  = internal.WARN  // Niveau de log pour les avertissements
	ERROR = internal.ERROR // Niveau de log pour les erreurs
)

// Variables exportées
var SupportedFormats = internal.SupportedFormats // Formats audio supportés

// NewClient crée un nouveau client AI.YOU
// Supporte deux méthodes d'authentification :
// 1. Email/Mot de passe : NewClient("email@example.com", "password")
// 2. Bearer Token : NewClient("", "", WithBearerToken("your-token"))
func NewClient(options ...ClientOption) (*Client, error) {
	return internal.NewClient(options...)
}

// NewDefaultLogger crée un nouveau logger avec la configuration par défaut
func NewDefaultLogger(w io.Writer) Logger {
	return internal.NewDefaultLogger(w)
}

// NewMessageBuilder crée un nouveau builder pour construire des messages complexes
func NewMessageBuilder(role string, logger Logger) *MessageBuilder {
	return internal.NewMessageBuilder(role, logger)
}

// NewRateLimiter crée un nouveau rate limiter avec la configuration spécifiée
func NewRateLimiter(config RateLimiterConfig, logger Logger) *RateLimiter {
	return internal.NewRateLimiter(config, logger)
}

// Fonctions utilitaires pour la création de messages
func NewTextMessage(role, text string) Message {
	return internal.NewTextMessage(role, text)
}

func NewImageMessage(role, imageURL string) Message {
	return internal.NewImageMessage(role, imageURL)
}

// Fonctions utilitaires de sécurité et logging
func MaskSensitiveInfo(input string) string {
	return internal.MaskSensitiveInfo(input)
}

func SafeLog(logger Logger) func(level LogLevel, format string, args ...interface{}) {
	return internal.SafeLog(logger)
}

// Options de configuration du client
func WithLogger(logger Logger) ClientOption {
	return internal.WithLogger(logger)
}

func WithBaseURL(url string) ClientOption {
	return internal.WithBaseURL(url)
}

// WithBearerToken configure le client pour utiliser l'authentification par bearer token
func WithBearerToken(token string) ClientOption {
	return internal.WithBearerToken(token)
}

// WithEmailPassword configure le client pour utiliser l'authentification par email/mot de passe
func WithEmailPassword(email, password string) ClientOption {
	return internal.WithEmailPassword(email, password)
}

func WithRateLimiter(config RateLimiterConfig) ClientOption {
	return internal.WithRateLimiter(config)
}

func WithRetry(maxRetries int, initialDelay time.Duration) ClientOption {
	return internal.WithRetry(maxRetries, initialDelay)
}

// Interface du Client définissant toutes les opérations disponibles
type ClientInterface interface {
	// Configuration
	SetBaseURL(url string)
	SetLogger(logger Logger)

	// Opérations de chat
	CreateChatCompletion(ctx context.Context, messages []Message, assistantID string) (*ChatCompletionResponse, error)
	CreateChatCompletionStream(ctx context.Context, messages []Message, assistantID string) (*StreamReader, error)
	ChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error)
	ChatCompletionStream(ctx context.Context, req ChatCompletionRequest) (*StreamReader, error)

	// Opérations sur les assistants et modèles
	GetUserAssistants(ctx context.Context) (*AssistantsResponse, error)
	CreateModel(ctx context.Context, req ModelRequest) (*ModelResponse, error)
	GetModels(ctx context.Context) (*ModelsResponse, error)

	// Opérations sur les conversations
	SaveConversation(ctx context.Context, req SaveConversationRequest) (*SaveConversationResponse, error)
	GetConversation(ctx context.Context, threadID string) (*ConversationThread, error)
	GetUserThreads(ctx context.Context, params *UserThreadsParams) (*UserThreadsOutput, error)
	DeleteThread(ctx context.Context, threadID string) error

	// Opérations audio
	TranscribeAudioFile(ctx context.Context, filePath string, opts *AudioTranscriptionRequest) (*AudioTranscriptionResponse, error)
}

// Vérification à la compilation que Client implémente ClientInterface
var _ ClientInterface = (*Client)(nil)
