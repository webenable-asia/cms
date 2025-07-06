package adapters

import (
	"fmt"
	"webenable-cms-backend/adapters/auth"
	"webenable-cms-backend/adapters/cache"
	"webenable-cms-backend/adapters/database"
	"webenable-cms-backend/adapters/email"
	"webenable-cms-backend/adapters/storage"
	"webenable-cms-backend/config"
)

// AdapterFactory creates adapters based on configuration
type AdapterFactory struct {
	config *config.AdapterConfig
}

// NewAdapterFactory creates a new adapter factory
func NewAdapterFactory(config *config.AdapterConfig) *AdapterFactory {
	return &AdapterFactory{
		config: config,
	}
}

// CreateDatabaseAdapter creates a database adapter based on configuration
func (f *AdapterFactory) CreateDatabaseAdapter() (database.DatabaseAdapter, error) {
	switch f.config.Database.Type {
	case database.DatabaseTypeCouchDB:
		return database.NewCouchDBAdapter(f.config.GetDatabaseConfig())
	case database.DatabaseTypePostgres:
		return nil, fmt.Errorf("postgres adapter not implemented yet")
	case database.DatabaseTypeMongoDB:
		return nil, fmt.Errorf("mongodb adapter not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported database adapter type: %s", f.config.Database.Type)
	}
}

// CreateCacheAdapter creates a cache adapter based on configuration
func (f *AdapterFactory) CreateCacheAdapter() (cache.CacheAdapter, error) {
	switch f.config.Cache.Type {
	case cache.CacheTypeValkey:
		return cache.NewValkeyAdapter(f.config.GetCacheConfig())
	case cache.CacheTypeRedis:
		return nil, fmt.Errorf("redis adapter not implemented yet")
	case cache.CacheTypeMemcached:
		return nil, fmt.Errorf("memcached adapter not implemented yet")
	case cache.CacheTypeInMemory:
		return nil, fmt.Errorf("in-memory adapter not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported cache adapter type: %s", f.config.Cache.Type)
	}
}

// CreateAuthAdapter creates an auth adapter based on configuration
func (f *AdapterFactory) CreateAuthAdapter() (auth.AuthAdapter, error) {
	switch f.config.Auth.Type {
	case auth.AuthTypeJWT:
		return auth.NewJWTAdapter(f.config.GetAuthConfig())
	case auth.AuthTypeOAuth2:
		return nil, fmt.Errorf("oauth2 adapter not implemented yet")
	case auth.AuthTypeSAML:
		return nil, fmt.Errorf("saml adapter not implemented yet")
	case auth.AuthTypeBasic:
		return nil, fmt.Errorf("basic auth adapter not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported auth adapter type: %s", f.config.Auth.Type)
	}
}

// CreateEmailAdapter creates an email adapter based on configuration
func (f *AdapterFactory) CreateEmailAdapter() (email.EmailAdapter, error) {
	switch f.config.Email.Type {
	case email.EmailTypeSMTP:
		return email.NewSMTPAdapter(f.config.GetEmailConfig())
	case email.EmailTypeSendGrid:
		return nil, fmt.Errorf("sendgrid adapter not implemented yet")
	case email.EmailTypeSES:
		return nil, fmt.Errorf("ses adapter not implemented yet")
	case email.EmailTypeMailgun:
		return nil, fmt.Errorf("mailgun adapter not implemented yet")
	case email.EmailTypePostmark:
		return nil, fmt.Errorf("postmark adapter not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported email adapter type: %s", f.config.Email.Type)
	}
}

// CreateStorageAdapter creates a storage adapter based on configuration
func (f *AdapterFactory) CreateStorageAdapter() (storage.StorageAdapter, error) {
	switch f.config.Storage.Type {
	case storage.StorageTypeLocal:
		return storage.NewLocalAdapter(f.config.GetStorageConfig())
	case storage.StorageTypeS3:
		return nil, fmt.Errorf("s3 adapter not implemented yet")
	case storage.StorageTypeGCS:
		return nil, fmt.Errorf("gcs adapter not implemented yet")
	case storage.StorageTypeAzureBlob:
		return nil, fmt.Errorf("azure blob adapter not implemented yet")
	case storage.StorageTypeMinIO:
		return nil, fmt.Errorf("minio adapter not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported storage adapter type: %s", f.config.Storage.Type)
	}
}

// CreateAllAdapters creates all adapters and returns them
func (f *AdapterFactory) CreateAllAdapters() (*AdapterSet, error) {
	db, err := f.CreateDatabaseAdapter()
	if err != nil {
		return nil, fmt.Errorf("failed to create database adapter: %w", err)
	}

	cache, err := f.CreateCacheAdapter()
	if err != nil {
		return nil, fmt.Errorf("failed to create cache adapter: %w", err)
	}

	auth, err := f.CreateAuthAdapter()
	if err != nil {
		return nil, fmt.Errorf("failed to create auth adapter: %w", err)
	}

	email, err := f.CreateEmailAdapter()
	if err != nil {
		return nil, fmt.Errorf("failed to create email adapter: %w", err)
	}

	storage, err := f.CreateStorageAdapter()
	if err != nil {
		return nil, fmt.Errorf("failed to create storage adapter: %w", err)
	}

	return &AdapterSet{
		Database: db,
		Cache:    cache,
		Auth:     auth,
		Email:    email,
		Storage:  storage,
	}, nil
}

// AdapterSet holds all adapters
type AdapterSet struct {
	Database database.DatabaseAdapter
	Cache    cache.CacheAdapter
	Auth     auth.AuthAdapter
	Email    email.EmailAdapter
	Storage  storage.StorageAdapter
}

// Close closes all adapters
func (a *AdapterSet) Close() error {
	var errors []error

	if err := a.Database.Close(); err != nil {
		errors = append(errors, fmt.Errorf("failed to close database adapter: %w", err))
	}

	if err := a.Cache.Close(); err != nil {
		errors = append(errors, fmt.Errorf("failed to close cache adapter: %w", err))
	}

	if err := a.Storage.Health(); err != nil {
		// Storage adapter may not have a close method, just check health
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors closing adapters: %v", errors)
	}

	return nil
}

// Health checks the health of all adapters
func (a *AdapterSet) Health() error {
	if err := a.Database.Health(); err != nil {
		return fmt.Errorf("database adapter health check failed: %w", err)
	}

	if err := a.Cache.Health(); err != nil {
		return fmt.Errorf("cache adapter health check failed: %w", err)
	}

	if err := a.Auth.Health(); err != nil {
		return fmt.Errorf("auth adapter health check failed: %w", err)
	}

	if err := a.Email.Health(); err != nil {
		return fmt.Errorf("email adapter health check failed: %w", err)
	}

	if err := a.Storage.Health(); err != nil {
		return fmt.Errorf("storage adapter health check failed: %w", err)
	}

	return nil
}