package main

import (
	"crypto/tls"
	"net/http"

	openai "github.com/sashabaranov/go-openai"
)

// NewClient creates a new OpenAI client with the given configuration
func NewClient(config *Config) *openai.Client {
	clientConfig := openai.DefaultConfig(config.APIKey)
	clientConfig.BaseURL = config.BaseURL

	// Handle SSL verification
	if config.SkipSSLVerify {
		clientConfig.HTTPClient = createInsecureHTTPClient()
	}

	return openai.NewClientWithConfig(clientConfig)
}

// createInsecureHTTPClient creates an HTTP client that skips SSL verification
// #nosec G402 - This is intentionally configurable by the user for local/self-hosted LLM services
func createInsecureHTTPClient() *http.Client {
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // #nosec G402
	}
	return &http.Client{
		Transport: customTransport,
	}
}
