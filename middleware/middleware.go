package middleware

import (
	"mcp-go-server/helper"
	"mcp-go-server/models"
	"mcp-go-server/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errRes := response.ErrorClientResponse(http.StatusUnauthorized, "Authorization header is required", nil)
			c.JSON(http.StatusUnauthorized, errRes)
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			errRes := response.ErrorClientResponse(http.StatusUnauthorized, "Invalid authorization format", nil)
			c.JSON(http.StatusUnauthorized, errRes)
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			errRes := response.ErrorClientResponse(http.StatusUnauthorized, "Token is required", nil)
			c.JSON(http.StatusUnauthorized, errRes)
			c.Abort()
			return
		}

		// Validate and parse token
		userID, err := helper.ValidateJWTToken(token)
		if err != nil {
			errRes := response.ErrorClientResponse(http.StatusUnauthorized, "Invalid token", err.Error())
			c.JSON(http.StatusUnauthorized, errRes)
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set(models.UserIDKey, userID)
		c.Next()
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
