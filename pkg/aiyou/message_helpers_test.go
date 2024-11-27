/*
Copyright (C) 2024 Cloud Temple

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
*/

package aiyou

import (
	"io"
	"testing"
)

func TestMessageBuilder(t *testing.T) {
	logger := NewDefaultLogger(io.Discard)

	t.Run("Build text only message", func(t *testing.T) {
		builder := NewMessageBuilder("user", logger)
		msg := builder.
			AddText("Hello").
			Build()

		if msg.Role != "user" {
			t.Errorf("Expected role 'user', got %s", msg.Role)
		}
		if len(msg.Content) != 1 {
			t.Errorf("Expected 1 content part, got %d", len(msg.Content))
		}
		if msg.Content[0].Type != "text" {
			t.Errorf("Expected type 'text', got %s", msg.Content[0].Type)
		}
	})

	t.Run("Build complex message", func(t *testing.T) {
		builder := NewMessageBuilder("assistant", logger)
		msg := builder.
			AddText("Here's an image").
			AddImage("https://example.com/image.jpg").
			AddText("What do you think?").
			Build()

		if len(msg.Content) != 3 {
			t.Errorf("Expected 3 content parts, got %d", len(msg.Content))
		}

		expectedTypes := []string{"text", "image", "text"}
		for i, expectedType := range expectedTypes {
			if msg.Content[i].Type != expectedType {
				t.Errorf("Content part %d: expected type %s, got %s",
					i, expectedType, msg.Content[i].Type)
			}
		}
	})
}

func TestHelperFunctions(t *testing.T) {
	t.Run("NewTextMessage", func(t *testing.T) {
		msg := NewTextMessage("user", "Hello")
		if msg.Role != "user" {
			t.Errorf("Expected role 'user', got %s", msg.Role)
		}
		if len(msg.Content) != 1 {
			t.Errorf("Expected 1 content part, got %d", len(msg.Content))
		}
		if msg.Content[0].Type != "text" {
			t.Errorf("Expected type 'text', got %s", msg.Content[0].Type)
		}
	})

	t.Run("NewImageMessage", func(t *testing.T) {
		msg := NewImageMessage("user", "https://example.com/image.jpg")
		if msg.Role != "user" {
			t.Errorf("Expected role 'user', got %s", msg.Role)
		}
		if len(msg.Content) != 1 {
			t.Errorf("Expected 1 content part, got %d", len(msg.Content))
		}
		if msg.Content[0].Type != "image" {
			t.Errorf("Expected type 'image', got %s", msg.Content[0].Type)
		}
	})
}
