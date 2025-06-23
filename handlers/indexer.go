package handlers

import (
	"log"
	"mcp-go-server/models"
	"mcp-go-server/response"
	"mcp-go-server/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// IndexRepository indexes a GitHub repository
func IndexRepository(c *gin.Context) {
	var indexReq models.IndexRequest

	if err := c.ShouldBindJSON(&indexReq); err != nil {
		errRes := response.ErrorClientResponse(http.StatusBadRequest, "Invalid request format", err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	// Validate the request
	if err := validator.New().Struct(indexReq); err != nil {
		errRes := response.ErrorClientResponse(http.StatusBadRequest, "Validation failed", err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	// Set default branch
	if indexReq.Branch == "" {
		indexReq.Branch = "main"
	}

	// Log the start of indexing process
	log.Printf("🎯 Indexing request received for repository: %s (branch: %s)", indexReq.RepoURL, indexReq.Branch)
	log.Printf("⏱️  This process may take 5-10 minutes depending on repository size...")

	// Index repository
	result, err := usecase.IndexRepository(indexReq)
	if err != nil {
		log.Printf("❌ Indexing failed: %v", err)
		errRes := response.ErrorClientResponse(http.StatusInternalServerError, "Repository indexing failed", err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	log.Printf("🎉 Indexing completed successfully for: %s", result.Repository)
	successRes := response.ClientResponse(http.StatusOK, "Repository indexed successfully", result, nil)
	c.JSON(http.StatusOK, successRes)
}

// GetRepositories retrieves list of indexed repositories
func GetRepositories(c *gin.Context) {
	userID, exists := c.Get(models.UserIDKey)
	if !exists {
		errRes := response.ErrorClientResponse(http.StatusUnauthorized, "User not authenticated", nil)
		c.JSON(http.StatusUnauthorized, errRes)
		return
	}

	repositories, err := usecase.GetRepositories(userID.(string))
	if err != nil {
		errRes := response.ErrorClientResponse(http.StatusInternalServerError, "Failed to retrieve repositories", err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Repositories retrieved successfully", repositories, nil)
	c.JSON(http.StatusOK, successRes)
}
