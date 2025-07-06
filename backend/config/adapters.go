package config

import (
	"os"
	"strings"
)

// AdapterConfig holds configuration for all adapters
type AdapterConfig struct {
	Database DatabaseAdapterConfig `json:"database"`
	Cache    CacheAdapterConfig    `json:"cache"`
	Auth     AuthAdapterConfig     `json:"auth"`
	Email    EmailAdapterConfig    `json:"email"`
	Storage  StorageAdapterConfig  `json:"storage"`
}

// DatabaseAdapterConfig holds configuration for database adapters
type DatabaseAdapterConfig struct {
	Type   string                 `json:"type"`   // "couchdb", "postgres", "mongodb"
	Config map[string]interface{} `json:"config"`
}

// CacheAdapterConfig holds configuration for cache adapters
type CacheAdapterConfig struct {
	Type   string                 `json:"type"`   // "valkey", "redis", "memcached"
	Config map[string]interface{} `json:"config"`
}

// AuthAdapterConfig holds configuration for auth adapters
type AuthAdapterConfig struct {
	Type   string                 `json:"type"`   // "jwt", "oauth2", "saml"
	Config map[string]interface{} `json:"config"`
}

// EmailAdapterConfig holds configuration for email adapters
type EmailAdapterConfig struct {
	Type   string                 `json:"type"`   // "smtp", "sendgrid", "ses"
	Config map[string]interface{} `json:"config"`
}

// StorageAdapterConfig holds configuration for storage adapters
type StorageAdapterConfig struct {
	Type   string                 `json:"type"`   // "local", "s3", "gcs"
	Config map[string]interface{} `json:"config"`
}

// InitAdapterConfig initializes adapter configuration from environment variables
func InitAdapterConfig() *AdapterConfig {
	return &AdapterConfig{
		Database: DatabaseAdapterConfig{
			Type: getEnvOrDefault("DATABASE_ADAPTER", "couchdb"),
			Config: map[string]interface{}{
				"url": getEnvOrDefault("COUCHDB_URL", "http://admin:password@localhost:5984/"),
			},
		},
		Cache: CacheAdapterConfig{
			Type: getEnvOrDefault("CACHE_ADAPTER", "valkey"),
			Config: map[string]interface{}{
				"url": getEnvOrDefault("VALKEY_URL", "valkey://valkeypassword@localhost:6379"),
			},
		},
		Auth: AuthAdapterConfig{
			Type: getEnvOrDefault("AUTH_ADAPTER", "jwt"),
			Config: map[string]interface{}{
				"secret":     getRequiredEnv("JWT_SECRET"),
				"expiration": getEnvOrDefault("JWT_EXPIRATION", "24h"),
			},
		},
		Email: EmailAdapterConfig{
			Type: getEnvOrDefault("EMAIL_ADAPTER", "smtp"),
			Config: map[string]interface{}{
				"host":     getEnvOrDefault("SMTP_HOST", "localhost"),
				"port":     getEnvOrDefault("SMTP_PORT", "1025"),
				"username": getEnvOrDefault("SMTP_USER", "hello@webenable.asia"),
				"password": os.Getenv("SMTP_PASS"),
				"from":     getEnvOrDefault("SMTP_FROM", "hello@webenable.asia"),
			},
		},
		Storage: StorageAdapterConfig{
			Type: getEnvOrDefault("STORAGE_ADAPTER", "local"),
			Config: map[string]interface{}{
				"base_path": getEnvOrDefault("STORAGE_BASE_PATH", "./uploads"),
				"base_url":  getEnvOrDefault("STORAGE_BASE_URL", "http://localhost:8080/uploads"),
			},
		},
	}
}

// GetDatabaseConfig returns database-specific configuration
func (c *AdapterConfig) GetDatabaseConfig() map[string]interface{} {
	switch c.Database.Type {
	case "couchdb":
		return map[string]interface{}{
			"url": c.Database.Config["url"],
		}
	case "postgres":
		return map[string]interface{}{
			"host":     getEnvOrDefault("POSTGRES_HOST", "localhost"),
			"port":     getEnvOrDefault("POSTGRES_PORT", "5432"),
			"user":     getEnvOrDefault("POSTGRES_USER", "postgres"),
			"password": getEnvOrDefault("POSTGRES_PASSWORD", "password"),
			"dbname":   getEnvOrDefault("POSTGRES_DB", "webenable_cms"),
			"sslmode":  getEnvOrDefault("POSTGRES_SSLMODE", "disable"),
		}
	case "mongodb":
		return map[string]interface{}{
			"uri":      getEnvOrDefault("MONGODB_URI", "mongodb://localhost:27017"),
			"database": getEnvOrDefault("MONGODB_DATABASE", "webenable_cms"),
		}
	default:
		return c.Database.Config
	}
}

// GetCacheConfig returns cache-specific configuration
func (c *AdapterConfig) GetCacheConfig() map[string]interface{} {
	switch c.Cache.Type {
	case "valkey", "redis":
		return map[string]interface{}{
			"url": c.Cache.Config["url"],
		}
	case "memcached":
		return map[string]interface{}{
			"servers": strings.Split(getEnvOrDefault("MEMCACHED_SERVERS", "localhost:11211"), ","),
		}
	case "inmemory":
		return map[string]interface{}{
			"max_size": getEnvOrDefault("INMEMORY_CACHE_MAX_SIZE", "100MB"),
		}
	default:
		return c.Cache.Config
	}
}

// GetAuthConfig returns auth-specific configuration
func (c *AdapterConfig) GetAuthConfig() map[string]interface{} {
	switch c.Auth.Type {
	case "jwt":
		return map[string]interface{}{
			"secret":     c.Auth.Config["secret"],
			"expiration": c.Auth.Config["expiration"],
		}
	case "oauth2":
		return map[string]interface{}{
			"client_id":     getEnvOrDefault("OAUTH2_CLIENT_ID", ""),
			"client_secret": getEnvOrDefault("OAUTH2_CLIENT_SECRET", ""),
			"redirect_url":  getEnvOrDefault("OAUTH2_REDIRECT_URL", ""),
			"scopes":        strings.Split(getEnvOrDefault("OAUTH2_SCOPES", "read,write"), ","),
		}
	case "saml":
		return map[string]interface{}{
			"entity_id":     getEnvOrDefault("SAML_ENTITY_ID", ""),
			"sso_url":       getEnvOrDefault("SAML_SSO_URL", ""),
			"certificate":   getEnvOrDefault("SAML_CERTIFICATE", ""),
		}
	default:
		return c.Auth.Config
	}
}

// GetEmailConfig returns email-specific configuration
func (c *AdapterConfig) GetEmailConfig() map[string]interface{} {
	switch c.Email.Type {
	case "smtp":
		return map[string]interface{}{
			"host":     c.Email.Config["host"],
			"port":     c.Email.Config["port"],
			"username": c.Email.Config["username"],
			"password": c.Email.Config["password"],
			"from":     c.Email.Config["from"],
		}
	case "sendgrid":
		return map[string]interface{}{
			"api_key": getEnvOrDefault("SENDGRID_API_KEY", ""),
			"from":    getEnvOrDefault("SENDGRID_FROM", "hello@webenable.asia"),
		}
	case "ses":
		return map[string]interface{}{
			"region":     getEnvOrDefault("AWS_REGION", "us-east-1"),
			"access_key": getEnvOrDefault("AWS_ACCESS_KEY_ID", ""),
			"secret_key": getEnvOrDefault("AWS_SECRET_ACCESS_KEY", ""),
			"from":       getEnvOrDefault("SES_FROM", "hello@webenable.asia"),
		}
	default:
		return c.Email.Config
	}
}

// GetStorageConfig returns storage-specific configuration
func (c *AdapterConfig) GetStorageConfig() map[string]interface{} {
	switch c.Storage.Type {
	case "local":
		return map[string]interface{}{
			"base_path": c.Storage.Config["base_path"],
			"base_url":  c.Storage.Config["base_url"],
		}
	case "s3":
		return map[string]interface{}{
			"bucket":     getEnvOrDefault("S3_BUCKET", ""),
			"region":     getEnvOrDefault("AWS_REGION", "us-east-1"),
			"access_key": getEnvOrDefault("AWS_ACCESS_KEY_ID", ""),
			"secret_key": getEnvOrDefault("AWS_SECRET_ACCESS_KEY", ""),
		}
	case "gcs":
		return map[string]interface{}{
			"bucket":         getEnvOrDefault("GCS_BUCKET", ""),
			"project_id":     getEnvOrDefault("GCS_PROJECT_ID", ""),
			"credentials":    getEnvOrDefault("GCS_CREDENTIALS", ""),
		}
	default:
		return c.Storage.Config
	}
}

// Helper function to get required environment variable
func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("Required environment variable " + key + " is not set")
	}
	return value
}