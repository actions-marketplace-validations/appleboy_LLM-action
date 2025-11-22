package main

import (
	"context"
	"fmt"
	"os"

	"github.com/appleboy/com/gh"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	// Create OpenAI client
	client := NewClient(config)

	// Build messages
	messages := BuildMessages(config)

	// Create chat completion request
	req := openai.ChatCompletionRequest{
		Model:       config.Model,
		Messages:    messages,
		Temperature: float32(config.Temperature),
		MaxTokens:   config.MaxTokens,
	}

	fmt.Println("Sending request to LLM...")
	fmt.Printf("Model: %s\n", config.Model)
	fmt.Printf("Base URL: %s\n", config.BaseURL)

	// Call the API
	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return fmt.Errorf("chat completion error: %v", err)
	}

	// Extract response
	if len(resp.Choices) == 0 {
		return fmt.Errorf("no response from LLM")
	}

	response := resp.Choices[0].Message.Content

	// Print response for debugging
	fmt.Println("--- LLM Response ---")
	fmt.Println(response)
	fmt.Println("--- End Response ---")

	// Set GitHub Actions output
	if err := gh.SetOutput(map[string]string{
		"response": response,
	}); err != nil {
		return fmt.Errorf("failed to set output: %v", err)
	}

	return nil
}
