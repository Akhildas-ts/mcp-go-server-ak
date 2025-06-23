package router

import (
	"mcp-go-server/database"
	"mcp-go-server/handlers"
	"mcp-go-server/middleware"

	"github.com/gin-gonic/gin"
)

func MCPRoutes(rg *gin.RouterGroup, db *database.Database) {
	// Public routes
	rg.GET("/health", handlers.Health)
	rg.GET("/auth/login", handlers.Login)
	rg.GET("/auth/callback", handlers.Callback)

	// Protected routes - require authentication
	protected := rg.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Vector search endpoints
		protected.POST("/search", handlers.VectorSearch)
		protected.POST("/search/summary", handlers.VectorSearchWithSummary)

		// Repository indexing endpoints
		protected.POST("/index", handlers.IndexRepository)
		protected.GET("/repositories", handlers.GetRepositories)

		// User management endpoints
		protected.GET("/profile", handlers.GetProfile)
	}
}
