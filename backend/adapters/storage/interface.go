package storage

import (
	"io"
	"time"
)

// StorageAdapter defines the interface for file storage operations
type StorageAdapter interface {
	// File Operations
	Upload(file StorageFile) (*StorageResult, error)
	Download(path string) (*StorageFile, error)
	Delete(path string) error
	Exists(path string) (bool, error)

	// Directory Operations
	CreateDirectory(path string) error
	ListFiles(path string) ([]StorageInfo, error)

	// Metadata Operations
	GetMetadata(path string) (*StorageMetadata, error)
	SetMetadata(path string, metadata map[string]string) error

	// URL Operations
	GetPublicURL(path string) (string, error)
	GetSignedURL(path string, expiration time.Duration) (string, error)

	// Configuration
	Configure(config StorageConfig) error
	Health() error
}

// StorageFile represents a file to be stored or retrieved
type StorageFile struct {
	Path        string            `json:"path"`
	Content     io.Reader         `json:"-"`
	ContentType string            `json:"content_type"`
	Size        int64             `json:"size"`
	Metadata    map[string]string `json:"metadata"`
}

// StorageResult represents the result of a file upload
type StorageResult struct {
	Path       string            `json:"path"`
	URL        string            `json:"url"`
	Size       int64             `json:"size"`
	Metadata   map[string]string `json:"metadata"`
	UploadedAt time.Time         `json:"uploaded_at"`
}

// StorageInfo represents information about a stored file
type StorageInfo struct {
	Path        string            `json:"path"`
	Size        int64             `json:"size"`
	ContentType string            `json:"content_type"`
	ModifiedAt  time.Time         `json:"modified_at"`
	Metadata    map[string]string `json:"metadata"`
	IsDirectory bool              `json:"is_directory"`
}

// StorageMetadata represents metadata about a stored file
type StorageMetadata struct {
	ContentType string            `json:"content_type"`
	Size        int64             `json:"size"`
	ModifiedAt  time.Time         `json:"modified_at"`
	Custom      map[string]string `json:"custom"`
}

// StorageConfig holds configuration for storage adapters
type StorageConfig struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

// StorageType constants for supported storage types
const (
	StorageTypeLocal      = "local"
	StorageTypeS3         = "s3"
	StorageTypeGCS        = "gcs"
	StorageTypeAzureBlob  = "azure_blob"
	StorageTypeMinIO      = "minio"
)

// Common storage errors
const (
	ErrFileNotFound      = "file_not_found"
	ErrInvalidPath       = "invalid_path"
	ErrUploadFailed      = "upload_failed"
	ErrDownloadFailed    = "download_failed"
	ErrDeleteFailed      = "delete_failed"
	ErrInvalidFile       = "invalid_file"
	ErrPermissionDenied  = "permission_denied"
	ErrStorageQuotaExceeded = "storage_quota_exceeded"
)