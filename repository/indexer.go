package repository

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"mcp-go-server/database"
	"mcp-go-server/domain"
	"mcp-go-server/helper"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/protobuf/types/known/structpb"
)

// CloneRepository clones a Git repository to temporary directory
func CloneRepository(repoURL, branch string) (string, error) {
	log.Printf("📥 Creating temporary directory for repository...")
	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "repo-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	log.Printf("📁 Temporary directory created: %s", tempDir)

	// Clone repository
	log.Printf("🔗 Cloning repository from: %s", repoURL)
	cmd := exec.Command("git", "clone", repoURL, tempDir)
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to clone repository: %w", err)
	}
	log.Printf("✅ Repository cloned successfully")

	// Checkout specific branch if specified
	if branch != "" && branch != "main" && branch != "master" {
		log.Printf("🌿 Checking out branch: %s", branch)
		cmd = exec.Command("git", "checkout", branch)
		cmd.Dir = tempDir
		if err := cmd.Run(); err != nil {
			os.RemoveAll(tempDir)
			return "", fmt.Errorf("failed to checkout branch %s: %w", branch, err)
		}
		log.Printf("✅ Branch checkout completed")
	}

	return tempDir, nil
}

// ProcessRepositoryFiles processes all files in repository
func ProcessRepositoryFiles(repoPath, repoURL, branch string) (int, int, error) {
	log.Printf("🔍 Scanning repository for files to process...")
	fileCount := 0
	chunkCount := 0
	processedFiles := 0
	skippedFiles := 0

	// First pass: count total files
	totalFiles := 0
	filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		if !helper.IsBinaryFile(path) {
			totalFiles++
		}
		return nil
	})

	log.Printf("📊 Found %d files to process", totalFiles)

	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and hidden files
		if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() && strings.Contains(path, ".git") {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip binary files
		if helper.IsBinaryFile(path) {
			skippedFiles++
			return nil
		}

		// Read file content
		content, err := ioutil.ReadFile(path)
		if err != nil {
			skippedFiles++
			return nil // Skip files we can't read
		}

		// Skip large files and binary content
		if len(content) > 100000 || helper.ContainsBinaryData(content) {
			skippedFiles++
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(repoPath, path)
		if err != nil {
			relPath = path
		}

		processedFiles++
		log.Printf("📄 Processing file %d/%d: %s", processedFiles, totalFiles, relPath)

		// Process file
		chunks, err := processFile(string(content), relPath, repoURL, branch)
		if err != nil {
			log.Printf("⚠️  Failed to process file %s: %v", relPath, err)
			return nil // Skip files that fail processing
		}

		fileCount++
		chunkCount += chunks
		log.Printf("✅ Processed %s (%d chunks)", relPath, chunks)
		return nil
	})

	log.Printf("📈 Processing completed:")
	log.Printf("   - Files processed: %d", fileCount)
	log.Printf("   - Files skipped: %d", skippedFiles)
	log.Printf("   - Total chunks created: %d", chunkCount)

	return fileCount, chunkCount, err
}

// processFile processes a single file and stores chunks
func processFile(content, filePath, repoURL, branch string) (int, error) {
	if database.DB == nil {
		return 0, fmt.Errorf("database not initialized")
	}

	// Extract repository name from URL
	repoName := helper.ExtractRepoName(repoURL)

	// Determine language
	language := helper.GetLanguageFromExtension(filepath.Ext(filePath))

	// Split content into chunks
	chunks := helper.SplitIntoChunks(content, 1000)
	log.Printf("   📝 Split into %d chunks", len(chunks))

	index, err := database.DB.PineconeClient.Index(pinecone.NewIndexConnParams{
		Host: database.DB.Config.PineconeHost,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to connect to index: %w", err)
	}

	successfulChunks := 0
	// Process each chunk
	for i, chunk := range chunks {
		// Get embedding
		embedding, err := getEmbedding(chunk)
		if err != nil {
			log.Printf("   ⚠️  Failed to generate embedding for chunk %d: %v", i+1, err)
			continue // Skip chunks that fail embedding
		}

		// Create metadata
		metadata, err := structpb.NewStruct(map[string]interface{}{
			"content":    chunk,
			"filePath":   filePath,
			"repository": repoName,
			"branch":     branch,
			"language":   language,
		})
		if err != nil {
			log.Printf("   ⚠️  Failed to create metadata for chunk %d: %v", i+1, err)
			continue
		}

		// Create vector ID
		vectorID := fmt.Sprintf("%s-%s-%d", repoName, strings.ReplaceAll(filePath, "/", "-"), i)
		if len(vectorID) > 100 {
			vectorID = vectorID[:100]
		}

		// Store in Pinecone
		vectors := []*pinecone.Vector{
			{
				Id:       vectorID,
				Values:   embedding,
				Metadata: metadata,
			},
		}

		_, err = index.UpsertVectors(context.Background(), vectors)
		if err != nil {
			log.Printf("   ⚠️  Failed to store chunk %d in Pinecone: %v", i+1, err)
			continue // Skip chunks that fail to store
		}

		successfulChunks++
	}

	log.Printf("   💾 Successfully stored %d/%d chunks in Pinecone", successfulChunks, len(chunks))
	return successfulChunks, nil
}

// getEmbedding generates embedding for text
func getEmbedding(text string) ([]float32, error) {
	if database.DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	resp, err := database.DB.OpenAIClient.CreateEmbeddings(
		context.Background(),
		openai.EmbeddingRequest{
			Model: openai.AdaEmbeddingV2,
			Input: []string{text},
		},
	)
	if err != nil {
		return nil, err
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

// GetRepositoryInfo retrieves information about indexed repositories
func GetRepositoryInfo(userID string) ([]domain.Repository, error) {
	// In a real implementation, this would query a database
	// For now, return empty slice as this is a demo
	return []domain.Repository{}, nil
}

// SaveRepositoryInfo saves repository indexing information
func SaveRepositoryInfo(repo domain.Repository) error {
	// In a real implementation, this would save to a database
	// For now, just return nil as this is a demo
	return nil
}
