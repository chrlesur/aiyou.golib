package aiyou

import (
	"context"
	"io"
	"os"
	"testing"
)

func TestChatCompletion(t *testing.T) {
	logger := NewDefaultLogger(os.Stderr)
	logger.SetLevel(DEBUG)

	client, err := NewClient(
		WithEmailPassword(testConfig.Email, testConfig.Password),
		WithBaseURL(testConfig.BaseURL),
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Message de test simple
	req := ChatCompletionRequest{
		Messages: []Message{
			{
				Role: "user",
				Content: []ContentPart{
					{Type: "text", Text: "Hello, how are you?"},
				},
			},
		},
		AssistantID: "287", // Utiliser un ID d'assistant valide de votre système
	}

	t.Log("Sending chat completion request...")
	resp, err := client.ChatCompletion(context.Background(), req)
	if err != nil {
		t.Fatalf("ChatCompletion failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected non-nil response")
	}

	// Log de la réponse
	t.Log("Got response:")
	t.Logf("  ID: %s", resp.ID)
	t.Logf("  Model: %s", resp.Model)
	t.Logf("  Created: %d", resp.Created)
	if len(resp.Choices) > 0 {
		t.Logf("  First response: %s", resp.Choices[0].Message.Content[0].Text)
	}
}

func TestChatCompletionStream(t *testing.T) {
	logger := NewDefaultLogger(os.Stderr)
	logger.SetLevel(DEBUG)

	client, err := NewClient(
		WithEmailPassword(testConfig.Email, testConfig.Password),
		WithBaseURL(testConfig.BaseURL),
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req := ChatCompletionRequest{
		Messages: []Message{
			{
				Role: "user",
				Content: []ContentPart{
					{Type: "text", Text: "Count from 1 to 5 slowly"},
				},
			},
		},
		AssistantID: "287", // Utiliser un ID d'assistant valide
		Stream:      true,
	}

	t.Log("Starting streaming chat completion...")
	stream, err := client.ChatCompletionStream(context.Background(), req)
	if err != nil {
		t.Fatalf("ChatCompletionStream failed: %v", err)
	}
	defer stream.Close()

	t.Log("Reading stream chunks:")
	for {
		chunk, err := stream.ReadChunk()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Error reading stream: %v", err)
		}

		if chunk != nil && len(chunk.Choices) > 0 {
			choice := chunk.Choices[0]
			if choice.Delta != nil && choice.Delta.Content != "" {
				t.Logf("Received chunk: %s", choice.Delta.Content)
			}
		}
	}
}

func TestChatCompletionWithComplexMessage(t *testing.T) {
	logger := NewDefaultLogger(os.Stderr)
	logger.SetLevel(DEBUG)

	client, err := NewClient(
		WithEmailPassword(testConfig.Email, testConfig.Password),
		WithBaseURL(testConfig.BaseURL),
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Créer un message plus complexe avec le MessageBuilder
	builder := NewMessageBuilder("user", logger)
	builder.AddText("Explain these three things:\n")
	builder.AddText("1. Go interfaces\n")
	builder.AddText("2. Error handling\n")
	builder.AddText("3. Concurrency patterns")

	req := ChatCompletionRequest{
		Messages:    []Message{builder.Build()},
		AssistantID: "287", // Utiliser un ID d'assistant valide
		Temperature: 7,
		TopP:        1.0,
	}

	t.Log("Sending complex chat completion request...")
	resp, err := client.ChatCompletion(context.Background(), req)
	if err != nil {
		t.Fatalf("ChatCompletion failed: %v", err)
	}

	if len(resp.Choices) > 0 {
		t.Log("Response to complex query:")
		t.Logf("Content: %s", resp.Choices[0].Message.Content[0].Text)
	}
}

func TestChatCompletionUnauthorized(t *testing.T) {
	logger := NewDefaultLogger(os.Stderr)
	logger.SetLevel(DEBUG)

	client, err := NewClient(
		WithEmailPassword("invalid@example.com", "wrong_password"),
		WithBaseURL(testConfig.BaseURL),
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req := ChatCompletionRequest{
		Messages: []Message{
			{
				Role: "user",
				Content: []ContentPart{
					{Type: "text", Text: "Hello"},
				},
			},
		},
		AssistantID: "287",
	}

	_, err = client.ChatCompletion(context.Background(), req)
	if err == nil {
		t.Error("Expected error for unauthorized request")
	} else {
		t.Logf("Got expected error: %v", err)
	}
}

func TestChatCompletionWithInvalidAssistant(t *testing.T) {
	logger := NewDefaultLogger(os.Stderr)
	logger.SetLevel(DEBUG)

	client, err := NewClient(
		WithEmailPassword(testConfig.Email, testConfig.Password),
		WithBaseURL(testConfig.BaseURL),
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	req := ChatCompletionRequest{
		Messages: []Message{
			{
				Role: "user",
				Content: []ContentPart{
					{Type: "text", Text: "Hello"},
				},
			},
		},
		AssistantID: "invalid_assistant_id",
	}

	_, err = client.ChatCompletion(context.Background(), req)
	if err == nil {
		t.Error("Expected error for invalid assistant ID")
	} else {
		t.Logf("Got expected error: %v", err)
	}
}
