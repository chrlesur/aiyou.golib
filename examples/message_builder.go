/*
Copyright (C) 2024 Cloud Temple

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
*/

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/chrlesur/aiyou.golib"
)

func main() {
	client, err := aiyou.NewClient("your-email@example.com", "your-password")
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Create a complex message using MessageBuilder
	builder := aiyou.NewMessageBuilder("user", client.Logger())
	message := builder.
		AddText("I have a question about this image:").
		AddImage("https://example.com/image.jpg").
		AddText("What can you tell me about it?").
		Build()

	// Create a chat completion request with the complex message
	req := aiyou.ChatCompletionRequest{
		Messages:    []aiyou.Message{message},
		AssistantID: "your-assistant-id",
	}

	// Send the request
	resp, err := client.ChatCompletion(context.Background(), req)
	if err != nil {
		log.Fatalf("Error in ChatCompletion: %v", err)
	}

	fmt.Printf("AI response: %s\n", resp.Choices[0].Message.Content[0].Text)
}
