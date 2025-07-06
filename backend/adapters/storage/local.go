package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LocalAdapter implements StorageAdapter for local file system
type LocalAdapter struct {
	basePath string
	baseURL  string
	config   map[string]interface{}
}

// NewLocalAdapter creates a new local file system adapter
func NewLocalAdapter(config map[string]interface{}) (StorageAdapter, error) {
	adapter := &LocalAdapter{
		config: config,
	}

	if err := adapter.Configure(StorageConfig{
		Type:   StorageTypeLocal,
		Config: config,
	}); err != nil {
		return nil, err
	}

	return adapter, nil
}

// Configure configures the local storage adapter
func (l *LocalAdapter) Configure(config StorageConfig) error {
	basePath, ok := config.Config["base_path"].(string)
	if !ok {
		basePath = "./uploads"
	}

	baseURL, ok := config.Config["base_url"].(string)
	if !ok {
		baseURL = "http://localhost:8080/uploads"
	}

	// Ensure base path exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return fmt.Errorf("failed to create base directory: %w", err)
	}

	l.basePath = basePath
	l.baseURL = baseURL
	return nil
}

// Upload uploads a file to local storage
func (l *LocalAdapter) Upload(file StorageFile) (*StorageResult, error) {
	// Sanitize file path
	cleanPath := l.sanitizePath(file.Path)
	fullPath := filepath.Join(l.basePath, cleanPath)

	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	f, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	// Copy content
	size, err := io.Copy(f, file.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to write file content: %w", err)
	}

	// Create public URL
	url := strings.TrimSuffix(l.baseURL, "/") + "/" + strings.TrimPrefix(cleanPath, "/")

	return &StorageResult{
		Path:       cleanPath,
		URL:        url,
		Size:       size,
		Metadata:   file.Metadata,
		UploadedAt: time.Now(),
	}, nil
}

// Download downloads a file from local storage
func (l *LocalAdapter) Download(path string) (*StorageFile, error) {
	cleanPath := l.sanitizePath(path)
	fullPath := filepath.Join(l.basePath, cleanPath)

	// Check if file exists
	stat, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	if stat.IsDir() {
		return nil, fmt.Errorf("path is a directory: %s", path)
	}

	// Open file
	f, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	// Detect content type
	contentType := l.detectContentType(cleanPath)

	return &StorageFile{
		Path:        cleanPath,
		Content:     f,
		ContentType: contentType,
		Size:        stat.Size(),
		Metadata:    make(map[string]string),
	}, nil
}

// Delete deletes a file from local storage
func (l *LocalAdapter) Delete(path string) error {
	cleanPath := l.sanitizePath(path)
	fullPath := filepath.Join(l.basePath, cleanPath)

	err := os.Remove(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", path)
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// Exists checks if a file exists in local storage
func (l *LocalAdapter) Exists(path string) (bool, error) {
	cleanPath := l.sanitizePath(path)
	fullPath := filepath.Join(l.basePath, cleanPath)

	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}

	return true, nil
}

// CreateDirectory creates a directory in local storage
func (l *LocalAdapter) CreateDirectory(path string) error {
	cleanPath := l.sanitizePath(path)
	fullPath := filepath.Join(l.basePath, cleanPath)

	err := os.MkdirAll(fullPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}

// ListFiles lists files in a directory
func (l *LocalAdapter) ListFiles(path string) ([]StorageInfo, error) {
	cleanPath := l.sanitizePath(path)
	fullPath := filepath.Join(l.basePath, cleanPath)

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var files []StorageInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		filePath := filepath.Join(cleanPath, entry.Name())
		contentType := ""
		if !entry.IsDir() {
			contentType = l.detectContentType(filePath)
		}

		files = append(files, StorageInfo{
			Path:        filePath,
			Size:        info.Size(),
			ContentType: contentType,
			ModifiedAt:  info.ModTime(),
			Metadata:    make(map[string]string),
			IsDirectory: entry.IsDir(),
		})
	}

	return files, nil
}

// GetMetadata retrieves metadata for a file
func (l *LocalAdapter) GetMetadata(path string) (*StorageMetadata, error) {
	cleanPath := l.sanitizePath(path)
	fullPath := filepath.Join(l.basePath, cleanPath)

	stat, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	contentType := ""
	if !stat.IsDir() {
		contentType = l.detectContentType(cleanPath)
	}

	return &StorageMetadata{
		ContentType: contentType,
		Size:        stat.Size(),
		ModifiedAt:  stat.ModTime(),
		Custom:      make(map[string]string),
	}, nil
}

// SetMetadata sets metadata for a file (local storage doesn't support custom metadata)
func (l *LocalAdapter) SetMetadata(path string, metadata map[string]string) error {
	// Local file system doesn't support custom metadata
	// In a production system, you might store metadata in a separate file or database
	return fmt.Errorf("local storage does not support custom metadata")
}

// GetPublicURL returns the public URL for a file
func (l *LocalAdapter) GetPublicURL(path string) (string, error) {
	cleanPath := l.sanitizePath(path)
	url := strings.TrimSuffix(l.baseURL, "/") + "/" + strings.TrimPrefix(cleanPath, "/")
	return url, nil
}

// GetSignedURL returns a signed URL (local storage doesn't support signed URLs)
func (l *LocalAdapter) GetSignedURL(path string, expiration time.Duration) (string, error) {
	// Local storage doesn't support signed URLs
	// Return the public URL instead
	return l.GetPublicURL(path)
}

// Health checks the health of the local storage
func (l *LocalAdapter) Health() error {
	// Check if base path is accessible
	_, err := os.Stat(l.basePath)
	if err != nil {
		return fmt.Errorf("base path not accessible: %w", err)
	}

	// Try to create a test file
	testPath := filepath.Join(l.basePath, ".health_check")
	f, err := os.Create(testPath)
	if err != nil {
		return fmt.Errorf("cannot write to storage: %w", err)
	}
	f.Close()

	// Clean up test file
	os.Remove(testPath)

	return nil
}

// sanitizePath sanitizes a file path to prevent directory traversal
func (l *LocalAdapter) sanitizePath(path string) string {
	// Remove any directory traversal attempts
	path = filepath.Clean(path)
	
	// Remove leading slashes and dots
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimPrefix(path, "./")
	path = strings.TrimPrefix(path, "../")
	
	// Replace any remaining .. with _
	path = strings.ReplaceAll(path, "..", "_")
	
	return path
}

// detectContentType detects content type based on file extension
func (l *LocalAdapter) detectContentType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	
	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".svg":  "image/svg+xml",
		".pdf":  "application/pdf",
		".txt":  "text/plain",
		".html": "text/html",
		".css":  "text/css",
		".js":   "application/javascript",
		".json": "application/json",
		".xml":  "application/xml",
		".zip":  "application/zip",
		".mp4":  "video/mp4",
		".mp3":  "audio/mpeg",
		".wav":  "audio/wav",
	}
	
	if contentType, exists := contentTypes[ext]; exists {
		return contentType
	}
	
	return "application/octet-stream"
}