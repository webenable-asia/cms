package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// ValkeyClient wraps the Redis client for Valkey
type ValkeyClient struct {
	client *redis.Client
	ctx    context.Context
}

// NewValkeyClient creates a new Valkey client connection
func NewValkeyClient(url string) (*ValkeyClient, error) {
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

	return &ValkeyClient{
		client: client,
		ctx:    ctx,
	}, nil
}

// Set stores a key-value pair with expiration
func (v *ValkeyClient) Set(key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	err = v.client.Set(v.ctx, key, jsonValue, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}

	return nil
}

// Get retrieves a value by key
func (v *ValkeyClient) Get(key string, dest interface{}) error {
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
func (v *ValkeyClient) Delete(key string) error {
	err := v.client.Del(v.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}

	return nil
}

// Exists checks if a key exists
func (v *ValkeyClient) Exists(key string) (bool, error) {
	count, err := v.client.Exists(v.ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %w", key, err)
	}

	return count > 0, nil
}

// SetExpiration updates the expiration time for a key
func (v *ValkeyClient) SetExpiration(key string, expiration time.Duration) error {
	err := v.client.Expire(v.ctx, key, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set expiration for key %s: %w", key, err)
	}

	return nil
}

// GetTTL returns the time to live for a key
func (v *ValkeyClient) GetTTL(key string) (time.Duration, error) {
	ttl, err := v.client.TTL(v.ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get TTL for key %s: %w", key, err)
	}

	return ttl, nil
}

// IncrementCounter increments a counter and returns the new value
func (v *ValkeyClient) IncrementCounter(key string, expiration time.Duration) (int64, error) {
	// Use pipeline for atomic operation
	pipe := v.client.Pipeline()
	incr := pipe.Incr(v.ctx, key)
	pipe.Expire(v.ctx, key, expiration)

	_, err := pipe.Exec(v.ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to increment counter %s: %w", key, err)
	}

	return incr.Val(), nil
}

// SetSession stores session data
func (v *ValkeyClient) SetSession(sessionID string, data interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return v.Set(key, data, expiration)
}

// GetSession retrieves session data
func (v *ValkeyClient) GetSession(sessionID string, dest interface{}) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return v.Get(key, dest)
}

// DeleteSession removes session data
func (v *ValkeyClient) DeleteSession(sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return v.Delete(key)
}

// CachePost caches a post for faster retrieval
func (v *ValkeyClient) CachePost(postID string, post interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("post:%s", postID)
	return v.Set(key, post, expiration)
}

// GetCachedPost retrieves a cached post
func (v *ValkeyClient) GetCachedPost(postID string, dest interface{}) error {
	key := fmt.Sprintf("post:%s", postID)
	return v.Get(key, dest)
}

// InvalidatePostCache removes a cached post
func (v *ValkeyClient) InvalidatePostCache(postID string) error {
	key := fmt.Sprintf("post:%s", postID)
	return v.Delete(key)
}

// SetRateLimit sets a rate limit counter
func (v *ValkeyClient) SetRateLimit(identifier string, limit int, window time.Duration) (bool, error) {
	key := fmt.Sprintf("rate_limit:%s", identifier)

	current, err := v.IncrementCounter(key, window)
	if err != nil {
		return false, err
	}

	return current <= int64(limit), nil
}

// ResetRateLimit removes a specific rate limit counter
func (v *ValkeyClient) ResetRateLimit(identifier string) error {
	key := fmt.Sprintf("rate_limit:%s", identifier)
	return v.Delete(key)
}

// ResetRateLimitByPattern removes all rate limit counters matching a pattern
func (v *ValkeyClient) ResetRateLimitByPattern(pattern string) error {
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
func (v *ValkeyClient) ResetAllRateLimits() error {
	return v.ResetRateLimitByPattern("*")
}

// GetRateLimitInfo returns information about a rate limit
func (v *ValkeyClient) GetRateLimitInfo(identifier string) (current int64, ttl time.Duration, err error) {
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

// Close closes the Valkey connection
func (v *ValkeyClient) Close() error {
	return v.client.Close()
}

// Health checks the health of the Valkey connection
func (v *ValkeyClient) Health() error {
	_, err := v.client.Ping(v.ctx).Result()
	return err
}

// --- Application State Management ---

// SetApplicationState stores application-wide state
func (v *ValkeyClient) SetApplicationState(key string, value interface{}, expiration time.Duration) error {
	stateKey := fmt.Sprintf("app_state:%s", key)
	return v.Set(stateKey, value, expiration)
}

// GetApplicationState retrieves application-wide state
func (v *ValkeyClient) GetApplicationState(key string, dest interface{}) error {
	stateKey := fmt.Sprintf("app_state:%s", key)
	return v.Get(stateKey, dest)
}

// SetUserState stores user-specific state
func (v *ValkeyClient) SetUserState(userID, key string, value interface{}, expiration time.Duration) error {
	stateKey := fmt.Sprintf("user_state:%s:%s", userID, key)
	return v.Set(stateKey, value, expiration)
}

// GetUserState retrieves user-specific state
func (v *ValkeyClient) GetUserState(userID, key string, dest interface{}) error {
	stateKey := fmt.Sprintf("user_state:%s:%s", userID, key)
	return v.Get(stateKey, dest)
}

// SetTemporaryData stores temporary data with auto-expiration
func (v *ValkeyClient) SetTemporaryData(key string, value interface{}, expiration time.Duration) error {
	tempKey := fmt.Sprintf("temp:%s", key)
	return v.Set(tempKey, value, expiration)
}

// GetTemporaryData retrieves temporary data
func (v *ValkeyClient) GetTemporaryData(key string, dest interface{}) error {
	tempKey := fmt.Sprintf("temp:%s", key)
	return v.Get(tempKey, dest)
}

// SetCounterWithExpiry sets a counter with expiry (useful for rate limiting, stats)
func (v *ValkeyClient) SetCounterWithExpiry(key string, value int64, expiration time.Duration) error {
	counterKey := fmt.Sprintf("counter:%s", key)
	err := v.client.Set(v.ctx, counterKey, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set counter %s: %w", key, err)
	}
	return nil
}

// GetCounter retrieves a counter value
func (v *ValkeyClient) GetCounter(key string) (int64, error) {
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

// IncrementCounterBy increments a counter by a specific amount
func (v *ValkeyClient) IncrementCounterBy(key string, increment int64, expiration time.Duration) (int64, error) {
	counterKey := fmt.Sprintf("counter:%s", key)

	pipe := v.client.Pipeline()
	incr := pipe.IncrBy(v.ctx, counterKey, increment)
	pipe.Expire(v.ctx, counterKey, expiration)

	_, err := pipe.Exec(v.ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to increment counter %s: %w", key, err)
	}

	return incr.Val(), nil
}

// GetCacheStats returns cache statistics
func (v *ValkeyClient) GetCacheStats() (map[string]interface{}, error) {
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

// SetNotification stores a notification for a user
func (v *ValkeyClient) SetNotification(userID, notificationID string, notification interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("notification:%s:%s", userID, notificationID)
	return v.Set(key, notification, expiration)
}

// GetUserNotifications retrieves all notifications for a user (simplified version)
func (v *ValkeyClient) GetUserNotifications(userID string) ([]string, error) {
	pattern := fmt.Sprintf("notification:%s:*", userID)
	keys, err := v.client.Keys(v.ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications for user %s: %w", userID, err)
	}
	return keys, nil
}

// PublishEvent publishes an event to a channel (for real-time features)
func (v *ValkeyClient) PublishEvent(channel string, message interface{}) error {
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

// --- Page Cache Management ---

// CachePage stores a full page response with headers
func (v *ValkeyClient) CachePage(cacheKey string, response []byte, contentType string, expiration time.Duration) error {
	pageData := map[string]interface{}{
		"content":      string(response),
		"content_type": contentType,
		"cached_at":    time.Now().Unix(),
	}

	key := fmt.Sprintf("page_cache:%s", cacheKey)
	return v.Set(key, pageData, expiration)
}

// GetCachedPage retrieves a cached page response
func (v *ValkeyClient) GetCachedPage(cacheKey string) ([]byte, string, error) {
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
func (v *ValkeyClient) InvalidatePageCache(pattern string) error {
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
func (v *ValkeyClient) InvalidateAllPageCache() error {
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

// CachePostsList caches the posts list with query parameters
func (v *ValkeyClient) CachePostsList(queryHash string, posts interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("posts_list:%s", queryHash)
	return v.Set(key, posts, expiration)
}

// GetCachedPostsList retrieves cached posts list
func (v *ValkeyClient) GetCachedPostsList(queryHash string, dest interface{}) error {
	key := fmt.Sprintf("posts_list:%s", queryHash)
	return v.Get(key, dest)
}

// InvalidatePostsListCache removes cached posts lists
func (v *ValkeyClient) InvalidatePostsListCache() error {
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
