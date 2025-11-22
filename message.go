package main

import openai "github.com/sashabaranov/go-openai"

// BuildMessages builds the chat completion messages from configuration
func BuildMessages(config *Config) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{}

	// Add system prompt if provided
	if config.SystemPrompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: config.SystemPrompt,
		})
	}

	// Add user prompt
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: config.InputPrompt,
	})

	return messages
}
