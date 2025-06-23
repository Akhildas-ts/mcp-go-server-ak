package usecase

import (
	"errors"
	"mcp-go-server/models"
	"mcp-go-server/repository"
)

// PerformVectorSearch executes vector search on repository code
func PerformVectorSearch(searchReq models.SearchRequest) (models.SearchResponse, error) {
	// Validate repository exists
	exists, err := repository.CheckRepositoryExists(searchReq.Repository, searchReq.Branch)
	if err != nil {
		return models.SearchResponse{}, err
	}
	if !exists {
		return models.SearchResponse{}, models.ErrRepositoryNotFound
	}

	// Get query embedding
	embedding, err := repository.GetQueryEmbedding(searchReq.Query)
	if err != nil {
		return models.SearchResponse{}, errors.New("failed to generate query embedding")
	}

	// Perform vector search
	results, err := repository.SearchVectors(embedding, searchReq.Repository, searchReq.Branch, searchReq.Limit)
	if err != nil {
		return models.SearchResponse{}, errors.New("vector search failed")
	}

	// Convert to response format
	var searchResults []models.SearchResult
	for _, result := range results {
		searchResults = append(searchResults, models.SearchResult{
			Content:    result.Content,
			FilePath:   result.FilePath,
			Repository: result.Repository,
			Branch:     result.Branch,
			Language:   result.Language,
			Score:      result.Score,
		})
	}

	return models.SearchResponse{
		Results: searchResults,
		Total:   len(searchResults),
	}, nil
}

// PerformSearchWithSummary executes search and generates AI summary
func PerformSearchWithSummary(searchReq models.SearchRequest) (models.SearchWithSummaryResponse, error) {
	// First perform regular search
	searchResponse, err := PerformVectorSearch(searchReq)
	if err != nil {
		return models.SearchWithSummaryResponse{}, err
	}

	// Generate AI summary
	summary, err := repository.GenerateAISummary(searchResponse.Results, searchReq.Query)
	if err != nil {
		return models.SearchWithSummaryResponse{}, errors.New("failed to generate summary")
	}

	return models.SearchWithSummaryResponse{
		Summary: summary,
		Results: searchResponse.Results,
		Total:   searchResponse.Total,
	}, nil
}
