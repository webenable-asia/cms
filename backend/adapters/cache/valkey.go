package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// ValkeyAdapter implements CacheAdapter for Valkey/Redis
type ValkeyAdapter struct {
	client *redis.Client
	ctx    context.Context
	config map[string]interface{}
}

// NewValkeyAdapter creates a new Valkey adapter
func NewValkeyAdapter(config map[string]interface{}) (CacheAdapter, error) {
	url, ok := config["url"].(string)
	if !ok {
		return nil, fmt.Errorf("valkey url is required")
	}

	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Valkey URL: %w", err)
	}

	client := redis.NewClient(opts)
	ctx := context.Background()

	// Test the connection
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Valkey: %w", err)
	}

	log.Println("Successfully connected to Valkey")

	return &ValkeyAdapter{
		client: client,
		ctx:    ctx,
		config: config,
	}, nil
}

// Basic Operations

// Set stores a key-value pair with expiration
func (v *ValkeyAdapter) Set(key string, value interface{}, ttl time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	err = v.client.Set(v.ctx, key, jsonValue, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}

	return nil
}

// Get retrieves a value by key
func (v *ValkeyAdapter) Get(key string, dest interface{}) error {
	value, err := v.client.Get(v.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key %s not found", key)
		}
		return fmt.Errorf("failed to get key %s: %w", key, err)
	}

	err = json.Unmarshal([]byte(value), dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal value for key %s: %w", key, err)
	}

	return nil
}

// Delete removes a key
func (v *ValkeyAdapter) Delete(key string) error {
	err := v.client.Del(v.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}

	return nil
}

// Exists checks if a key exists
func (v *ValkeyAdapter) Exists(key string) (bool, error) {
	count, err := v.client.Exists(v.ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %w", key, err)
	}

	return count > 0, nil
}

// Advanced Operations

// SetExpiration updates the expiration time for a key
func (v *ValkeyAdapter) SetExpiration(key string, ttl time.Duration) error {
	err := v.client.Expire(v.ctx, key, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set expiration for key %s: %w", key, err)
	}

	return nil
}

// GetTTL returns the time to live for a key
func (v *ValkeyAdapter) GetTTL(key string) (time.Duration, error) {
	ttl, err := v.client.TTL(v.ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get TTL for key %s: %w", key, err)
	}

	return ttl, nil
}

// IncrementCounter increments a counter and returns the new value
func (v *ValkeyAdapter) IncrementCounter(key string, ttl time.Duration) (int64, error) {
	// Use pipeline for atomic operation
	pipe := v.client.Pipeline()
	incr := pipe.Incr(v.ctx, key)
	pipe.Expire(v.ctx, key, ttl)

	_, err := pipe.Exec(v.ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to increment counter %s: %w", key, err)
	}

	return incr.Val(), nil
}

// IncrementCounterBy increments a counter by a specific amount
func (v *ValkeyAdapter) IncrementCounterBy(key string, increment int64, ttl time.Duration) (int64, error) {
	pipe := v.client.Pipeline()
	incr := pipe.IncrBy(v.ctx, key, increment)
	pipe.Expire(v.ctx, key, ttl)

	_, err := pipe.Exec(v.ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to increment counter %s: %w", key, err)
	}

	return incr.Val(), nil
}

// Application-Specific Operations

// SetSession stores session data
func (v *ValkeyAdapter) SetSession(sessionID string, data interface{}, ttl time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return v.Set(key, data, ttl)
}

// GetSession retrieves session data
func (v *ValkeyAdapter) GetSession(sessionID string, dest interface{}) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return v.Get(key, dest)
}

// DeleteSession removes session data
func (v *ValkeyAdapter) DeleteSession(sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return v.Delete(key)
}

// CachePost caches a post for faster retrieval
func (v *ValkeyAdapter) CachePost(postID string, post interface{}, ttl time.Duration) error {
	key := fmt.Sprintf("post:%s", postID)
	return v.Set(key, post, ttl)
}

// GetCachedPost retrieves a cached post
func (v *ValkeyAdapter) GetCachedPost(postID string, dest interface{}) error {
	key := fmt.Sprintf("post:%s", postID)
	return v.Get(key, dest)
}

// InvalidatePostCache removes a cached post
func (v *ValkeyAdapter) InvalidatePostCache(postID string) error {
	key := fmt.Sprintf("post:%s", postID)
	return v.Delete(key)
}

// Rate Limiting

// SetRateLimit sets a rate limit counter
func (v *ValkeyAdapter) SetRateLimit(identifier string, limit int, window time.Duration) (bool, error) {
	key := fmt.Sprintf("rate_limit:%s", identifier)

	current, err := v.IncrementCounter(key, window)
	if err != nil {
		return false, err
	}

	return current <= int64(limit), nil
}

// ResetRateLimit removes a specific rate limit counter
func (v *ValkeyAdapter) ResetRateLimit(identifier string) error {
	key := fmt.Sprintf("rate_limit:%s", identifier)
	return v.Delete(key)
}

// ResetRateLimitByPattern removes all rate limit counters matching a pattern
func (v *ValkeyAdapter) ResetRateLimitByPattern(pattern string) error {
	searchKey := fmt.Sprintf("rate_limit:%s", pattern)

	keys, err := v.client.Keys(v.ctx, searchKey).Result()
	if err != nil {
		return fmt.Errorf("failed to find rate limit keys for pattern %s: %w", pattern, err)
	}

	if len(keys) > 0 {
		err = v.client.Del(v.ctx, keys...).Err()
		if err != nil {
			return fmt.Errorf("failed to delete rate limit keys: %w", err)
		}
	}

	return nil
}

// ResetAllRateLimits removes all rate limit counters
func (v *ValkeyAdapter) ResetAllRateLimits() error {
	return v.ResetRateLimitByPattern("*")
}

// GetRateLimitInfo returns information about a rate limit
func (v *ValkeyAdapter) GetRateLimitInfo(identifier string) (current int64, ttl time.Duration, err error) {
	key := fmt.Sprintf("rate_limit:%s", identifier)

	current, err = v.GetCounter(identifier)
	if err != nil {
		return 0, 0, err
	}

	ttl, err = v.GetTTL(key)
	if err != nil {
		return current, 0, err
	}

	return current, ttl, nil
}

// Page Caching

// CachePage stores a full page response with headers
func (v *ValkeyAdapter) CachePage(cacheKey string, response []byte, contentType string, ttl time.Duration) error {
	pageData := map[string]interface{}{
		"content":      string(response),
		"content_type": contentType,
		"cached_at":    time.Now().Unix(),
	}

	key := fmt.Sprintf("page_cache:%s", cacheKey)
	return v.Set(key, pageData, ttl)
}

// GetCachedPage retrieves a cached page response
func (v *ValkeyAdapter) GetCachedPage(cacheKey string) ([]byte, string, error) {
	key := fmt.Sprintf("page_cache:%s", cacheKey)

	var pageData map[string]interface{}
	err := v.Get(key, &pageData)
	if err != nil {
		return nil, "", err
	}

	content, ok := pageData["content"].(string)
	if !ok {
		return nil, "", fmt.Errorf("invalid cached content format")
	}

	contentType, ok := pageData["content_type"].(string)
	if !ok {
		contentType = "application/json" // default
	}

	return []byte(content), contentType, nil
}

// InvalidatePageCache removes cached pages based on pattern
func (v *ValkeyAdapter) InvalidatePageCache(pattern string) error {
	key := fmt.Sprintf("page_cache:%s", pattern)

	// If pattern contains wildcards, delete multiple keys
	if len(pattern) == 0 || pattern == "*" || pattern[len(pattern)-1] == '*' {
		keys, err := v.client.Keys(v.ctx, key).Result()
		if err != nil {
			return fmt.Errorf("failed to find keys for pattern %s: %w", pattern, err)
		}

		if len(keys) > 0 {
			err = v.client.Del(v.ctx, keys...).Err()
			if err != nil {
				return fmt.Errorf("failed to delete keys: %w", err)
			}
		}
		return nil
	}

	// Single key deletion
	return v.Delete(key)
}

// InvalidateAllPageCache clears all page cache
func (v *ValkeyAdapter) InvalidateAllPageCache() error {
	keys, err := v.client.Keys(v.ctx, "page_cache:*").Result()
	if err != nil {
		return fmt.Errorf("failed to find page cache keys: %w", err)
	}

	if len(keys) > 0 {
		err = v.client.Del(v.ctx, keys...).Err()
		if err != nil {
			return fmt.Errorf("failed to delete page cache keys: %w", err)
		}
	}

	return nil
}

// Posts List Caching

// CachePostsList caches the posts list with query parameters
func (v *ValkeyAdapter) CachePostsList(queryHash string, posts interface{}, ttl time.Duration) error {
	key := fmt.Sprintf("posts_list:%s", queryHash)
	return v.Set(key, posts, ttl)
}

// GetCachedPostsList retrieves cached posts list
func (v *ValkeyAdapter) GetCachedPostsList(queryHash string, dest interface{}) error {
	key := fmt.Sprintf("posts_list:%s", queryHash)
	return v.Get(key, dest)
}

// InvalidatePostsListCache removes cached posts lists
func (v *ValkeyAdapter) InvalidatePostsListCache() error {
	keys, err := v.client.Keys(v.ctx, "posts_list:*").Result()
	if err != nil {
		return fmt.Errorf("failed to find posts list cache keys: %w", err)
	}

	if len(keys) > 0 {
		err = v.client.Del(v.ctx, keys...).Err()
		if err != nil {
			return fmt.Errorf("failed to delete posts list cache keys: %w", err)
		}
	}

	return nil
}

// Application State Management

// SetApplicationState stores application-wide state
func (v *ValkeyAdapter) SetApplicationState(key string, value interface{}, ttl time.Duration) error {
	stateKey := fmt.Sprintf("app_state:%s", key)
	return v.Set(stateKey, value, ttl)
}

// GetApplicationState retrieves application-wide state
func (v *ValkeyAdapter) GetApplicationState(key string, dest interface{}) error {
	stateKey := fmt.Sprintf("app_state:%s", key)
	return v.Get(stateKey, dest)
}

// SetUserState stores user-specific state
func (v *ValkeyAdapter) SetUserState(userID, key string, value interface{}, ttl time.Duration) error {
	stateKey := fmt.Sprintf("user_state:%s:%s", userID, key)
	return v.Set(stateKey, value, ttl)
}

// GetUserState retrieves user-specific state
func (v *ValkeyAdapter) GetUserState(userID, key string, dest interface{}) error {
	stateKey := fmt.Sprintf("user_state:%s:%s", userID, key)
	return v.Get(stateKey, dest)
}

// SetTemporaryData stores temporary data with auto-expiration
func (v *ValkeyAdapter) SetTemporaryData(key string, value interface{}, ttl time.Duration) error {
	tempKey := fmt.Sprintf("temp:%s", key)
	return v.Set(tempKey, value, ttl)
}

// GetTemporaryData retrieves temporary data
func (v *ValkeyAdapter) GetTemporaryData(key string, dest interface{}) error {
	tempKey := fmt.Sprintf("temp:%s", key)
	return v.Get(tempKey, dest)
}

// Counter Operations

// SetCounterWithExpiry sets a counter with expiry
func (v *ValkeyAdapter) SetCounterWithExpiry(key string, value int64, ttl time.Duration) error {
	counterKey := fmt.Sprintf("counter:%s", key)
	err := v.client.Set(v.ctx, counterKey, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set counter %s: %w", key, err)
	}
	return nil
}

// GetCounter retrieves a counter value
func (v *ValkeyAdapter) GetCounter(key string) (int64, error) {
	counterKey := fmt.Sprintf("counter:%s", key)
	value, err := v.client.Get(v.ctx, counterKey).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // Return 0 if key doesn't exist
		}
		return 0, fmt.Errorf("failed to get counter %s: %w", key, err)
	}
	return value, nil
}

// Notifications

// SetNotification stores a notification for a user
func (v *ValkeyAdapter) SetNotification(userID, notificationID string, notification interface{}, ttl time.Duration) error {
	key := fmt.Sprintf("notification:%s:%s", userID, notificationID)
	return v.Set(key, notification, ttl)
}

// GetUserNotifications retrieves all notifications for a user
func (v *ValkeyAdapter) GetUserNotifications(userID string) ([]string, error) {
	pattern := fmt.Sprintf("notification:%s:*", userID)
	keys, err := v.client.Keys(v.ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications for user %s: %w", userID, err)
	}
	return keys, nil
}

// PublishEvent publishes an event to a channel
func (v *ValkeyAdapter) PublishEvent(channel string, message interface{}) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = v.client.Publish(v.ctx, channel, jsonMessage).Err()
	if err != nil {
		return fmt.Errorf("failed to publish to channel %s: %w", channel, err)
	}

	return nil
}

// Health & Stats

// Health checks the health of the Valkey connection
func (v *ValkeyAdapter) Health() error {
	_, err := v.client.Ping(v.ctx).Result()
	return err
}

// GetStats returns cache statistics
func (v *ValkeyAdapter) GetStats() (map[string]interface{}, error) {
	info, err := v.client.Info(v.ctx, "memory", "stats").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache stats: %w", err)
	}

	// Parse basic info
	stats := map[string]interface{}{
		"connected":   true,
		"info":        info,
		"ping_result": "PONG",
	}

	// Get database size
	dbSize, err := v.client.DBSize(v.ctx).Result()
	if err == nil {
		stats["db_size"] = dbSize
	}

	return stats, nil
}

// Close closes the Valkey connection
func (v *ValkeyAdapter) Close() error {
	return v.client.Close()
}

// NewRedisAdapter creates a new Redis adapter (alias for Valkey)
func NewRedisAdapter(config map[string]interface{}) (CacheAdapter, error) {
	return NewValkeyAdapter(config)
}