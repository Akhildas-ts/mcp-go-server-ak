package helper

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// ValidateURL validates URL format
func ValidateURL(urlStr string) error {
	_, err := url.Parse(urlStr)
	if err != nil {
		return errors.New("invalid URL format")
	}
	return nil
}

// ValidateGitHubRepoURL validates GitHub repository URL
func ValidateGitHubRepoURL(repoURL string) error {
	if err := ValidateURL(repoURL); err != nil {
		return err
	}

	// Check if it's a GitHub URL
	if !strings.Contains(repoURL, "github.com") {
		return errors.New("only GitHub repositories are supported")
	}

	// Check URL format
	githubRegex := regexp.MustCompile(`^https://github\.com/[\w\-\.]+/[\w\-\.]+(?:\.git)?/?$`)
	if !githubRegex.MatchString(repoURL) {
		return errors.New("invalid GitHub repository URL format")
	}

	return nil
}

// ValidateBranch validates branch name
func ValidateBranch(branch string) error {
	if branch == "" {
		return nil // Empty branch is valid (defaults to main)
	}

	// Git branch name validation rules
	if strings.HasPrefix(branch, "/") || strings.HasSuffix(branch, "/") {
		return errors.New("branch name cannot start or end with /")
	}

	if strings.Contains(branch, "//") {
		return errors.New("branch name cannot contain consecutive slashes")
	}

	if strings.Contains(branch, " ") {
		return errors.New("branch name cannot contain spaces")
	}

	// Check for invalid characters
	invalidChars := []string{"~", "^", ":", "?", "*", "[", "\\", "..", "@{"}
	for _, char := range invalidChars {
		if strings.Contains(branch, char) {
			return errors.New("branch name contains invalid characters")
		}
	}

	return nil
}

// ValidateSearchQuery validates search query
func ValidateSearchQuery(query string) error {
	if strings.TrimSpace(query) == "" {
		return errors.New("search query cannot be empty")
	}

	if len(query) > 500 {
		return errors.New("search query too long (max 500 characters)")
	}

	return nil
}

// ValidateLimit validates pagination limit
func ValidateLimit(limit int) error {
	if limit < 1 {
		return errors.New("limit must be at least 1")
	}

	if limit > 100 {
		return errors.New("limit cannot exceed 100")
	}

	return nil
}

// SanitizeInput sanitizes user input by removing potentially harmful characters
func SanitizeInput(input string) string {
	// Remove control characters
	input = regexp.MustCompile(`[\x00-\x1f\x7f-\x9f]`).ReplaceAllString(input, "")
	
	// Trim whitespace
	input = strings.TrimSpace(input)
	
	return input
}

// ValidateJWTSecret validates JWT secret strength
func ValidateJWTSecret(secret string) error {
	if len(secret) < 32 {
		return errors.New("JWT secret must be at least 32 characters long")
	}
	return nil
}