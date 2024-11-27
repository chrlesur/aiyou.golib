/*
Copyright (C) 2024 Cloud Temple

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
*/

package aiyou

// MessageBuilder helps construct complex messages with multiple content types.
type MessageBuilder struct {
	message Message
	logger  Logger
}

// NewMessageBuilder creates a new MessageBuilder with the specified role.
func NewMessageBuilder(role string, logger Logger) *MessageBuilder {
	return &MessageBuilder{
		message: Message{
			Role:    role,
			Content: make([]ContentPart, 0),
		},
		logger: logger,
	}
}

// AddText adds a text content part to the message.
func (mb *MessageBuilder) AddText(text string) *MessageBuilder {
	mb.logger.Debugf("Adding text content: %s", MaskSensitiveInfo(text))
	mb.message.Content = append(mb.message.Content, ContentPart{
		Type: "text",
		Text: text,
	})
	return mb
}

// AddImage adds an image content part to the message.
func (mb *MessageBuilder) AddImage(imageURL string) *MessageBuilder {
	mb.logger.Debugf("Adding image content: %s", MaskSensitiveInfo(imageURL))
	mb.message.Content = append(mb.message.Content, ContentPart{
		Type: "image",
		Text: imageURL,
	})
	return mb
}

// Build returns the constructed Message.
func (mb *MessageBuilder) Build() Message {
	mb.logger.Infof("Building message with %d content parts", len(mb.message.Content))
	return mb.message
}

// Helper functions for creating simple messages

// NewTextMessage creates a new Message with a single text content part.
func NewTextMessage(role, text string) Message {
	return Message{
		Role: role,
		Content: []ContentPart{
			{
				Type: "text",
				Text: text,
			},
		},
	}
}

// NewImageMessage creates a new Message with a single image content part.
func NewImageMessage(role, imageURL string) Message {
	return Message{
		Role: role,
		Content: []ContentPart{
			{
				Type: "image",
				Text: imageURL,
			},
		},
	}
}
