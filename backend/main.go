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

	"webenable-cms-backend/adapters"
	"webenable-cms-backend/cache"
	"webenable-cms-backend/config"
	"webenable-cms-backend/container"
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

	// Initialize legacy database (for backward compatibility during migration)
	database.Init()

	// Create adapter factory
	factory := adapters.NewAdapterFactory(config.AppConfig.Adapters)

	// Create all adapters
	adapterSet, err := factory.CreateAllAdapters()
	if err != nil {
		utils.LogError(err, "Failed to create adapters", logrus.Fields{})
		panic(err)
	}
	defer adapterSet.Close()

	// Check adapter health
	if err := adapterSet.Health(); err != nil {
		utils.LogError(err, "Adapter health check failed", logrus.Fields{})
		panic(err)
	}

	// Create service container
	serviceContainer, err := container.NewContainer(config.AppConfig.Adapters)
	if err != nil {
		utils.LogError(err, "Failed to create service container", logrus.Fields{})
		panic(err)
	}
	defer serviceContainer.Close()

	// Initialize legacy Valkey client for backward compatibility
	valkeyClient, err := cache.NewValkeyClient(config.AppConfig.ValkeyURL)
	if err != nil {
		utils.LogError(err, "Failed to connect to Valkey", logrus.Fields{
			"valkey_url": config.AppConfig.ValkeyURL,
		})
		panic(err)
	}
	defer valkeyClient.Close()

	// Set global dependencies for handlers (legacy support)
	handlers.SetGlobalCache(valkeyClient)

	// Initialize middleware using adapters
	rateLimiter := middleware.NewRateLimiter(valkeyClient)
	pageCache := middleware.NewPageCache(valkeyClient).WithTTL(10 * time.Minute)

	// Set global rate limiter for handlers
	handlers.SetGlobalRateLimiter(rateLimiter)

	// Set service container for handlers
	handlers.SetServiceContainer(serviceContainer)

	// Set service container for middleware
	middleware.SetServiceContainer(serviceContainer)

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
	//	@Description	Check API and all adapter health status
	//	@Tags			System
	//	@Accept			json
	//	@Produce		json
	//	@Success		200	{object}	object{status=string,adapters=object}
	//	@Failure		503	{object}	models.ErrorResponse
	//	@Router			/health [get]
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Check all adapter health
		if err := adapterSet.Health(); err != nil {
			http.Error(w, "One or more adapters unavailable", http.StatusServiceUnavailable)
			return
		}

		// Create detailed health response
		response := map[string]interface{}{
			"status": "healthy",
			"adapters": map[string]string{
				"database": "connected",
				"cache":    "connected",
				"auth":     "connected",
				"email":    "connected",
				"storage":  "connected",
			},
			"timestamp": time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}).Methods("GET")

	// Real-time routes (no caching) for admin features
	realtime := api.PathPrefix("").Subrouter()
	realtime.Use(rateLimiter.RateLimit(100))           // 100 requests per minute
	realtime.Use(middleware.NoCache())                 // No caching for real-time data
	realtime.HandleFunc("/posts", handlers.GetPosts).Methods("GET")
	
	// Public routes with lighter rate limiting and page caching
	public := api.PathPrefix("").Subrouter()
	public.Use(rateLimiter.RateLimit(100))             // 100 requests per minute for public routes
	public.Use(pageCache.PageCacheMiddleware())        // Add page caching for public routes
	public.Use(middleware.CacheControlMiddleware(600)) // 10 minutes browser cache
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

	// Admin routes with real-time headers and no caching
	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.AuthMiddleware)
	admin.Use(middleware.AdminRealtimeHeaders())       // Real-time headers for admin
	admin.Use(middleware.NoCache())                    // Complete cache bypass
	admin.Use(middleware.AdminSecurityHeaders())       // Enhanced security for admin
	admin.Use(rateLimiter.UserRateLimit(200))         // Higher rate limit for admin operations

	// Admin-specific API routes
	admin.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	admin.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	admin.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
	admin.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")
	admin.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")
	admin.HandleFunc("/contacts", handlers.GetContacts).Methods("GET")
	admin.HandleFunc("/contacts/{id}", handlers.GetContact).Methods("GET")
	admin.HandleFunc("/contacts/{id}", handlers.UpdateContactStatus).Methods("PUT")
	admin.HandleFunc("/contacts/{id}/reply", handlers.ReplyToContact).Methods("POST")
	admin.HandleFunc("/contacts/{id}", handlers.DeleteContact).Methods("DELETE")

	// Legacy protected routes for backward compatibility
	protected.HandleFunc("/contacts", handlers.GetContacts).Methods("GET")
	protected.HandleFunc("/contacts/{id}", handlers.GetContact).Methods("GET")
	protected.HandleFunc("/contacts/{id}", handlers.UpdateContactStatus).Methods("PUT")
	protected.HandleFunc("/contacts/{id}/reply", handlers.ReplyToContact).Methods("POST")
	protected.HandleFunc("/contacts/{id}", handlers.DeleteContact).Methods("DELETE")

	// User management routes (admin only) - keep for backward compatibility
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

		// Get cache stats using adapter
		cacheAdapter := serviceContainer.Cache()
		cacheStats, err := cacheAdapter.GetStats()
		if err != nil {
			http.Error(w, "Failed to get cache stats", http.StatusInternalServerError)
			return
		}

		// Get adapter health status
		adapterHealth := make(map[string]string)
		if err := serviceContainer.Database().Health(); err != nil {
			adapterHealth["database"] = "unhealthy"
		} else {
			adapterHealth["database"] = "healthy"
		}

		if err := serviceContainer.Cache().Health(); err != nil {
			adapterHealth["cache"] = "unhealthy"
		} else {
			adapterHealth["cache"] = "healthy"
		}

		if err := serviceContainer.Auth().Health(); err != nil {
			adapterHealth["auth"] = "unhealthy"
		} else {
			adapterHealth["auth"] = "healthy"
		}

		if err := serviceContainer.Email().Health(); err != nil {
			adapterHealth["email"] = "unhealthy"
		} else {
			adapterHealth["email"] = "healthy"
		}

		if err := serviceContainer.Storage().Health(); err != nil {
			adapterHealth["storage"] = "unhealthy"
		} else {
			adapterHealth["storage"] = "healthy"
		}

		stats := map[string]interface{}{
			"cache":          cacheStats,
			"auth":           "JWT-based",
			"adapter_health": adapterHealth,
			"adapter_types": map[string]string{
				"database": config.AppConfig.Adapters.Database.Type,
				"cache":    config.AppConfig.Adapters.Cache.Type,
				"auth":     config.AppConfig.Adapters.Auth.Type,
				"email":    config.AppConfig.Adapters.Email.Type,
				"storage":  config.AppConfig.Adapters.Storage.Type,
			},
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
