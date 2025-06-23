package config

import (
	"errors"
	"os"
)

type Config struct {
	Port                   string
	PineconeAPIKey         string
	PineconeIndexName      string
	PineconeHost           string
	OpenAIAPIKey           string
	GitHubClientID         string
	GitHubClientSecret     string
	GitHubOAuthRedirectURL string
	JWTSecret              string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Port:                   getEnv("PORT", "8081"),
		PineconeAPIKey:         getEnv("PINECONE_API_KEY", ""),
		PineconeIndexName:      getEnv("PINECONE_INDEX_NAME", "default-index"),
		PineconeHost:           getEnv("PINECONE_HOST", ""),
		OpenAIAPIKey:           getEnv("OPENAI_API_KEY", ""),
		GitHubClientID:         getEnv("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret:     getEnv("GITHUB_CLIENT_SECRET", ""),
		GitHubOAuthRedirectURL: getEnv("GITHUB_OAUTH_REDIRECT_URL", "http://localhost:8081/auth/github/callback"),
		JWTSecret:              getEnv("JWT_SECRET", "mcp-secret-key"),
	}

	// Validate required fields with helpful error messages
	if cfg.PineconeAPIKey == "" {
		return nil, errors.New("PINECONE_API_KEY is required. Please set it in your environment variables or .env file")
	}
	if cfg.OpenAIAPIKey == "" {
		return nil, errors.New("OPENAI_API_KEY is required. Please set it in your environment variables or .env file")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
