package handlers

import (
	"net/http"
	"strconv"
	"webenable-cms-backend/cache"
	"webenable-cms-backend/middleware"
)

// Handlers holds dependencies for all handlers
type Handlers struct {
	Cache *cache.ValkeyClient
}

// NewHandlers creates a new handlers instance
func NewHandlers(valkeyClient *cache.ValkeyClient) *Handlers {
	return &Handlers{
		Cache: valkeyClient,
	}
}

// Global variables for backward compatibility
var (
	globalCache       *cache.ValkeyClient
	globalRateLimiter *middleware.RateLimiter
)

// SetGlobalCache sets the global cache instance
func SetGlobalCache(cache *cache.ValkeyClient) {
	globalCache = cache
}

// SetGlobalRateLimiter sets the global rate limiter instance
func SetGlobalRateLimiter(rateLimiter *middleware.RateLimiter) {
	globalRateLimiter = rateLimiter
}

// getPaginationParams extracts pagination parameters from request
func getPaginationParams(r *http.Request) (page, limit int) {
	page = 1
	limit = 10

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	return page, limit
}
