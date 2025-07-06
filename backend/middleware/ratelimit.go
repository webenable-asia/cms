package middleware

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"webenable-cms-backend/cache"
)

// RateLimiter handles rate limiting using Valkey
type RateLimiter struct {
	cache *cache.ValkeyClient
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(valkeyClient *cache.ValkeyClient) *RateLimiter {
	return &RateLimiter{
		cache: valkeyClient,
	}
}

// getClientIP extracts the client IP address
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}

// RateLimit middleware for general API endpoints
func (rl *RateLimiter) RateLimit(requestsPerMinute int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)
			identifier := fmt.Sprintf("rate_limit:api:%s", clientIP)

			allowed, err := rl.cache.SetRateLimit(identifier, requestsPerMinute, time.Minute)
			if err != nil {
				// If cache is down, allow the request but log the error
				fmt.Printf("Rate limit check failed: %v\n", err)
				next.ServeHTTP(w, r)
				return
			}

			if !allowed {
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", requestsPerMinute))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("Retry-After", "60")
				http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
				return
			}

			// Get current count for headers
			current, _ := rl.getCurrentCount(identifier)
			remaining := requestsPerMinute - current
			if remaining < 0 {
				remaining = 0
			}

			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", requestsPerMinute))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

			next.ServeHTTP(w, r)
		})
	}
}

// AuthRateLimit middleware specifically for authentication endpoints
func (rl *RateLimiter) AuthRateLimit(attemptsPerHour int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)
			identifier := fmt.Sprintf("rate_limit:auth:%s", clientIP)

			allowed, err := rl.cache.SetRateLimit(identifier, attemptsPerHour, time.Hour)
			if err != nil {
				// If cache is down, allow the request but log the error
				fmt.Printf("Auth rate limit check failed: %v\n", err)
				next.ServeHTTP(w, r)
				return
			}

			if !allowed {
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", attemptsPerHour))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("Retry-After", "3600")
				http.Error(w, "Too many authentication attempts. Please try again later.", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// UserRateLimit middleware for user-specific rate limiting
func (rl *RateLimiter) UserRateLimit(requestsPerMinute int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := GetSessionFromContext(r)
			if session == nil {
				// No session, fall back to IP-based rate limiting
				fallback := rl.RateLimit(requestsPerMinute)
				fallback(next).ServeHTTP(w, r)
				return
			}

			identifier := fmt.Sprintf("rate_limit:user:%s", session.UserID)

			allowed, err := rl.cache.SetRateLimit(identifier, requestsPerMinute, time.Minute)
			if err != nil {
				fmt.Printf("User rate limit check failed: %v\n", err)
				next.ServeHTTP(w, r)
				return
			}

			if !allowed {
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", requestsPerMinute))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("Retry-After", "60")
				http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getCurrentCount gets the current count for rate limiting headers
func (rl *RateLimiter) getCurrentCount(identifier string) (int, error) {
	var count int
	err := rl.cache.Get(identifier, &count)
	if err != nil {
		return 0, nil // Key doesn't exist yet
	}
	return count, nil
}

// ClearRateLimit clears rate limit for a specific identifier (admin function)
func (rl *RateLimiter) ClearRateLimit(identifier string) error {
	return rl.cache.Delete(identifier)
}

// ResetRateLimitForIP resets rate limit for a specific IP address
func (rl *RateLimiter) ResetRateLimitForIP(ipAddress string) error {
	// Reset API rate limit
	apiKey := fmt.Sprintf("rate_limit:api:%s", ipAddress)
	err := rl.cache.Delete(apiKey)
	if err != nil {
		return fmt.Errorf("failed to reset API rate limit for IP %s: %w", ipAddress, err)
	}

	// Reset auth rate limit
	authKey := fmt.Sprintf("rate_limit:auth:%s", ipAddress)
	err = rl.cache.Delete(authKey)
	if err != nil {
		return fmt.Errorf("failed to reset auth rate limit for IP %s: %w", ipAddress, err)
	}

	return nil
}

// ResetRateLimitForUser resets rate limit for a specific user
func (rl *RateLimiter) ResetRateLimitForUser(userID string) error {
	userKey := fmt.Sprintf("rate_limit:user:%s", userID)
	return rl.cache.Delete(userKey)
}

// ResetAllAPIRateLimits resets all API rate limits
func (rl *RateLimiter) ResetAllAPIRateLimits() error {
	return rl.cache.ResetRateLimitByPattern("api:*")
}

// ResetAllAuthRateLimits resets all authentication rate limits
func (rl *RateLimiter) ResetAllAuthRateLimits() error {
	return rl.cache.ResetRateLimitByPattern("auth:*")
}

// ResetAllUserRateLimits resets all user-specific rate limits
func (rl *RateLimiter) ResetAllUserRateLimits() error {
	return rl.cache.ResetRateLimitByPattern("user:*")
}

// ResetAllRateLimits resets all rate limits
func (rl *RateLimiter) ResetAllRateLimits() error {
	return rl.cache.ResetAllRateLimits()
}

// GetRateLimitStatus returns the current rate limit status
func (rl *RateLimiter) GetRateLimitStatus(identifier string, limit int) (current int, remaining int, resetTime time.Time, err error) {
	current, err = rl.getCurrentCount(identifier)
	if err != nil {
		return 0, 0, time.Time{}, err
	}

	remaining = limit - current
	if remaining < 0 {
		remaining = 0
	}

	// Get TTL for reset time
	ttl, err := rl.cache.GetTTL(identifier)
	if err != nil {
		return current, remaining, time.Time{}, err
	}

	resetTime = time.Now().Add(ttl)
	return current, remaining, resetTime, nil
}
