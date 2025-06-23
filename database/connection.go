package database

import (
	"context"
	"fmt"
	"log"
	"mcp-go-server/config"

	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/sashabaranov/go-openai"
)

type Database struct {
	PineconeClient *pinecone.Client
	OpenAIClient   *openai.Client
	Config         *config.Config
}

var DB *Database

func ConnectDatabase(cfg *config.Config) (*Database, error) {
	log.Println("Connecting to external services...")

	// Initialize Pinecone client
	pineconeClient, err := pinecone.NewClient(pinecone.NewClientParams{
		ApiKey: cfg.PineconeAPIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Pinecone client: %w", err)
	}

	// Initialize OpenAI client
	openaiClient := openai.NewClient(cfg.OpenAIAPIKey)

	// Test connections
	if err := testConnections(pineconeClient, openaiClient, cfg); err != nil {
		return nil, fmt.Errorf("connection test failed: %w", err)
	}

	DB = &Database{
		PineconeClient: pineconeClient,
		OpenAIClient:   openaiClient,
		Config:         cfg,
	}

	log.Println("Successfully connected to all services")
	return DB, nil
}

func testConnections(pineconeClient *pinecone.Client, openaiClient *openai.Client, cfg *config.Config) error {
	// Test Pinecone connection
	_, err := pineconeClient.Index(pinecone.NewIndexConnParams{
		Host: cfg.PineconeHost,
	})
	if err != nil {
		return fmt.Errorf("Pinecone connection test failed: %w", err)
	}

	// Test OpenAI connection with a small request
	_, err = openaiClient.CreateEmbeddings(
		context.Background(),
		openai.EmbeddingRequest{
			Model: openai.AdaEmbeddingV2,
			Input: []string{"test"},
		},
	)
	if err != nil {
		return fmt.Errorf("OpenAI connection test failed: %w", err)
	}

	return nil
}
