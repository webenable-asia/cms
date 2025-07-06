package handlers

import (
	"net/http"
	"strconv"
	"webenable-cms-backend/cache"
	"webenable-cms-backend/container"
	"webenable-cms-backend/middleware"
)

// Handlers holds dependencies for all handlers
type Handlers struct {
	Cache     *cache.ValkeyClient
	Container *container.Container
}

// NewHandlers creates a new handlers instance
func NewHandlers(valkeyClient *cache.ValkeyClient) *Handlers {
	return &Handlers{
		Cache: valkeyClient,
	}
}

// NewHandlersWithContainer creates a new handlers instance with service container
func NewHandlersWithContainer(valkeyClient *cache.ValkeyClient, container *container.Container) *Handlers {
	return &Handlers{
		Cache:     valkeyClient,
		Container: container,
	}
}

// Global variables for backward compatibility
var (
	globalCache       *cache.ValkeyClient
	globalRateLimiter *middleware.RateLimiter
	globalContainer   *container.Container
)

// SetGlobalCache sets the global cache instance
func SetGlobalCache(cache *cache.ValkeyClient) {
	globalCache = cache
}

// SetGlobalRateLimiter sets the global rate limiter instance
func SetGlobalRateLimiter(rateLimiter *middleware.RateLimiter) {
	globalRateLimiter = rateLimiter
}

// SetServiceContainer sets the global service container instance
func SetServiceContainer(container *container.Container) {
	globalContainer = container
}

// GetServiceContainer returns the global service container instance
func GetServiceContainer() *container.Container {
	return globalContainer
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
