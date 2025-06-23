package usecase

import (
	"errors"
	"mcp-go-server/domain"
	"mcp-go-server/helper"
	"mcp-go-server/models"
	"mcp-go-server/repository"

	"github.com/jinzhu/copier"
)

// GetLoginURL returns the GitHub OAuth login URL
func GetLoginURL() (string, error) {
	url, err := repository.GetGitHubLoginURL()
	if err != nil {
		return "", errors.New("failed to generate login URL")
	}
	return url, nil
}

// HandleCallback processes GitHub OAuth callback
func HandleCallback(callbackData models.CallbackData) (domain.Session, error) {
	// Exchange code for access token
	accessToken, err := repository.ExchangeCodeForToken(callbackData.Code)
	if err != nil {
		return domain.Session{}, errors.New("failed to exchange code for token")
	}

	// Get user info from GitHub
	userInfo, err := repository.GetGitHubUserInfo(accessToken)
	if err != nil {
		return domain.Session{}, errors.New("failed to get user info")
	}

	// Create or update user
	user := domain.User{
		ID:          userInfo.ID,
		Email:       userInfo.Email,
		AccessToken: accessToken,
	}

	// Save user to storage
	err = repository.SaveUser(user)
	if err != nil {
		return domain.Session{}, errors.New("failed to save user")
	}

	// Create user response
	var userResponse domain.UserResponse
	err = copier.Copy(&userResponse, &user)
	if err != nil {
		return domain.Session{}, err
	}

	// Generate JWT token
	tokenString, err := helper.GenerateJWTToken(userResponse)
	if err != nil {
		return domain.Session{}, errors.New("failed to generate token")
	}

	return domain.Session{
		User:  userResponse,
		Token: tokenString,
	}, nil
}

// GetUserProfile retrieves user profile information
func GetUserProfile(userID string) (domain.UserResponse, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return domain.UserResponse{}, models.ErrUserNotAuthenticated
	}

	var userResponse domain.UserResponse
	err = copier.Copy(&userResponse, &user)
	if err != nil {
		return domain.UserResponse{}, err
	}

	return userResponse, nil
}
