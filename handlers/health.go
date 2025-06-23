package handlers

import (
	"mcp-go-server/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health returns the health status of the server
func Health(c *gin.Context) {
	healthData := map[string]interface{}{
		"status":    "healthy",
		"service":   "MCP Vector Search Server",
		"version":   "1.0.0",
		"timestamp": "2025-06-23T00:00:00Z",
	}

	successRes := response.ClientResponse(http.StatusOK, "Server is healthy", healthData, nil)
	c.JSON(http.StatusOK, successRes)
}