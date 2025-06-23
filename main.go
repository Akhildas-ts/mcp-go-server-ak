package main

import (
	"log"
	"mcp-go-server/config"
	"mcp-go-server/database"
	"mcp-go-server/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title MCP Server API
// @version 1.0
// @description API for MCP Vector Search Server
// @host localhost:8081
// @BasePath /
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize database connections
	db, err := database.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Setup routes
	router.MCPRoutes(r.Group("/"), db)

	// Start server
	log.Printf("MCP Server starting on port %s...", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
