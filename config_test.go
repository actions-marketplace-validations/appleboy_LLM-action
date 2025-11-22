package main

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectError bool
		validate    func(*testing.T, *Config)
	}{
		{
			name: "Valid config with all fields",
			envVars: map[string]string{
				"INPUT_BASE_URL":        "http://localhost:8080/v1",
				"INPUT_API_KEY":         "test-key",
				"INPUT_MODEL":           "gpt-4",
				"INPUT_SKIP_SSL_VERIFY": "true",
				"INPUT_SYSTEM_PROMPT":   "You are helpful",
				"INPUT_INPUT_PROMPT":    "Hello",
				"INPUT_TEMPERATURE":     "0.5",
				"INPUT_MAX_TOKENS":      "500",
			},
			expectError: false,
			validate: func(t *testing.T, c *Config) {
				if c.BaseURL != "http://localhost:8080/v1" {
					t.Errorf("expected base_url 'http://localhost:8080/v1', got '%s'", c.BaseURL)
				}
				if c.APIKey != "test-key" {
					t.Errorf("expected api_key 'test-key', got '%s'", c.APIKey)
				}
				if c.Model != "gpt-4" {
					t.Errorf("expected model 'gpt-4', got '%s'", c.Model)
				}
				if !c.SkipSSLVerify {
					t.Error("expected skip_ssl_verify to be true")
				}
				if c.SystemPrompt != "You are helpful" {
					t.Errorf("expected system_prompt 'You are helpful', got '%s'", c.SystemPrompt)
				}
				if c.InputPrompt != "Hello" {
					t.Errorf("expected input_prompt 'Hello', got '%s'", c.InputPrompt)
				}
				if c.Temperature != 0.5 {
					t.Errorf("expected temperature 0.5, got %f", c.Temperature)
				}
				if c.MaxTokens != 500 {
					t.Errorf("expected max_tokens 500, got %d", c.MaxTokens)
				}
			},
		},
		{
			name: "Default values",
			envVars: map[string]string{
				"INPUT_API_KEY":      "test-key",
				"INPUT_INPUT_PROMPT": "Hello",
			},
			expectError: false,
			validate: func(t *testing.T, c *Config) {
				if c.BaseURL != "https://api.openai.com/v1" {
					t.Errorf("expected default base_url, got '%s'", c.BaseURL)
				}
				if c.Temperature != 0.7 {
					t.Errorf("expected default temperature 0.7, got %f", c.Temperature)
				}
				if c.MaxTokens != 1000 {
					t.Errorf("expected default max_tokens 1000, got %d", c.MaxTokens)
				}
				if c.SkipSSLVerify {
					t.Error("expected default skip_ssl_verify to be false")
				}
			},
		},
		{
			name: "Missing API key",
			envVars: map[string]string{
				"INPUT_INPUT_PROMPT": "Hello",
			},
			expectError: true,
		},
		{
			name: "Missing input prompt",
			envVars: map[string]string{
				"INPUT_API_KEY": "test-key",
			},
			expectError: true,
		},
		{
			name: "Invalid temperature",
			envVars: map[string]string{
				"INPUT_API_KEY":      "test-key",
				"INPUT_INPUT_PROMPT": "Hello",
				"INPUT_TEMPERATURE":  "invalid",
			},
			expectError: true,
		},
		{
			name: "Invalid max tokens",
			envVars: map[string]string{
				"INPUT_API_KEY":      "test-key",
				"INPUT_INPUT_PROMPT": "Hello",
				"INPUT_MAX_TOKENS":   "abc",
			},
			expectError: true,
		},
		{
			name: "Negative max tokens",
			envVars: map[string]string{
				"INPUT_API_KEY":      "test-key",
				"INPUT_INPUT_PROMPT": "Hello",
				"INPUT_MAX_TOKENS":   "-100",
			},
			expectError: true,
		},
		{
			name: "Invalid skip SSL verify",
			envVars: map[string]string{
				"INPUT_API_KEY":         "test-key",
				"INPUT_INPUT_PROMPT":    "Hello",
				"INPUT_SKIP_SSL_VERIFY": "invalid",
			},
			expectError: true,
		},
		{
			name: "Multiline system prompt",
			envVars: map[string]string{
				"INPUT_API_KEY":       "test-key",
				"INPUT_INPUT_PROMPT":  "Hello",
				"INPUT_SYSTEM_PROMPT": "You are a helpful assistant.\nProvide clear responses.\n\nAlways be concise.",
			},
			expectError: false,
			validate: func(t *testing.T, c *Config) {
				expected := "You are a helpful assistant.\nProvide clear responses.\n\nAlways be concise."
				if c.SystemPrompt != expected {
					t.Errorf("expected multiline system_prompt, got '%s'", c.SystemPrompt)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all env vars
			clearEnvVars()

			// Set test env vars
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Load config
			config, err := LoadConfig()

			// Check error expectation
			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Validate config if no error expected
			if !tt.expectError && tt.validate != nil {
				tt.validate(t, config)
			}

			// Cleanup
			clearEnvVars()
		})
	}
}

func TestConfigParseTemperature(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    float64
		expectError bool
	}{
		{"Valid temperature", "0.5", 0.5, false},
		{"Max temperature", "2.0", 2.0, false},
		{"Min temperature", "0.0", 0.0, false},
		{"Empty string", "", 0.7, false}, // should keep default
		{"Invalid temperature", "invalid", 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{Temperature: 0.7}
			err := config.parseTemperature(tt.input)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectError && config.Temperature != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, config.Temperature)
			}
		})
	}
}

func TestConfigParseMaxTokens(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    int
		expectError bool
	}{
		{"Valid tokens", "500", 500, false},
		{"Large tokens", "4000", 4000, false},
		{"Empty string", "", 1000, false}, // should keep default
		{"Invalid tokens", "abc", 0, true},
		{"Negative tokens", "-100", 0, true},
		{"Zero tokens", "0", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{MaxTokens: 1000}
			err := config.parseMaxTokens(tt.input)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectError && config.MaxTokens != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, config.MaxTokens)
			}
		})
	}
}

func TestConfigParseSkipSSL(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    bool
		expectError bool
	}{
		{"True lowercase", "true", true, false},
		{"True uppercase", "TRUE", true, false},
		{"False lowercase", "false", false, false},
		{"Numeric true", "1", true, false},
		{"Numeric false", "0", false, false},
		{"Empty string", "", false, false}, // should keep default
		{"Invalid value", "invalid", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{SkipSSLVerify: false}
			err := config.parseSkipSSL(tt.input)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectError && config.SkipSSLVerify != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, config.SkipSSLVerify)
			}
		})
	}
}

// Helper function to clear all env vars
func clearEnvVars() {
	os.Unsetenv("INPUT_BASE_URL")
	os.Unsetenv("INPUT_API_KEY")
	os.Unsetenv("INPUT_MODEL")
	os.Unsetenv("INPUT_SKIP_SSL_VERIFY")
	os.Unsetenv("INPUT_SYSTEM_PROMPT")
	os.Unsetenv("INPUT_INPUT_PROMPT")
	os.Unsetenv("INPUT_TEMPERATURE")
	os.Unsetenv("INPUT_MAX_TOKENS")
}
