package handlers

import (
	"mcp-go-server/models"
	"mcp-go-server/response"
	"mcp-go-server/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// VectorSearch performs vector search on repository code
func VectorSearch(c *gin.Context) {
	var searchReq models.SearchRequest

	if err := c.ShouldBindJSON(&searchReq); err != nil {
		errRes := response.ErrorClientResponse(http.StatusBadRequest, "Invalid request format", err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	// Validate the request
	if err := validator.New().Struct(searchReq); err != nil {
		errRes := response.ErrorClientResponse(http.StatusBadRequest, "Validation failed", err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	// Set default values
	if searchReq.Branch == "" {
		searchReq.Branch = "main"
	}
	if searchReq.Limit <= 0 {
		searchReq.Limit = 10
	}

	// Perform search
	results, err := usecase.PerformVectorSearch(searchReq)
	if err != nil {
		errRes := response.ErrorClientResponse(http.StatusInternalServerError, "Search failed", err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Search completed successfully", results, nil)
	c.JSON(http.StatusOK, successRes)
}

// VectorSearchWithSummary performs vector search and generates AI summary
func VectorSearchWithSummary(c *gin.Context) {
	var searchReq models.SearchRequest

	if err := c.ShouldBindJSON(&searchReq); err != nil {
		errRes := response.ErrorClientResponse(http.StatusBadRequest, "Invalid request format", err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	// Validate the request
	if err := validator.New().Struct(searchReq); err != nil {
		errRes := response.ErrorClientResponse(http.StatusBadRequest, "Validation failed", err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	// Set default values
	if searchReq.Branch == "" {
		searchReq.Branch = "main"
	}
	if searchReq.Limit <= 0 {
		searchReq.Limit = 5
	}

	// Perform search with summary
	summary, err := usecase.PerformSearchWithSummary(searchReq)
	if err != nil {
		errRes := response.ErrorClientResponse(http.StatusInternalServerError, "Search with summary failed", err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Search with summary completed", summary, nil)
	c.JSON(http.StatusOK, successRes)
}
