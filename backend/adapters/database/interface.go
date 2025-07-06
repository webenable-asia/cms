package database

import (
	"context"
	"webenable-cms-backend/models"
)

// DatabaseAdapter defines the interface for database operations
type DatabaseAdapter interface {
	// Connection Management
	Connect(config DatabaseConfig) error
	Close() error
	Health() error

	// Post Operations
	CreatePost(post *models.Post) error
	GetPost(id string) (*models.Post, error)
	GetPosts(limit, offset int) ([]models.Post, error)
	UpdatePost(id string, post *models.Post) error
	DeletePost(id string) error

	// User Operations
	CreateUser(user *models.User) error
	GetUser(id string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUsers(limit, offset int) ([]models.User, error)
	UpdateUser(id string, user *models.User) error
	DeleteUser(id string) error

	// Contact Operations
	CreateContact(contact *models.Contact) error
	GetContact(id string) (*models.Contact, error)
	GetContacts(limit, offset int) ([]models.Contact, error)
	UpdateContact(id string, contact *models.Contact) error
	DeleteContact(id string) error

	// Transaction Support
	BeginTransaction() (Transaction, error)
}

// Transaction defines the interface for database transactions
type Transaction interface {
	Commit() error
	Rollback() error
	Context() context.Context
}

// DatabaseConfig holds configuration for database adapters
type DatabaseConfig struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

// DatabaseType constants for supported database types
const (
	DatabaseTypeCouchDB  = "couchdb"
	DatabaseTypePostgres = "postgres"
	DatabaseTypeMongoDB  = "mongodb"
)