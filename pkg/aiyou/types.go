// Package aiyou defines types used across the AI.YOU API client.
package aiyou

import (
	"context"
	"time"
)

// Authenticator manages authentication for API requests.
type Authenticator interface {
	Authenticate(ctx context.Context) error
	Token() string
}

// LoginRequest represents the login request payload.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User represents a user in the AI.YOU system.
type User struct {
	ID           int    `json:"id"` // Changé de string à int
	Email        string `json:"email"`
	ProfileImage string `json:"profileImage"`
	FirstName    string `json:"firstName"`
}

// LoginResponse represents the login response from the API.
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      User      `json:"user"`
}

// ChatCompletionRequest represents a request to the chat completions endpoint.
type ChatCompletionRequest struct {
	Messages     []Message `json:"messages"`               // Liste des messages de la conversation
	AssistantID  string    `json:"assistantId"`            // ID de l'assistant à utiliser
	Temperature  float64   `json:"temperature"`            // Contrôle de la créativité (1-10)
	TopP         float64   `json:"top_p"`                  // Contrôle de la diversité des réponses
	Stream       bool      `json:"stream"`                 // Activer le mode streaming
	PromptSystem string    `json:"promptSystem,omitempty"` // Message système personnalisé
	Form         string    `json:"form,omitempty"`         // Format de sortie
	Stop         []string  `json:"stop,omitempty"`         // Séquences d'arrêt
	ThreadId     string    `json:"threadId,omitempty"`     // ID du thread de conversation
	MaxTokens    *int      `json:"max_tokens,omitempty"`   // Nombre maximum de tokens
}

// Message represents a single message in the conversation.
type Message struct {
	Role    string        `json:"role"`    // Rôle du message (user, assistant, system)
	Content []ContentPart `json:"content"` // Contenu du message
}

// ContentPart represents a part of the message content.
type ContentPart struct {
	Type string `json:"type"` // Type de contenu (text, image, etc.)
	Text string `json:"text"` // Texte du contenu
}

// ChatCompletionResponse represents the response from the chat completions endpoint.
type ChatCompletionResponse struct {
	ID      string   `json:"id"`              // Identifiant unique de la réponse
	Object  string   `json:"object"`          // Type d'objet retourné
	Created int64    `json:"created"`         // Timestamp de création
	Model   string   `json:"model"`           // Modèle utilisé
	Choices []Choice `json:"choices"`         // Liste des choix de réponse
	Usage   *Usage   `json:"usage,omitempty"` // Statistiques d'utilisation
}

// Choice represents a single completion choice in the response.
type Choice struct {
	Index        int     `json:"index"`                   // Index du choix
	Message      Message `json:"message,omitempty"`       // Message pour mode non-streaming
	Delta        *Delta  `json:"delta,omitempty"`         // Delta pour mode streaming
	FinishReason string  `json:"finish_reason,omitempty"` // Raison de fin de génération
}

// Delta represents a streaming response delta.
type Delta struct {
	Role    string `json:"role,omitempty"`    // Rôle du message
	Content string `json:"content,omitempty"` // Contenu du message
}

// Usage represents token usage information.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`     // Nombre de tokens dans la requête
	CompletionTokens int `json:"completion_tokens"` // Nombre de tokens dans la réponse
	TotalTokens      int `json:"total_tokens"`      // Nombre total de tokens
}

// AssistantsResponse represents the response from the assistants endpoint.
type AssistantsResponse struct {
	Context    string      `json:"@context"`
	ID         string      `json:"@id"`
	Type       string      `json:"@type"`
	TotalItems int         `json:"hydra:totalItems"`
	Members    []Assistant `json:"hydra:member"`
}

// Assistant represents an AI assistant.
type Assistant struct {
	ID              string          `json:"id"`
	AssistantID     string          `json:"assistantId"`
	Name            string          `json:"name"`
	ThreadHistories []ThreadHistory `json:"threadHistories"`
	Image           string          `json:"image"`
	Model           string          `json:"model"`
	ModelAi         string          `json:"modelAi"`
	Instructions    string          `json:"instructions"`
}

// ThreadHistory represents the history of a conversation thread.
type ThreadHistory struct {
	ThreadID     string `json:"threadId"`
	FirstMessage string `json:"firstMessage"`
}

// SaveConversationRequest represents a request to save a conversation.
type SaveConversationRequest struct {
	AssistantID    string `json:"assistantId"`
	Conversation   string `json:"conversation"`
	ThreadID       string `json:"threadId,omitempty"`
	FirstMessage   string `json:"firstMessage"`
	ContentJson    string `json:"contentJson"`
	ModelName      string `json:"modelName"`
	IsNewAppThread bool   `json:"isNewAppThread"`
}

// SaveConversationResponse represents the response after saving a conversation.
type SaveConversationResponse struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	CreatedAt int64  `json:"createdAt"`
}

// ConversationThread represents a conversation thread.
type ConversationThread struct {
	ID                   string    `json:"id"`
	ThreadIdParam        int       `json:"threadIdParam"`
	Content              string    `json:"content"`
	AssistantContentJson string    `json:"assistantContentJson"`
	AssistantName        string    `json:"assistantName"`
	AssistantModel       *string   `json:"assistantModel"`
	AssistantId          int       `json:"assistantId"`
	AssistantIdOpenAi    string    `json:"assistantIdOpenAi"`
	FirstMessage         string    `json:"firstMessage"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	IsNewAppThread       bool      `json:"isNewAppThread"`
}

// UserThreadsParams represents the parameters for retrieving user threads.
type UserThreadsParams struct {
	Page         int    `json:"page,omitempty"`
	ItemsPerPage int    `json:"itemsPerPage,omitempty"`
	Search       string `json:"search,omitempty"`
}

// UserThreadsOutput represents the response containing user threads.
type UserThreadsOutput struct {
	Threads      []ConversationThread `json:"threads"`
	TotalItems   int                  `json:"totalItems"`
	ItemsPerPage int                  `json:"itemsPerPage"`
	CurrentPage  int                  `json:"currentPage"`
}

// AudioTranscriptionRequest represents an audio transcription request.
type AudioTranscriptionRequest struct {
	FileName string `json:"fileName"`           // Nom du fichier audio
	Language string `json:"language,omitempty"` // Code langue ISO (ex: "fr", "en")
	Format   string `json:"format,omitempty"`   // Format de sortie souhaité
}

// AudioTranscriptionResponse represents an audio transcription response.
type AudioTranscriptionResponse struct {
	Transcription string `json:"transcription"` // Texte transcrit
}

// SupportedAudioFormat represents a supported audio format.
type SupportedAudioFormat struct {
	Extension string   `json:"extension"` // Extension du fichier (.mp3, .wav, etc.)
	MimeTypes []string `json:"mimeTypes"` // Types MIME supportés
	MaxSize   int64    `json:"maxSize"`   // Taille maximale en bytes
}

// Model represents an AI model in the system.
type Model struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Version     string          `json:"version"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	Properties  ModelProperties `json:"properties"`
}

// ModelProperties represents the properties of an AI model.
type ModelProperties struct {
	MaxTokens    int      `json:"maxTokens"`
	Temperature  float64  `json:"temperature"`
	Provider     string   `json:"provider"`
	Capabilities []string `json:"capabilities"`
}

// ModelRequest represents a request to create a new model.
type ModelRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Properties  ModelProperties `json:"properties"`
}

// ModelResponse represents the response when creating or retrieving a model.
type ModelResponse struct {
	Model Model `json:"model"`
}

// ModelsResponse represents the response when retrieving multiple models.
type ModelsResponse struct {
	Models []Model `json:"models"`
	Total  int     `json:"total"`
}

// StreamOptions represents options for streaming responses.
type StreamOptions struct {
	IncludeUsage         bool `json:"include_usage"`
	ContinuousUsageStats bool `json:"continuous_usage_stats"`
}

// ThreadFilter represents filter options for retrieving threads.
type ThreadFilter struct {
	AssistantID string    `json:"assistantId,omitempty"`
	StartDate   time.Time `json:"startDate,omitempty"`
	EndDate     time.Time `json:"endDate,omitempty"`
	Page        int       `json:"page,omitempty"`
	Limit       int       `json:"limit,omitempty"`
}
