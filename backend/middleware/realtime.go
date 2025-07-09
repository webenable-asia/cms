package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// AdminRealtimeHeaders adds headers for real-time admin updates
func AdminRealtimeHeaders() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if this is an admin request
			if isAdminRoute(r) {
				// Disable all caching with comprehensive headers
				w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0, s-maxage=0")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
				w.Header().Set("Last-Modified", "Thu, 01 Jan 1970 00:00:00 GMT")
				
				// Remove any existing ETag to prevent conditional requests
				w.Header().Set("ETag", "")
				
				// Enable real-time features
				w.Header().Set("X-Admin-Realtime", "enabled")
				w.Header().Set("X-Content-Type-Options", "nosniff")
				w.Header().Set("X-Frame-Options", "DENY")
				
				// Add timestamp for freshness verification
				w.Header().Set("X-Server-Time", time.Now().UTC().Format(time.RFC3339))
				
				// Force fresh content with unique ETag
				w.Header().Set("ETag", fmt.Sprintf(`"admin-%d"`, time.Now().UnixNano()))
				
				// Prevent proxy caching
				w.Header().Set("Surrogate-Control", "no-store")
				w.Header().Set("Vary", "*")
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

// RealtimeMiddleware combines real-time headers with cache invalidation
func RealtimeMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Apply real-time headers first
			AdminRealtimeHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})).ServeHTTP(w, r)
		})
	}
}

// isAdminRoute checks if the request is for an admin route
func isAdminRoute(r *http.Request) bool {
	path := r.URL.Path
	referer := r.Header.Get("Referer")
	
	// Direct admin routes
	if strings.HasPrefix(path, "/admin/") || 
	   strings.HasPrefix(path, "/api/admin/") ||
	   strings.HasPrefix(path, "/api/users/") ||
	   strings.HasPrefix(path, "/api/contacts/") {
		return true
	}
	
	// Requests coming from admin panel
	if strings.Contains(referer, "/admin") {
		return true
	}
	
	// Admin query parameter
	if r.URL.Query().Get("admin") != "" {
		return true
	}
	
	return false
}

// NoCache middleware for complete cache bypass
func NoCache() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Ultra-strict no-cache headers
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0, s-maxage=0, proxy-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "Mon, 01 Jan 1990 00:00:00 GMT")
			w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
			w.Header().Set("ETag", fmt.Sprintf(`"nocache-%d"`, time.Now().UnixNano()))
			w.Header().Set("Vary", "*")
			
			next.ServeHTTP(w, r)
		})
	}
}

// AdminSecurityHeaders adds security headers specifically for admin routes
func AdminSecurityHeaders() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isAdminRoute(r) {
				// Enhanced security headers for admin
				w.Header().Set("X-Frame-Options", "DENY")
				w.Header().Set("X-Content-Type-Options", "nosniff")
				w.Header().Set("X-XSS-Protection", "1; mode=block")
				w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
				w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'")
				w.Header().Set("X-Admin-Route", "true")
			}
			
			next.ServeHTTP(w, r)
		})
	}
}
