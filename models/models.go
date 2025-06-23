package models

import "errors"

// Constants
const UserIDKey = "user_id"

// Custom errors
var (
	ErrEmailNotFound        = errors.New("email not found")
	ErrPasswordIncorrect    = errors.New("password is not correct")
	ErrUserNotAuthenticated = errors.New("user not authenticated")
	ErrRepositoryNotFound   = errors.New("repository not found")
	ErrInvalidBranch        = errors.New("invalid branch")
)

// Auth models
type CallbackData struct {
	Code string `json:"code" validate:"required"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// Search models
type SearchRequest struct {
	Query      string `json:"query" validate:"required,min=1"`
	Repository string `json:"repository" validate:"required"`
	Branch     string `json:"branch"`
	Limit      int    `json:"limit"`
}

type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
}

type SearchResult struct {
	Content    string  `json:"content"`
	FilePath   string  `json:"file_path"`
	Repository string  `json:"repository"`
	Branch     string  `json:"branch"`
	Language   string  `json:"language"`
	Score      float32 `json:"score"`
}

type SearchWithSummaryResponse struct {
	Summary string         `json:"summary"`
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
}

// Repository indexing models
type IndexRequest struct {
	RepoURL string `json:"repo_url" validate:"required,url"`
	Branch  string `json:"branch"`
}

type IndexResponse struct {
	Repository string `json:"repository"`
	Branch     string `json:"branch"`
	FileCount  int    `json:"file_count"`
	ChunkCount int    `json:"chunk_count"`
	Status     string `json:"status"`
}

type IndexProgress struct {
	Repository   string `json:"repository"`
	Branch       string `json:"branch"`
	Status       string `json:"status"` // "processing", "completed", "failed"
	CurrentFile  int    `json:"current_file"`
	TotalFiles   int    `json:"total_files"`
	CurrentChunk int    `json:"current_chunk"`
	TotalChunks  int    `json:"total_chunks"`
	Message      string `json:"message"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time,omitempty"`
}

type RepositoryInfo struct {
	Name       string `json:"name"`
	Owner      string `json:"owner"`
	URL        string `json:"url"`
	Branch     string `json:"branch"`
	IndexedAt  string `json:"indexed_at"`
	FileCount  int    `json:"file_count"`
	ChunkCount int    `json:"chunk_count"`
}

// Code chunk model
type CodeChunk struct {
	ID         string    `json:"id"`
	Content    string    `json:"content"`
	FilePath   string    `json:"file_path"`
	Repository string    `json:"repository"`
	Branch     string    `json:"branch"`
	Language   string    `json:"language"`
	Embedding  []float32 `json:"embedding"`
}
