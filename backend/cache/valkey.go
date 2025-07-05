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

// Close closes the Valkey connection
func (v *ValkeyClient) Close() error {
	return v.client.Close()
}

// Health checks the health of the Valkey connection
func (v *ValkeyClient) Health() error {
	_, err := v.client.Ping(v.ctx).Result()
	return err
}
