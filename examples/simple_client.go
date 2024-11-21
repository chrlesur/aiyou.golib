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

    ctx := context.Background()

    // Exemple d'utilisation de ChatCompletion
    req := aiyou.ChatCompletionRequest{
        Messages: []aiyou.Message{
            {
                Role: "user",
                Content: []aiyou.ContentPart{
                    {Type: "text", Text: "What is the capital of France?"},
                },
            },
        },
        AssistantID: "your-assistant-id",
    }

    resp, err := client.ChatCompletion(ctx, req)
    if err != nil {
        log.Fatalf("Error in ChatCompletion: %v", err)
    }

    fmt.Printf("AI response: %s\n", resp.Choices[0].Message.Content[0].Text)

    // Exemple d'utilisation de ChatCompletionStream
    streamReq := aiyou.ChatCompletionRequest{
        Messages: []aiyou.Message{
            {
                Role: "user",
                Content: []aiyou.ContentPart{
                    {Type: "text", Text: "Tell me a short story."},
                },
            },
        },
        AssistantID: "your-assistant-id",
        Stream:      true,
    }

    stream, err := client.ChatCompletionStream(ctx, streamReq)
    if err != nil {
        log.Fatalf("Error in ChatCompletionStream: %v", err)
    }

    fmt.Println("Streaming response:")
    for {
        chunk, err := stream.ReadChunk()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatalf("Error reading chunk: %v", err)
        }
        fmt.Print(chunk.Choices[0].Message.Content[0].Text)
    }
}