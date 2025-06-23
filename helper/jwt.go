package helper

import (
	"errors"
	"mcp-go-server/database"
	"mcp-go-server/domain"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims represents JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateJWTToken generates a JWT token for user
func GenerateJWTToken(user domain.UserResponse) (string, error) {
	if database.DB == nil || database.DB.Config == nil {
		return "", errors.New("database not initialized")
	}

	// Create claims
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "mcp-server",
			Subject:   user.ID,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(database.DB.Config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWTToken validates and parses JWT token
func ValidateJWTToken(tokenString string) (string, error) {
	if database.DB == nil || database.DB.Config == nil {
		return "", errors.New("database not initialized")
	}

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(database.DB.Config.JWTSecret), nil
	})

	if err != nil {
		return "", err
	}

	// Validate token and extract claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Check if token is expired
		if claims.ExpiresAt.Time.Before(time.Now()) {
			return "", errors.New("token expired")
		}
		return claims.UserID, nil
	}

	return "", errors.New("invalid token")
}

// RefreshJWTToken refreshes an existing JWT token
func RefreshJWTToken(tokenString string) (string, error) {
	// Validate current token
	userID, err := ValidateJWTToken(tokenString)
	if err != nil {
		return "", err
	}

	// For refresh, we need user details
	// In a real implementation, you'd fetch from database
	user := domain.UserResponse{
		ID: userID,
		// Email would be fetched from database
	}

	// Generate new token
	return GenerateJWTToken(user)
}