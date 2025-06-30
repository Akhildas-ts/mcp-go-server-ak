package main

import (
	"log"
	"mcp-go-server/config"
	"mcp-go-server/database"
	"mcp-go-server/router"
	"os"

	"github.com/gin-contrib/cors"
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

	// Add CORS middleware before other middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://mcp-node-server-ui.replit.app", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		MaxAge:           300,
	}))

	// Setup routes
	router.MCPRoutes(r.Group("/"), db)

	// Use PORT env var if set (for Replit compatibility)
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}

	log.Printf("MCP Server starting on port %s...", port)
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
