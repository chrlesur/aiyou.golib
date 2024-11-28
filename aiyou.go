// Package aiyou provides a client for interacting with the AI.YOU API from Cloud Temple.
package aiyou

import (
	"context"
	"io"
	"time"

	internal "github.com/chrlesur/aiyou.golib/pkg/aiyou"
)

// Réexportation des types publics
type (
	// Client et options
	Client            = internal.Client
	ClientOption      = internal.ClientOption
	MessageBuilder    = internal.MessageBuilder
	StreamReader      = internal.StreamReader
	RateLimiter       = internal.RateLimiter
	RateLimiterConfig = internal.RateLimiterConfig

	// Interfaces
	Authenticator = internal.Authenticator
	Logger        = internal.Logger
	ChatCompleter = internal.ChatCompleter

	// Structures de données Input/Output
	ChatCompletionInput  = internal.ChatCompletionInput
	ChatCompletionOutput = internal.ChatCompletionOutput

	// Structures de requêtes/réponses
	LoginRequest               = internal.LoginRequest
	LoginResponse              = internal.LoginResponse
	User                       = internal.User
	Message                    = internal.Message
	ContentPart                = internal.ContentPart
	ChatCompletionRequest      = internal.ChatCompletionRequest
	ChatCompletionResponse     = internal.ChatCompletionResponse
	Usage                      = internal.Usage
	Choice                     = internal.Choice
	Assistant                  = internal.Assistant
	AssistantsResponse         = internal.AssistantsResponse
	Model                      = internal.Model
	ModelProperties            = internal.ModelProperties
	ModelRequest               = internal.ModelRequest
	ModelResponse              = internal.ModelResponse
	ModelsResponse             = internal.ModelsResponse
	ConversationThread         = internal.ConversationThread
	SaveConversationRequest    = internal.SaveConversationRequest
	SaveConversationResponse   = internal.SaveConversationResponse
	ThreadsResponse            = internal.ThreadsResponse
	ThreadFilter               = internal.ThreadFilter
	AudioTranscriptionRequest  = internal.AudioTranscriptionRequest
	AudioTranscriptionResponse = internal.AudioTranscriptionResponse
	SupportedAudioFormat       = internal.SupportedAudioFormat

	// Types d'erreurs personnalisées
	APIError            = internal.APIError
	AuthenticationError = internal.AuthenticationError
	RateLimitError      = internal.RateLimitError
	NetworkError        = internal.NetworkError

	// Types de log
	LogLevel = internal.LogLevel
)

// Réexportation des constantes
const (
	DEBUG = internal.DEBUG
	INFO  = internal.INFO
	WARN  = internal.WARN
	ERROR = internal.ERROR
)

// Réexportation des variables
var SupportedFormats = internal.SupportedFormats

// Réexportation des fonctions de construction
func NewClient(email, password string, options ...ClientOption) (*Client, error) {
	return internal.NewClient(email, password, options...)
}

func NewDefaultLogger(w io.Writer) Logger {
	return internal.NewDefaultLogger(w)
}

func NewMessageBuilder(role string, logger Logger) *MessageBuilder {
	return internal.NewMessageBuilder(role, logger)
}

func NewRateLimiter(config RateLimiterConfig, logger Logger) *RateLimiter {
	return internal.NewRateLimiter(config, logger)
}

// Réexportation des fonctions utilitaires
func NewTextMessage(role, text string) Message {
	return internal.NewTextMessage(role, text)
}

func NewImageMessage(role, imageURL string) Message {
	return internal.NewImageMessage(role, imageURL)
}

func MaskSensitiveInfo(input string) string {
	return internal.MaskSensitiveInfo(input)
}

func SafeLog(logger Logger) func(level LogLevel, format string, args ...interface{}) {
	return internal.SafeLog(logger)
}

// Réexportation des options de configuration client
func WithLogger(logger Logger) ClientOption {
	return internal.WithLogger(logger)
}

func WithBaseURL(url string) ClientOption {
	return internal.WithBaseURL(url)
}

func WithRateLimiter(config RateLimiterConfig) ClientOption {
	return internal.WithRateLimiter(config)
}

func WithRetry(maxRetries int, initialDelay time.Duration) ClientOption {
	return internal.WithRetry(maxRetries, initialDelay)
}

// Extensions des interfaces pour les méthodes importantes du Client
type ClientInterface interface {
	SetBaseURL(url string)
	SetLogger(logger Logger)
	CreateChatCompletion(ctx context.Context, messages []Message, assistantID string) (*ChatCompletionResponse, error)
	CreateChatCompletionStream(ctx context.Context, messages []Message, assistantID string) (*StreamReader, error)
	ChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error)
	ChatCompletionStream(ctx context.Context, req ChatCompletionRequest) (*StreamReader, error)
	GetUserAssistants(ctx context.Context) (*AssistantsResponse, error)
	TranscribeAudioFile(ctx context.Context, filePath string, opts *AudioTranscriptionRequest) (*AudioTranscriptionResponse, error)
	SaveConversation(ctx context.Context, req SaveConversationRequest) (*SaveConversationResponse, error)
	GetConversation(ctx context.Context, threadID string) (*ConversationThread, error)
	CreateModel(ctx context.Context, req ModelRequest) (*ModelResponse, error)
	GetModels(ctx context.Context) (*ModelsResponse, error)
	GetUserThreads(ctx context.Context, filter *ThreadFilter) (*ThreadsResponse, error)
	DeleteThread(ctx context.Context, threadID string) error
}

// Vérification à la compilation que Client implémente ClientInterface
var _ ClientInterface = (*Client)(nil)
