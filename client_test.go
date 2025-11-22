package main

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
	}{
		{
			name: "Default OpenAI config",
			config: &Config{
				BaseURL:       "https://api.openai.com/v1",
				APIKey:        "test-key",
				SkipSSLVerify: false,
			},
		},
		{
			name: "Custom base URL",
			config: &Config{
				BaseURL:       "http://localhost:8080/v1",
				APIKey:        "test-key",
				SkipSSLVerify: false,
			},
		},
		{
			name: "Skip SSL verification",
			config: &Config{
				BaseURL:       "https://api.openai.com/v1",
				APIKey:        "test-key",
				SkipSSLVerify: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.config)

			if client == nil {
				t.Error("expected client to be created, got nil")
			}
		})
	}
}

func TestCreateInsecureHTTPClient(t *testing.T) {
	client := createInsecureHTTPClient()

	if client == nil {
		t.Fatal("expected HTTP client to be created, got nil")
	}

	if client.Transport == nil {
		t.Error("expected transport to be set, got nil")
	}
}
