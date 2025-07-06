package handlers

import (
	"webenable-cms-backend/cache"
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
