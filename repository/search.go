package repository

import (
	"context"
	"fmt"
	"mcp-go-server/database"
	"mcp-go-server/domain"
	"mcp-go-server/models"
	"strings"

	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/protobuf/types/known/structpb"
)

// CheckRepositoryExists checks if repository exists in vector database
func CheckRepositoryExists(repository, branch string) (bool, error) {
	if database.DB == nil {
		return false, fmt.Errorf("database not initialized")
	}

	index, err := database.DB.PineconeClient.Index(pinecone.NewIndexConnParams{
		Host: database.DB.Config.PineconeHost,
	})
	if err != nil {
		return false, fmt.Errorf("failed to connect to index: %w", err)
	}

	// Create a test query to check if repository exists
	testEmbedding := make([]float32, 1536) // OpenAI embedding dimension
	for i := range testEmbedding {
		testEmbedding[i] = 0.1
	}

	filterStruct, err := structpb.NewStruct(map[string]interface{}{
		"repository": repository,
		"branch":     branch,
	})
	if err != nil {
		return false, fmt.Errorf("failed to create filter: %w", err)
	}

	queryResp, err := index.QueryByVectorValues(context.Background(), &pinecone.QueryByVectorValuesRequest{
		Vector:          testEmbedding,
		TopK:            1,
		MetadataFilter:  filterStruct,
		IncludeMetadata: true,
	})
	if err != nil {
		return false, fmt.Errorf("query failed: %w", err)
	}

	return len(queryResp.Matches) > 0, nil
}

// GetQueryEmbedding generates embedding for search query
func GetQueryEmbedding(query string) ([]float32, error) {
	if database.DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	resp, err := database.DB.OpenAIClient.CreateEmbeddings(
		context.Background(),
		openai.EmbeddingRequest{
			Model: openai.AdaEmbeddingV2,
			Input: []string{query},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	// Convert []float64 to []float32
	embedding := make([]float32, len(resp.Data[0].Embedding))
	for i, v := range resp.Data[0].Embedding {
		embedding[i] = float32(v)
	}

	return embedding, nil
}

// SearchVectors performs vector search in Pinecone
func SearchVectors(queryEmbedding []float32, repository, branch string, limit int) ([]domain.SearchResult, error) {
	if database.DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	index, err := database.DB.PineconeClient.Index(pinecone.NewIndexConnParams{
		Host: database.DB.Config.PineconeHost,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to index: %w", err)
	}

	// Create filter for repository and branch
	filterStruct, err := structpb.NewStruct(map[string]interface{}{
		"repository": repository,
		"branch":     branch,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create filter: %w", err)
	}

	// Perform query
	queryResp, err := index.QueryByVectorValues(context.Background(), &pinecone.QueryByVectorValuesRequest{
		Vector:          queryEmbedding,
		TopK:            uint32(limit),
		MetadataFilter:  filterStruct,
		IncludeMetadata: true,
	})
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// Convert results
	var results []domain.SearchResult
	for _, match := range queryResp.Matches {
		if match == nil || match.Vector == nil || match.Vector.Metadata == nil {
			continue
		}

		metadata := match.Vector.Metadata.AsMap()

		result := domain.SearchResult{
			CodeChunk: domain.CodeChunk{
				ID:         match.Vector.Id,
				Content:    metadata["content"].(string),
				FilePath:   metadata["filePath"].(string),
				Repository: metadata["repository"].(string),
				Branch:     metadata["branch"].(string),
				Language:   metadata["language"].(string),
			},
			Score: match.Score,
		}

		results = append(results, result)
	}

	return results, nil
}

// GenerateAISummary generates AI summary for search results
func GenerateAISummary(results []models.SearchResult, query string) (string, error) {
	if database.DB == nil {
		return "", fmt.Errorf("database not initialized")
	}

	if len(results) == 0 {
		return "No results found for the query.", nil
	}

	// Build context from search results
	var contextBuilder strings.Builder
	contextBuilder.WriteString("Based on the following code search results:\n\n")

	for i, result := range results {
		contextBuilder.WriteString(fmt.Sprintf("Result %d - File: %s\n", i+1, result.FilePath))
		contextBuilder.WriteString(fmt.Sprintf("Language: %s\n", result.Language))
		contextBuilder.WriteString(fmt.Sprintf("Content:\n%s\n\n", result.Content))

		// Limit context size
		if contextBuilder.Len() > 3000 {
			contextBuilder.WriteString("... (truncated for brevity)\n")
			break
		}
	}

	// Generate summary using OpenAI
	completion, err := database.DB.OpenAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: "system",
					Content: `You are a technical expert analyzing code search results. 
					Provide a concise, helpful summary that directly answers the user's query.
					Focus on the most relevant information from the code snippets.
					Be specific and technical when appropriate.`,
				},
				{
					Role: "user",
					Content: fmt.Sprintf(`Query: %s

%s

Provide a concise summary that answers the query based on the code search results.`,
						query, contextBuilder.String()),
				},
			},
			Temperature: 0.3,
			MaxTokens:   300,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate summary: %w", err)
	}

	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	return completion.Choices[0].Message.Content, nil
}
