package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mcp-go-server/database"
	"mcp-go-server/domain"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	githubOauthConfig *oauth2.Config
	userStore         = make(map[string]domain.User) // In-memory store (replace with DB in production)
)

// Initialize OAuth config
func init() {
	if database.DB != nil && database.DB.Config != nil {
		githubOauthConfig = &oauth2.Config{
			ClientID:     database.DB.Config.GitHubClientID,
			ClientSecret: database.DB.Config.GitHubClientSecret,
			Scopes:       []string{"user:email", "repo"},
			Endpoint:     github.Endpoint,
			RedirectURL:  database.DB.Config.GitHubOAuthRedirectURL,
		}
	}
}

// InitializeOAuth initializes OAuth configuration
func InitializeOAuth() {
	if database.DB != nil && database.DB.Config != nil {
		githubOauthConfig = &oauth2.Config{
			ClientID:     database.DB.Config.GitHubClientID,
			ClientSecret: database.DB.Config.GitHubClientSecret,
			Scopes:       []string{"user:email", "repo"},
			Endpoint:     github.Endpoint,
			RedirectURL:  database.DB.Config.GitHubOAuthRedirectURL,
		}
	}
}

// GetGitHubLoginURL returns the GitHub OAuth login URL
func GetGitHubLoginURL() (string, error) {
	if githubOauthConfig == nil {
		InitializeOAuth()
	}
	if githubOauthConfig == nil {
		return "", errors.New("GitHub OAuth config not initialized")
	}

	state := "random-state" // TODO: Generate secure random state
	url := githubOauthConfig.AuthCodeURL(state)
	return url, nil
}

// ExchangeCodeForToken exchanges authorization code for access token
func ExchangeCodeForToken(code string) (string, error) {
	if githubOauthConfig == nil {
		InitializeOAuth()
	}
	if githubOauthConfig == nil {
		return "", errors.New("GitHub OAuth config not initialized")
	}

	token, err := githubOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code: %w", err)
	}

	return token.AccessToken, nil
}

// GitHubUser represents GitHub user info
type GitHubUser struct {
	ID    string `json:"login"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// GetGitHubUserInfo retrieves user info from GitHub API
func GetGitHubUserInfo(accessToken string) (GitHubUser, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return GitHubUser{}, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GitHubUser{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GitHubUser{}, fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return GitHubUser{}, err
	}

	// Get email if not provided in user info
	if user.Email == "" {
		email, err := getGitHubUserEmail(accessToken)
		if err == nil {
			user.Email = email
		}
	}

	return user, nil
}

// getGitHubUserEmail retrieves user email from GitHub API
func getGitHubUserEmail(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emails []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	for _, email := range emails {
		if email.Primary {
			return email.Email, nil
		}
	}

	if len(emails) > 0 {
		return emails[0].Email, nil
	}

	return "", errors.New("no email found")
}

// SaveUser saves user to storage
func SaveUser(user domain.User) error {
	userStore[user.ID] = user
	return nil
}

// GetUserByID retrieves user by ID
func GetUserByID(userID string) (domain.User, error) {
	user, exists := userStore[userID]
	if !exists {
		return domain.User{}, errors.New("user not found")
	}
	return user, nil
}

// ValidateAccessToken validates GitHub access token
func ValidateAccessToken(accessToken string) (bool, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}
