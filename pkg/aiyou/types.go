/*
Copyright (C) 2024 Cloud Temple

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

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

// LoginResponse represents the login response from the API.
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      User      `json:"user"`
}

// User represents a user in the AI.YOU system.
type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	ProfileImage string `json:"profileImage"`
	FirstName    string `json:"firstName"`
}

// Gardons uniquement ces structures pour le chat
type ChatCompletionRequest struct {
	Messages     []Message `json:"messages"`
	AssistantID  string    `json:"assistantId"`
	Temperature  float32   `json:"temperature"`
	TopP         float32   `json:"top_p"`
	PromptSystem string    `json:"promptSystem"`
	Stream       bool      `json:"stream"`
	Stop         []string  `json:"stop,omitempty"`
	ThreadID     string    `json:"threadId,omitempty"`
}

type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Usage   Usage    `json:"usage"`
	Choices []Choice `json:"choices"`
}

// Message represents a single message in the chat completion request
type Message struct {
	Role    string        `json:"role"`
	Content []ContentPart `json:"content"`
}

// ContentPart represents a part of the message content
type ContentPart struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Usage represents the token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Choice represents a single choice in the chat completion response
type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// AssistantsResponse représente la réponse de l'API pour la liste des assistants
type AssistantsResponse struct {
	Context    string      `json:"@context"`
	ID         string      `json:"@id"`
	Type       string      `json:"@type"`
	TotalItems int         `json:"hydra:totalItems"`
	Members    []Assistant `json:"hydra:member"`
}

// Assistant représente un assistant dans le système AI.YOU
type Assistant struct {
	ID              string          `json:"id"`
	AssistantID     string          `json:"assistantId"`
	Name            string          `json:"name"`
	ThreadHistories []ThreadHistory `json:"threadHistories"`
	Image           string          `json:"image"`
	Model           string          `json:"model"`   // Ajouté
	ModelAi         string          `json:"modelAi"` // Ajouté
	Instructions    string          `json:"instructions"`
}

type ThreadHistory struct {
	ThreadID     string `json:"threadId"`
	FirstMessage string `json:"firstMessage"`
}

// Model représente un modèle dans le système AI.YOU
type Model struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Version     string          `json:"version"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	Properties  ModelProperties `json:"properties"`
}

// ModelProperties représente les propriétés spécifiques d'un modèle
type ModelProperties struct {
	MaxTokens    int      `json:"maxTokens"`
	Temperature  float64  `json:"temperature"`
	Provider     string   `json:"provider"`
	Capabilities []string `json:"capabilities"`
}

// ModelRequest représente la requête de création d'un modèle
type ModelRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Properties  ModelProperties `json:"properties"`
}

// ModelResponse représente la réponse de l'API pour les opérations sur les modèles
type ModelResponse struct {
	Model Model `json:"model"`
}

// ModelsResponse représente la réponse pour la liste des modèles
type ModelsResponse struct {
	Models []Model `json:"models"`
	Total  int     `json:"total"`
}

// ConversationThread représente un fil de conversation
type ConversationThread struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Messages      []Message `json:"messages"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	AssistantID   string    `json:"assistantId"` // Back to string
	UserID        string    `json:"userId"`      // Back to string
	LastMessageAt time.Time `json:"lastMessageAt"`
}

// UserThreadsOutput représente la réponse de l'API pour la liste des threads
type UserThreadsOutput struct {
	Threads      []ConversationThread `json:"threads"`
	TotalItems   int                  `json:"totalItems"`
	ItemsPerPage int                  `json:"itemsPerPage"`
	CurrentPage  int                  `json:"currentPage"`
}

type UserThreadsParams struct {
	Page         int    `json:"page,omitempty"`
	ItemsPerPage int    `json:"itemsPerPage,omitempty"`
	Search       string `json:"search,omitempty"`
}

// SaveConversationRequest représente la requête pour sauvegarder une conversation
type SaveConversationRequest struct {
	AssistantID    string `json:"assistantId"`
	Conversation   string `json:"conversation"`
	ThreadID       string `json:"threadId,omitempty"`
	FirstMessage   string `json:"firstMessage"`
	ContentJson    string `json:"contentJson"`
	ModelName      string `json:"modelName"`
	IsNewAppThread bool   `json:"isNewAppThread"`
}

// SaveConversationResponse représente la réponse après sauvegarde d'une conversation
type SaveConversationResponse struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	CreatedAt int64  `json:"createdAt"`
}

// ThreadsResponse représente la réponse de l'API pour la liste des threads
type ThreadsResponse struct {
	Threads []ConversationThread `json:"threads"`
	Total   int                  `json:"total"`
	Page    int                  `json:"page"`
	Limit   int                  `json:"limit"`
}

// ThreadFilter représente les options de filtrage pour la récupération des threads
type ThreadFilter struct {
	AssistantID string    `json:"assistantId,omitempty"`
	StartDate   time.Time `json:"startDate,omitempty"`
	EndDate     time.Time `json:"endDate,omitempty"`
	Page        int       `json:"page,omitempty"`
	Limit       int       `json:"limit,omitempty"`
}

// AudioTranscriptionRequest représente la requête de transcription audio
type AudioTranscriptionRequest struct {
	FileName string `json:"fileName"`
	Language string `json:"language,omitempty"` // Code langue ISO (ex: "fr", "en")
	Format   string `json:"format,omitempty"`   // Format de sortie souhaité
}

// AudioTranscriptionResponse représente la réponse de transcription
type AudioTranscriptionResponse struct {
	Transcription string `json:"transcription"` // Changé de Text à Transcription
}

// SupportedAudioFormat liste les formats audio supportés
type SupportedAudioFormat struct {
	Extension string   `json:"extension"`
	MimeTypes []string `json:"mimeTypes"`
	MaxSize   int64    `json:"maxSize"` // Taille maximale en bytes
}
