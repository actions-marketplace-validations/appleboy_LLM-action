package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"

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
	// Get input parameters from environment variables
	baseURL := os.Getenv("INPUT_BASE_URL")
	apiKey := os.Getenv("INPUT_API_KEY")
	model := os.Getenv("INPUT_MODEL")
	skipSSLVerify := os.Getenv("INPUT_SKIP_SSL_VERIFY")
	systemPrompt := os.Getenv("INPUT_SYSTEM_PROMPT")
	inputPrompt := os.Getenv("INPUT_INPUT_PROMPT")
	temperatureStr := os.Getenv("INPUT_TEMPERATURE")
	maxTokensStr := os.Getenv("INPUT_MAX_TOKENS")

	// Set default base URL if not provided
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	// Validate required inputs
	if apiKey == "" {
		return fmt.Errorf("api_key is required")
	}
	if inputPrompt == "" {
		return fmt.Errorf("input_prompt is required")
	}

	// Parse optional parameters
	temperature := 0.7
	if temperatureStr != "" {
		temp, err := strconv.ParseFloat(temperatureStr, 32)
		if err != nil {
			return fmt.Errorf("invalid temperature value: %v", err)
		}
		temperature = temp
	}

	maxTokens := 1000
	if maxTokensStr != "" {
		tokens, err := strconv.Atoi(maxTokensStr)
		if err != nil {
			return fmt.Errorf("invalid max_tokens value: %v", err)
		}
		maxTokens = tokens
	}

	skipSSL := false
	if skipSSLVerify != "" {
		skip, err := strconv.ParseBool(skipSSLVerify)
		if err != nil {
			return fmt.Errorf("invalid skip_ssl_verify value: %v", err)
		}
		skipSSL = skip
	}

	// Configure OpenAI client
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL

	// Handle SSL verification
	if skipSSL {
		customTransport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		config.HTTPClient = &http.Client{
			Transport: customTransport,
		}
	}

	client := openai.NewClientWithConfig(config)

	// Prepare messages
	messages := []openai.ChatCompletionMessage{}

	// Add system prompt if provided
	if systemPrompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		})
	}

	// Add user prompt
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: inputPrompt,
	})

	// Create chat completion request
	req := openai.ChatCompletionRequest{
		Model:       model,
		Messages:    messages,
		Temperature: float32(temperature),
		MaxTokens:   maxTokens,
	}

	fmt.Println("Sending request to LLM...")
	fmt.Printf("Model: %s\n", model)
	fmt.Printf("Base URL: %s\n", baseURL)

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

	if err := gh.SetOutput(map[string]string{
		"response": response,
	}); err != nil {
		return fmt.Errorf("failed to set output: %v", err)
	}

	return nil
}
