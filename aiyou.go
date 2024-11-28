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

	// Structures de messages et contenus
	Message     = internal.Message
	ContentPart = internal.ContentPart

	// Structures d'authentification
	LoginRequest  = internal.LoginRequest
	LoginResponse = internal.LoginResponse
	User          = internal.User

	// Structures de chat completion
	ChatCompletionRequest  = internal.ChatCompletionRequest
	ChatCompletionResponse = internal.ChatCompletionResponse
	Usage                  = internal.Usage
	Choice                 = internal.Choice

	// Structures des assistants
	Assistant     = internal.Assistant
	ThreadHistory = internal.ThreadHistory

	// Structures des modèles
	Model           = internal.Model
	ModelProperties = internal.ModelProperties
	ModelRequest    = internal.ModelRequest
	ModelResponse   = internal.ModelResponse
	ModelsResponse  = internal.ModelsResponse

	// Structures des conversations et threads
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

// Fonctions de construction
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

// Fonctions utilitaires
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

// Options de configuration client
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

// Interface du Client
type ClientInterface interface {
	SetBaseURL(url string)
	SetLogger(logger Logger)
	CreateChatCompletion(ctx context.Context, messages []Message, assistantID string) (*ChatCompletionResponse, error)
	CreateChatCompletionStream(ctx context.Context, messages []Message, assistantID string) (*StreamReader, error)
	ChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error)
	ChatCompletionStream(ctx context.Context, req ChatCompletionRequest) (*StreamReader, error)
	GetUserAssistants(ctx context.Context) ([]Assistant, error)
	TranscribeAudioFile(ctx context.Context, filePath string, opts *AudioTranscriptionRequest) (*AudioTranscriptionResponse, error)
	SaveConversation(ctx context.Context, req SaveConversationRequest) (*SaveConversationResponse, error)
	GetConversation(ctx context.Context, threadID string) (*ConversationThread, error)
	CreateModel(ctx context.Context, req ModelRequest) (*ModelResponse, error)
	GetModels(ctx context.Context) (*ModelsResponse, error)
	GetUserThreads(ctx context.Context, params *UserThreadsParams) (*UserThreadsOutput, error)
	DeleteThread(ctx context.Context, threadID string) error
}

// Vérification à la compilation que Client implémente ClientInterface
var _ ClientInterface = (*Client)(nil)
