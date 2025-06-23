package domain

// User represents an authenticated user
type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
	CreatedAt   string `json:"created_at"`
}

// Session represents a user session with JWT token
type Session struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// UserResponse represents user data in API responses
type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// CodeChunk represents a code chunk stored in vector database
type CodeChunk struct {
	ID         string    `json:"id"`
	Content    string    `json:"content"`
	FilePath   string    `json:"file_path"`
	Repository string    `json:"repository"`
	Branch     string    `json:"branch"`
	Language   string    `json:"language"`
	Embedding  []float32 `json:"embedding"`
}

// SearchResult represents a search result with score
type SearchResult struct {
	CodeChunk
	Score float32 `json:"score"`
}

// Repository represents a Git repository
type Repository struct {
	URL        string `json:"url"`
	Name       string `json:"name"`
	Owner      string `json:"owner"`
	Branch     string `json:"branch"`
	IndexedAt  string `json:"indexed_at"`
	FileCount  int    `json:"file_count"`
	ChunkCount int    `json:"chunk_count"`
}