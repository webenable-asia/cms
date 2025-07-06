package middleware

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"webenable-cms-backend/cache"
)

// ResponseWriter wrapper to capture response
type responseWriter struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
	written    bool
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.written = true
	}
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// PageCacheConfig holds configuration for page caching
type PageCacheConfig struct {
	ValkeyClient    *cache.ValkeyClient
	DefaultTTL      time.Duration
	SkipMethods     []string
	SkipPaths       []string
	SkipQueryParams []string
	CachePrivate    bool // Cache responses for authenticated users
}

// NewPageCache creates a new page cache middleware
func NewPageCache(valkeyClient *cache.ValkeyClient) *PageCacheConfig {
	return &PageCacheConfig{
		ValkeyClient:    valkeyClient,
		DefaultTTL:      15 * time.Minute, // Default 15 minutes
		SkipMethods:     []string{"POST", "PUT", "DELETE", "PATCH"},
		SkipPaths:       []string{"/api/auth/", "/api/users/", "/swagger/"},
		SkipQueryParams: []string{"_", "timestamp", "nocache"},
		CachePrivate:    false, // Don't cache authenticated requests by default
	}
}

// generateCacheKey creates a unique cache key for the request
func (pc *PageCacheConfig) generateCacheKey(r *http.Request) string {
	// Start with method and path
	parts := []string{r.Method, r.URL.Path}

	// Add query parameters (sorted for consistency)
	if r.URL.RawQuery != "" {
		query := r.URL.Query()

		// Remove skip parameters
		for _, param := range pc.SkipQueryParams {
			query.Del(param)
		}

		if len(query) > 0 {
			// Sort parameters for consistent keys
			cleanQuery := url.Values{}
			for key, values := range query {
				cleanQuery[key] = values
			}
			parts = append(parts, cleanQuery.Encode())
		}
	}

	// Add user context if caching private responses
	if pc.CachePrivate {
		if claims, ok := r.Context().Value("user").(*Claims); ok {
			parts = append(parts, "user:"+claims.Username)
		}
	}

	// Create hash of the key parts
	keyString := strings.Join(parts, "|")
	hash := md5.Sum([]byte(keyString))
	return fmt.Sprintf("%x", hash)
}

// shouldCache determines if a request should be cached
func (pc *PageCacheConfig) shouldCache(r *http.Request) bool {
	// Skip non-GET requests by default
	for _, method := range pc.SkipMethods {
		if r.Method == method {
			return false
		}
	}

	// Skip certain paths
	for _, path := range pc.SkipPaths {
		if strings.HasPrefix(r.URL.Path, path) {
			return false
		}
	}

	// Skip if authenticated and not caching private responses
	if !pc.CachePrivate {
		if _, ok := r.Context().Value("user").(*Claims); ok {
			return false
		}
	}

	// Check for no-cache headers
	if r.Header.Get("Cache-Control") == "no-cache" {
		return false
	}

	// Check for nocache query parameter
	if r.URL.Query().Get("nocache") != "" {
		return false
	}

	return true
}

// shouldCacheResponse determines if a response should be cached based on status code
func (pc *PageCacheConfig) shouldCacheResponse(statusCode int) bool {
	// Only cache successful responses
	return statusCode >= 200 && statusCode < 300
}

// PageCacheMiddleware returns the page cache middleware function
func (pc *PageCacheConfig) PageCacheMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if we should cache this request
			if !pc.shouldCache(r) {
				next.ServeHTTP(w, r)
				return
			}

			// Generate cache key
			cacheKey := pc.generateCacheKey(r)

			// Try to get from cache first
			cachedContent, contentType, err := pc.ValkeyClient.GetCachedPage(cacheKey)
			if err == nil {
				// Cache hit - serve from cache
				w.Header().Set("Content-Type", contentType)
				w.Header().Set("X-Cache", "HIT")
				w.Header().Set("X-Cache-Key", cacheKey)
				w.Write(cachedContent)
				return
			}

			// Cache miss - generate response
			rw := &responseWriter{
				ResponseWriter: w,
				body:           &bytes.Buffer{},
				statusCode:     200,
			}

			// Process the request
			next.ServeHTTP(rw, r)

			// Cache the response if appropriate
			if pc.shouldCacheResponse(rw.statusCode) {
				contentType := rw.Header().Get("Content-Type")
				if contentType == "" {
					contentType = "application/json"
				}

				// Store in cache
				go func() {
					err := pc.ValkeyClient.CachePage(cacheKey, rw.body.Bytes(), contentType, pc.DefaultTTL)
					if err != nil {
						// Log error but don't fail the request
						fmt.Printf("Failed to cache page %s: %v\n", cacheKey, err)
					}
				}()
			}

			// Add cache headers
			w.Header().Set("X-Cache", "MISS")
			w.Header().Set("X-Cache-Key", cacheKey)
		})
	}
}

// WithTTL sets a custom TTL for the cache
func (pc *PageCacheConfig) WithTTL(ttl time.Duration) *PageCacheConfig {
	pc.DefaultTTL = ttl
	return pc
}

// WithSkipPaths adds paths to skip caching
func (pc *PageCacheConfig) WithSkipPaths(paths ...string) *PageCacheConfig {
	pc.SkipPaths = append(pc.SkipPaths, paths...)
	return pc
}

// WithCachePrivate enables caching of authenticated requests
func (pc *PageCacheConfig) WithCachePrivate(enabled bool) *PageCacheConfig {
	pc.CachePrivate = enabled
	return pc
}

// InvalidateCache provides methods to invalidate specific cache entries
type CacheInvalidator struct {
	ValkeyClient *cache.ValkeyClient
}

// NewCacheInvalidator creates a new cache invalidator
func NewCacheInvalidator(valkeyClient *cache.ValkeyClient) *CacheInvalidator {
	return &CacheInvalidator{
		ValkeyClient: valkeyClient,
	}
}

// InvalidateByPath invalidates cache entries for a specific path
func (ci *CacheInvalidator) InvalidateByPath(path string) error {
	return ci.ValkeyClient.InvalidatePageCache(fmt.Sprintf("*%s*", path))
}

// InvalidateByPattern invalidates cache entries matching a pattern
func (ci *CacheInvalidator) InvalidateByPattern(pattern string) error {
	return ci.ValkeyClient.InvalidatePageCache(pattern)
}

// InvalidateAll clears all page cache
func (ci *CacheInvalidator) InvalidateAll() error {
	return ci.ValkeyClient.InvalidateAllPageCache()
}

// CacheControlMiddleware adds cache control headers to responses
func CacheControlMiddleware(maxAge int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set cache control headers
			if maxAge > 0 {
				w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
				w.Header().Set("Expires", time.Now().Add(time.Duration(maxAge)*time.Second).Format(http.TimeFormat))
			} else {
				w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
			}

			next.ServeHTTP(w, r)
		})
	}
}
