package usecase

import (
	"errors"
	"fmt"
	"log"
	"mcp-go-server/helper"
	"mcp-go-server/models"
	"mcp-go-server/repository"
	"os"
	"time"
)

// IndexRepository indexes a GitHub repository
func IndexRepository(indexReq models.IndexRequest) (models.IndexResponse, error) {
	startTime := time.Now()
	log.Printf("🚀 Starting repository indexing for: %s (branch: %s)", indexReq.RepoURL, indexReq.Branch)

	// Validate repository URL
	log.Printf("📋 Validating repository URL...")
	if err := helper.ValidateGitHubRepoURL(indexReq.RepoURL); err != nil {
		log.Printf("❌ Repository URL validation failed: %v", err)
		return models.IndexResponse{}, fmt.Errorf("invalid repository URL: %w", err)
	}
	log.Printf("✅ Repository URL validation passed")

	// Validate branch name
	log.Printf("📋 Validating branch name...")
	if err := helper.ValidateBranch(indexReq.Branch); err != nil {
		log.Printf("❌ Branch validation failed: %v", err)
		return models.IndexResponse{}, fmt.Errorf("invalid branch name: %w", err)
	}
	log.Printf("✅ Branch validation passed")

	// Set default branch
	if indexReq.Branch == "" {
		indexReq.Branch = "main"
		log.Printf("📝 Using default branch: main")
	}

	// Clone repository
	log.Printf("📥 Cloning repository...")
	repoPath, err := repository.CloneRepository(indexReq.RepoURL, indexReq.Branch)
	if err != nil {
		log.Printf("❌ Repository cloning failed: %v", err)
		return models.IndexResponse{}, fmt.Errorf("failed to clone repository: %w", err)
	}
	defer os.RemoveAll(repoPath) // Clean up temp directory
	log.Printf("✅ Repository cloned successfully to: %s", repoPath)

	// Process repository files
	log.Printf("🔄 Processing repository files and generating embeddings...")
	fileCount, chunkCount, err := repository.ProcessRepositoryFiles(repoPath, indexReq.RepoURL, indexReq.Branch)
	if err != nil {
		log.Printf("❌ Repository processing failed: %v", err)
		return models.IndexResponse{}, fmt.Errorf("failed to process repository files: %w", err)
	}

	// Extract repository name
	repoName := helper.ExtractRepoName(indexReq.RepoURL)

	duration := time.Since(startTime)
	log.Printf("🎉 Repository indexing completed successfully!")
	log.Printf("📊 Summary:")
	log.Printf("   - Repository: %s", repoName)
	log.Printf("   - Branch: %s", indexReq.Branch)
	log.Printf("   - Files processed: %d", fileCount)
	log.Printf("   - Chunks created: %d", chunkCount)
	log.Printf("   - Total time: %v", duration)

	// Save repository info (in a real implementation, this would save to database)
	// For now, we'll skip this step

	return models.IndexResponse{
		Repository: repoName,
		Branch:     indexReq.Branch,
		FileCount:  fileCount,
		ChunkCount: chunkCount,
		Status:     "completed",
	}, nil
}

// GetRepositories retrieves list of indexed repositories for user
func GetRepositories(userID string) ([]models.RepositoryInfo, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Get repositories from storage
	repos, err := repository.GetRepositoryInfo(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve repositories: %w", err)
	}

	// Convert domain objects to response models
	var repoInfos []models.RepositoryInfo
	for _, repo := range repos {
		repoInfo := models.RepositoryInfo{
			Name:       repo.Name,
			Owner:      repo.Owner,
			URL:        repo.URL,
			Branch:     repo.Branch,
			IndexedAt:  repo.IndexedAt,
			FileCount:  repo.FileCount,
			ChunkCount: repo.ChunkCount,
		}
		repoInfos = append(repoInfos, repoInfo)
	}

	return repoInfos, nil
}

// DeleteRepository removes an indexed repository
func DeleteRepository(userID, repository, branch string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}

	if repository == "" {
		return errors.New("repository is required")
	}

	// In a real implementation, this would:
	// 1. Verify user owns the repository index
	// 2. Delete vectors from Pinecone
	// 3. Remove repository info from database

	return errors.New("delete repository functionality not implemented yet")
}

// GetIndexingStatus retrieves the status of a repository indexing operation
func GetIndexingStatus(userID, repository, branch string) (string, error) {
	if userID == "" {
		return "", errors.New("user ID is required")
	}

	if repository == "" {
		return "", errors.New("repository is required")
	}

	// In a real implementation, this would check the indexing status
	// For now, return a placeholder status
	return "completed", nil
}

// ReindexRepository re-indexes an existing repository
func ReindexRepository(userID string, indexReq models.IndexRequest) (models.IndexResponse, error) {
	if userID == "" {
		return models.IndexResponse{}, errors.New("user ID is required")
	}

	// First, delete existing index (in a real implementation)
	// Then perform new indexing
	return IndexRepository(indexReq)
}
