package cache

import (
	"time"
)

// CacheAdapter defines the interface for cache operations
type CacheAdapter interface {
	// Basic Operations
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string, dest interface{}) error
	Delete(key string) error
	Exists(key string) (bool, error)

	// Advanced Operations
	SetExpiration(key string, ttl time.Duration) error
	GetTTL(key string) (time.Duration, error)
	IncrementCounter(key string, ttl time.Duration) (int64, error)
	IncrementCounterBy(key string, increment int64, ttl time.Duration) (int64, error)

	// Application-Specific Operations
	SetSession(sessionID string, data interface{}, ttl time.Duration) error
	GetSession(sessionID string, dest interface{}) error
	DeleteSession(sessionID string) error

	CachePost(postID string, post interface{}, ttl time.Duration) error
	GetCachedPost(postID string, dest interface{}) error
	InvalidatePostCache(postID string) error

	// Rate Limiting
	SetRateLimit(identifier string, limit int, window time.Duration) (bool, error)
	ResetRateLimit(identifier string) error
	ResetRateLimitByPattern(pattern string) error
	ResetAllRateLimits() error
	GetRateLimitInfo(identifier string) (current int64, ttl time.Duration, err error)

	// Page Caching
	CachePage(cacheKey string, response []byte, contentType string, ttl time.Duration) error
	GetCachedPage(cacheKey string) ([]byte, string, error)
	InvalidatePageCache(pattern string) error
	InvalidateAllPageCache() error

	// Posts List Caching
	CachePostsList(queryHash string, posts interface{}, ttl time.Duration) error
	GetCachedPostsList(queryHash string, dest interface{}) error
	InvalidatePostsListCache() error

	// Application State Management
	SetApplicationState(key string, value interface{}, ttl time.Duration) error
	GetApplicationState(key string, dest interface{}) error
	SetUserState(userID, key string, value interface{}, ttl time.Duration) error
	GetUserState(userID, key string, dest interface{}) error
	SetTemporaryData(key string, value interface{}, ttl time.Duration) error
	GetTemporaryData(key string, dest interface{}) error

	// Counter Operations
	SetCounterWithExpiry(key string, value int64, ttl time.Duration) error
	GetCounter(key string) (int64, error)

	// Notifications
	SetNotification(userID, notificationID string, notification interface{}, ttl time.Duration) error
	GetUserNotifications(userID string) ([]string, error)
	PublishEvent(channel string, message interface{}) error

	// Health & Stats
	Health() error
	GetStats() (map[string]interface{}, error)
	Close() error
}

// CacheConfig holds configuration for cache adapters
type CacheConfig struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

// CacheType constants for supported cache types
const (
	CacheTypeValkey     = "valkey"
	CacheTypeRedis      = "redis"
	CacheTypeMemcached  = "memcached"
	CacheTypeInMemory   = "inmemory"
)