package main

import (
	"testing"

	openai "github.com/sashabaranov/go-openai"
)

func TestBuildMessages(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectedLen int
		checkSystem bool
	}{
		{
			name: "Both system and user prompts",
			config: &Config{
				SystemPrompt: "You are a helpful assistant",
				InputPrompt:  "Hello, how are you?",
			},
			expectedLen: 2,
			checkSystem: true,
		},
		{
			name: "Only user prompt",
			config: &Config{
				SystemPrompt: "",
				InputPrompt:  "Hello, how are you?",
			},
			expectedLen: 1,
			checkSystem: false,
		},
		{
			name: "Long system prompt",
			config: &Config{
				SystemPrompt: "You are a code reviewer. Provide constructive feedback on code quality, best practices, and potential issues.",
				InputPrompt:  "Review this code",
			},
			expectedLen: 2,
			checkSystem: true,
		},
		{
			name: "Multiline input prompt",
			config: &Config{
				SystemPrompt: "",
				InputPrompt:  "Line 1\nLine 2\nLine 3",
			},
			expectedLen: 1,
			checkSystem: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messages := BuildMessages(tt.config)

			if len(messages) != tt.expectedLen {
				t.Errorf("expected %d messages, got %d", tt.expectedLen, len(messages))
			}

			// Check if system message exists when expected
			if tt.checkSystem {
				if messages[0].Role != openai.ChatMessageRoleSystem {
					t.Error("expected first message to be system role")
				}
				if messages[0].Content != tt.config.SystemPrompt {
					t.Errorf(
						"expected system message content '%s', got '%s'",
						tt.config.SystemPrompt,
						messages[0].Content,
					)
				}
			}

			// Check user message (always last)
			lastMsg := messages[len(messages)-1]
			if lastMsg.Role != openai.ChatMessageRoleUser {
				t.Error("expected last message to be user role")
			}
			if lastMsg.Content != tt.config.InputPrompt {
				t.Errorf(
					"expected user message content '%s', got '%s'",
					tt.config.InputPrompt,
					lastMsg.Content,
				)
			}
		})
	}
}

func TestBuildMessagesEmpty(t *testing.T) {
	config := &Config{
		SystemPrompt: "",
		InputPrompt:  "",
	}

	messages := BuildMessages(config)

	// Should still create a user message even if input is empty
	if len(messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(messages))
	}

	if messages[0].Role != openai.ChatMessageRoleUser {
		t.Error("expected message to be user role")
	}
}

func TestBuildMessagesOrder(t *testing.T) {
	config := &Config{
		SystemPrompt: "System message",
		InputPrompt:  "User message",
	}

	messages := BuildMessages(config)

	if len(messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(messages))
	}

	// Verify order: system first, then user
	if messages[0].Role != openai.ChatMessageRoleSystem {
		t.Error("expected first message to be system role")
	}
	if messages[1].Role != openai.ChatMessageRoleUser {
		t.Error("expected second message to be user role")
	}
}
