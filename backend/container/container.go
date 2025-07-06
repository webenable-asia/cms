package container

import (
	"fmt"
	"webenable-cms-backend/adapters"
	"webenable-cms-backend/adapters/auth"
	"webenable-cms-backend/adapters/cache"
	"webenable-cms-backend/adapters/database"
	"webenable-cms-backend/adapters/email"
	"webenable-cms-backend/adapters/storage"
	"webenable-cms-backend/config"
)

// Container holds all application dependencies
type Container struct {
	adapters *adapters.AdapterSet
	config   *config.AdapterConfig
}

// NewContainer creates a new service container
func NewContainer(config *config.AdapterConfig) (*Container, error) {
	factory := adapters.NewAdapterFactory(config)
	
	adapters, err := factory.CreateAllAdapters()
	if err != nil {
		return nil, fmt.Errorf("failed to create adapters: %w", err)
	}

	container := &Container{
		adapters: adapters,
		config:   config,
	}

	return container, nil
}

// Database returns the database adapter
func (c *Container) Database() database.DatabaseAdapter {
	return c.adapters.Database
}

// Cache returns the cache adapter
func (c *Container) Cache() cache.CacheAdapter {
	return c.adapters.Cache
}

// Auth returns the auth adapter
func (c *Container) Auth() auth.AuthAdapter {
	return c.adapters.Auth
}

// Email returns the email adapter
func (c *Container) Email() email.EmailAdapter {
	return c.adapters.Email
}

// Storage returns the storage adapter
func (c *Container) Storage() storage.StorageAdapter {
	return c.adapters.Storage
}

// Config returns the adapter configuration
func (c *Container) Config() *config.AdapterConfig {
	return c.config
}

// Close closes all adapters
func (c *Container) Close() error {
	return c.adapters.Close()
}

// Health checks the health of all adapters
func (c *Container) Health() error {
	return c.adapters.Health()
}

// GetAdapterInfo returns information about all configured adapters
func (c *Container) GetAdapterInfo() map[string]interface{} {
	return map[string]interface{}{
		"database": c.config.Database.Type,
		"cache":    c.config.Cache.Type,
		"auth":     c.config.Auth.Type,
		"email":    c.config.Email.Type,
		"storage":  c.config.Storage.Type,
	}
}

// ServiceProvider interface for services that need container dependencies
type ServiceProvider interface {
	Register(container *Container) error
	Boot(container *Container) error
}

// RegisterService registers a service with the container
func (c *Container) RegisterService(service ServiceProvider) error {
	if err := service.Register(c); err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	if err := service.Boot(c); err != nil {
		return fmt.Errorf("failed to boot service: %w", err)
	}

	return nil
}