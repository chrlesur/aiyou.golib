/*
Copyright (C) 2023 Cloud Temple

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

// ChatCompleter gère les requêtes de complétion de chat.
type ChatCompleter interface {
    Complete(ctx context.Context, input ChatCompletionInput) (*ChatCompletionOutput, error)
}

// ChatCompletionInput représente l'entrée pour une requête de complétion de chat.
type ChatCompletionInput struct {
    Messages []Message
    // Autres champs à ajouter selon la spécification de l'API
}

// ChatCompletionOutput représente la sortie d'une requête de complétion de chat.
type ChatCompletionOutput struct {
    Response string
    // Autres champs à ajouter selon la spécification de l'API
}

// Message représente un message dans une conversation.
type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

// ChatCompletionRequest represents the request structure for chat completion
type ChatCompletionRequest struct {
    Messages      []Message `json:"messages"`
    AssistantID   string    `json:"assistantId"`
    Temperature   float32   `json:"temperature"`
    TopP          float32   `json:"top_p"`
    PromptSystem  string    `json:"promptSystem"`
    Stream        bool      `json:"stream"`
    Stop          []string  `json:"stop,omitempty"`
    ThreadID      string    `json:"threadId,omitempty"`
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

// ChatCompletionResponse represents the response structure for chat completion
type ChatCompletionResponse struct {
    ID      string    `json:"id"`
    Object  string    `json:"object"`
    Created int64     `json:"created"`
    Model   string    `json:"model"`
    Usage   Usage     `json:"usage"`
    Choices []Choice  `json:"choices"`
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
