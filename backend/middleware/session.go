package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"webenable-cms-backend/cache"
)

// SessionData represents session information
type SessionData struct {
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	LastSeen  time.Time `json:"last_seen"`
}

// SessionManager handles session operations
type SessionManager struct {
	cache      *cache.ValkeyClient
	cookieName string
	domain     string
	secure     bool
	httpOnly   bool
	sameSite   http.SameSite
	maxAge     time.Duration
}

// NewSessionManager creates a new session manager
func NewSessionManager(valkeyClient *cache.ValkeyClient, domain string, secure bool) *SessionManager {
	return &SessionManager{
		cache:      valkeyClient,
		cookieName: "webenable_session",
		domain:     domain,
		secure:     secure,
		httpOnly:   true,
		sameSite:   http.SameSiteStrictMode,
		maxAge:     24 * time.Hour, // 24 hours
	}
}

// generateSessionID creates a secure random session ID
func (sm *SessionManager) generateSessionID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// CreateSession creates a new session
func (sm *SessionManager) CreateSession(w http.ResponseWriter, userData SessionData) (string, error) {
	sessionID := sm.generateSessionID()

	// Set session data
	userData.CreatedAt = time.Now()
	userData.LastSeen = time.Now()

	err := sm.cache.SetSession(sessionID, userData, sm.maxAge)
	if err != nil {
		return "", err
	}

	// Set cookie
	cookie := &http.Cookie{
		Name:     sm.cookieName,
		Value:    sessionID,
		Domain:   sm.domain,
		Path:     "/",
		MaxAge:   int(sm.maxAge.Seconds()),
		Secure:   sm.secure,
		HttpOnly: sm.httpOnly,
		SameSite: sm.sameSite,
	}

	http.SetCookie(w, cookie)

	return sessionID, nil
}

// GetSession retrieves session data
func (sm *SessionManager) GetSession(r *http.Request) (*SessionData, error) {
	cookie, err := r.Cookie(sm.cookieName)
	if err != nil {
		return nil, err
	}

	var sessionData SessionData
	err = sm.cache.GetSession(cookie.Value, &sessionData)
	if err != nil {
		return nil, err
	}

	// Update last seen time
	sessionData.LastSeen = time.Now()
	sm.cache.SetSession(cookie.Value, sessionData, sm.maxAge)

	return &sessionData, nil
}

// DestroySession removes a session
func (sm *SessionManager) DestroySession(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(sm.cookieName)
	if err != nil {
		return nil // Cookie doesn't exist, already "destroyed"
	}

	// Remove from cache
	err = sm.cache.DeleteSession(cookie.Value)
	if err != nil {
		return err
	}

	// Clear cookie
	clearCookie := &http.Cookie{
		Name:     sm.cookieName,
		Value:    "",
		Domain:   sm.domain,
		Path:     "/",
		MaxAge:   -1,
		Secure:   sm.secure,
		HttpOnly: sm.httpOnly,
		SameSite: sm.sameSite,
	}

	http.SetCookie(w, clearCookie)

	return nil
}

// SessionMiddleware validates sessions for protected routes
func (sm *SessionManager) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip authentication for public routes
		if isPublicRoute(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		sessionData, err := sm.GetSession(r)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
			return
		}

		// Add session data to request context
		ctx := context.WithValue(r.Context(), "session", sessionData)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// isPublicRoute checks if a route is public
func isPublicRoute(path string) bool {
	publicRoutes := []string{
		"/api/auth/login",
		"/api/posts",
		"/api/health",
		"/",
	}

	for _, route := range publicRoutes {
		if path == route || strings.HasPrefix(path, route+"/") {
			return true
		}
	}

	return false
}

// GetSessionFromContext retrieves session data from request context
func GetSessionFromContext(r *http.Request) *SessionData {
	session, ok := r.Context().Value("session").(*SessionData)
	if !ok {
		return nil
	}
	return session
}

// RefreshSession extends session expiration
func (sm *SessionManager) RefreshSession(r *http.Request) error {
	cookie, err := r.Cookie(sm.cookieName)
	if err != nil {
		return err
	}

	var sessionData SessionData
	err = sm.cache.GetSession(cookie.Value, &sessionData)
	if err != nil {
		return err
	}

	// Extend session
	sessionData.LastSeen = time.Now()
	return sm.cache.SetSession(cookie.Value, sessionData, sm.maxAge)
}
