package helper

import (
	"path/filepath"
	"strings"
)

// SplitIntoChunks splits content into chunks of approximately the specified size
func SplitIntoChunks(content string, chunkSize int) []string {
	lines := strings.Split(content, "\n")
	chunks := []string{}
	currentChunk := ""
	currentSize := 0

	for _, line := range lines {
		lineSize := len(line)
		if currentSize+lineSize > chunkSize && currentSize > 0 {
			chunks = append(chunks, currentChunk)
			currentChunk = line
			currentSize = lineSize
		} else {
			if currentSize > 0 {
				currentChunk += "\n"
			}
			currentChunk += line
			currentSize += lineSize + 1 // +1 for newline
		}
	}

	if currentSize > 0 {
		chunks = append(chunks, currentChunk)
	}

	return chunks
}

// GetLanguageFromExtension returns the programming language based on file extension
func GetLanguageFromExtension(ext string) string {
	switch strings.ToLower(ext) {
	case ".go":
		return "Go"
	case ".js", ".jsx":
		return "JavaScript"
	case ".ts", ".tsx":
		return "TypeScript"
	case ".py":
		return "Python"
	case ".java":
		return "Java"
	case ".c", ".cpp", ".h", ".hpp":
		return "C/C++"
	case ".rb":
		return "Ruby"
	case ".php":
		return "PHP"
	case ".cs":
		return "C#"
	case ".html":
		return "HTML"
	case ".css":
		return "CSS"
	case ".scss", ".sass":
		return "SCSS"
	case ".sql":
		return "SQL"
	case ".json":
		return "JSON"
	case ".xml":
		return "XML"
	case ".yaml", ".yml":
		return "YAML"
	case ".md":
		return "Markdown"
	case ".sh":
		return "Shell"
	case ".dockerfile":
		return "Dockerfile"
	case ".rs":
		return "Rust"
	case ".kt":
		return "Kotlin"
	case ".swift":
		return "Swift"
	case ".dart":
		return "Dart"
	case ".r":
		return "R"
	case ".scala":
		return "Scala"
	case ".clj":
		return "Clojure"
	case ".hs":
		return "Haskell"
	case ".lua":
		return "Lua"
	case ".perl", ".pl":
		return "Perl"
	default:
		return "Text"
	}
}

// IsBinaryFile checks if a file is binary based on its extension
func IsBinaryFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	binaryExtensions := []string{
		// Images
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".ico", ".svg", ".webp", ".tiff",
		// Archives
		".zip", ".tar", ".gz", ".rar", ".7z", ".bz2", ".xz",
		// Documents
		".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
		// Media
		".mp3", ".mp4", ".wav", ".avi", ".mov", ".wmv", ".flv", ".mkv",
		// Executables
		".so", ".dll", ".exe", ".bin", ".deb", ".rpm", ".msi",
		// Fonts
		".ttf", ".otf", ".woff", ".woff2", ".eot",
		// Other binary formats
		".db", ".sqlite", ".sqlite3", ".pyc", ".class", ".jar",
	}

	for _, binaryExt := range binaryExtensions {
		if ext == binaryExt {
			return true
		}
	}

	// Check filename for Git internal files
	if strings.Contains(filename, ".git/") ||
		strings.HasPrefix(filename, ".git") ||
		filename == "DIRC" ||
		strings.Contains(filename, "index.lock") {
		return true
	}

	return false
}

// ContainsBinaryData checks if content contains binary data
func ContainsBinaryData(content []byte) bool {
	// Check first few bytes for NULL or other binary indicators
	for i := 0; i < min(len(content), 1000); i++ {
		if content[i] == 0 || content[i] > 127 {
			return true
		}
	}
	return false
}

// ExtractRepoName extracts repository name from URL
func ExtractRepoName(repoURL string) string {
	parts := strings.Split(repoURL, "/")
	if len(parts) < 2 {
		return repoURL
	}

	repoName := parts[len(parts)-1]
	if strings.HasSuffix(repoName, ".git") {
		repoName = repoName[:len(repoName)-4]
	}

	// Return owner/repo format
	if len(parts) >= 2 {
		owner := parts[len(parts)-2]
		return owner + "/" + repoName
	}

	return repoName
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}