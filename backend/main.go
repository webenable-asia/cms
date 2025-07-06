// Package main provides the entry point for the WebEnable CMS API
//
//	@title			WebEnable CMS API
//	@version		1.0
//	@description	A Content Management System API with JWT authentication
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	API Support
//	@contact.url	http://www.webenable.com/support
//	@contact.email	support@webenable.com
//
//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT
//
//	@host		localhost:8080
//	@BasePath	/api
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.
package main

import (
	"encoding/json"
	"net/http"
	"time"

	"webenable-cms-backend/cache"
	"webenable-cms-backend/config"
	"webenable-cms-backend/database"
	_ "webenable-cms-backend/docs"
	"webenable-cms-backend/handlers"
	"webenable-cms-backend/middleware"
	"webenable-cms-backend/utils"

	_ "github.com/go-kivik/kivik/v4/couchdb"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Initialize configuration
	config.Init()

	// Initialize database
	database.Init()

	// Initialize Valkey cache
	valkeyClient, err := cache.NewValkeyClient(config.AppConfig.ValkeyURL)
	if err != nil {
		utils.LogError(err, "Failed to connect to Valkey", logrus.Fields{
			"valkey_url": config.AppConfig.ValkeyURL,
		})
		panic(err)
	}
	defer valkeyClient.Close()

	// Set global cache for handlers
	handlers.SetGlobalCache(valkeyClient)

	// Initialize middleware
	rateLimiter := middleware.NewRateLimiter(valkeyClient)
	pageCache := middleware.NewPageCache(valkeyClient).WithTTL(10 * time.Minute)

	// Set global rate limiter for handlers
	handlers.SetGlobalRateLimiter(rateLimiter)

	// Initialize router
	r := mux.NewRouter()

	// Add security headers middleware (first)
	r.Use(middleware.SecurityHeaders)

	// Add XSS protection middleware
	r.Use(middleware.XSSProtection)

	// Add compression middleware
	r.Use(middleware.CompressionMiddleware)

	// Add global rate limiting
	r.Use(rateLimiter.RateLimit(100)) // 100 requests per minute per IP

	// Swagger documentation endpoint
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Health check endpoint
	// HealthCheck godoc
	//	@Summary		Health check
	//	@Description	Check API and cache health status
	//	@Tags			System
	//	@Accept			json
	//	@Produce		json
	//	@Success		200	{object}	object{status=string,cache=string}
	//	@Failure		503	{object}	models.ErrorResponse
	//	@Router			/health [get]
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Check Valkey health
		if err := valkeyClient.Health(); err != nil {
			http.Error(w, "Cache unavailable", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","cache":"connected"}`))
	}).Methods("GET")

	// Public routes with lighter rate limiting and page caching
	public := api.PathPrefix("").Subrouter()
	public.Use(rateLimiter.RateLimit(100))             // 100 requests per minute for public routes
	public.Use(pageCache.PageCacheMiddleware())        // Add page caching for public routes
	public.Use(middleware.CacheControlMiddleware(600)) // 10 minutes browser cache
	public.HandleFunc("/posts", handlers.GetPosts).Methods("GET")
	public.HandleFunc("/posts/{id}", handlers.GetPost).Methods("GET")
	public.HandleFunc("/contact", handlers.SubmitContact).Methods("POST")

	// Authentication routes with strict rate limiting
	auth := api.PathPrefix("/auth").Subrouter()
	auth.Use(rateLimiter.AuthRateLimit(100)) // 100 attempts per hour for auth (development)
	auth.HandleFunc("/login", handlers.Login).Methods("POST")
	auth.HandleFunc("/logout", handlers.Logout).Methods("POST")

	// Protected auth routes (require JWT authentication)
	authProtected := auth.PathPrefix("").Subrouter()
	authProtected.Use(middleware.AuthMiddleware)
	authProtected.HandleFunc("/me", handlers.GetCurrentUser).Methods("GET")

	// Protected routes (require JWT authentication)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.Use(rateLimiter.UserRateLimit(150)) // 150 requests per minute for authenticated users
	protected.HandleFunc("/posts", handlers.CreatePost).Methods("POST")
	protected.HandleFunc("/posts/{id}", handlers.UpdatePost).Methods("PUT")
	protected.HandleFunc("/posts/{id}", handlers.DeletePost).Methods("DELETE")
	protected.HandleFunc("/contacts", handlers.GetContacts).Methods("GET")
	protected.HandleFunc("/contacts/{id}", handlers.GetContact).Methods("GET")
	protected.HandleFunc("/contacts/{id}", handlers.UpdateContactStatus).Methods("PUT")
	protected.HandleFunc("/contacts/{id}/reply", handlers.ReplyToContact).Methods("POST")
	protected.HandleFunc("/contacts/{id}", handlers.DeleteContact).Methods("DELETE")

	// User management routes (admin only)
	protected.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	protected.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	protected.HandleFunc("/users/stats", handlers.GetUserStats).Methods("GET")
	protected.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
	protected.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")
	protected.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")

	// Admin routes for rate limit management (admin only)
	protected.HandleFunc("/admin/rate-limit/reset", handlers.ResetRateLimit).Methods("POST")
	protected.HandleFunc("/admin/rate-limit/status", handlers.GetRateLimitStatus).Methods("GET")

	// Stats endpoint for monitoring
	protected.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Get cache stats
		cacheStats, err := valkeyClient.GetCacheStats()
		if err != nil {
			http.Error(w, "Failed to get cache stats", http.StatusInternalServerError)
			return
		}

		stats := map[string]interface{}{
			"cache":     cacheStats,
			"auth":      "JWT-based",
			"timestamp": time.Now(),
		}

		json.NewEncoder(w).Encode(stats)
	}).Methods("GET")

	// Setup CORS with secure configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   config.AppConfig.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	utils.LogInfo("Server starting", logrus.Fields{
		"port":            config.AppConfig.Port,
		"valkey_url":      config.AppConfig.ValkeyURL,
		"auth_mode":       "JWT-based (cookieless)",
		"allowed_origins": config.AppConfig.AllowedOrigins,
	})

	if err := http.ListenAndServe(":"+config.AppConfig.Port, handler); err != nil {
		utils.LogError(err, "Server failed to start", logrus.Fields{
			"port": config.AppConfig.Port,
		})
		panic(err)
	}
}
