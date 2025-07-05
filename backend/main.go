package main

import (
	"log"
	"net/http"

	"webenable-cms-backend/cache"
	"webenable-cms-backend/config"
	"webenable-cms-backend/database"
	"webenable-cms-backend/handlers"
	"webenable-cms-backend/middleware"

	_ "github.com/go-kivik/kivik/v4/couchdb"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Initialize configuration
	config.Init()

	// Initialize database
	database.Init()

	// Initialize Valkey cache
	valkeyClient, err := cache.NewValkeyClient(config.AppConfig.ValkeyURL)
	if err != nil {
		log.Fatalf("Failed to connect to Valkey: %v", err)
	}
	defer valkeyClient.Close()

	// Initialize middleware
	sessionManager := middleware.NewSessionManager(valkeyClient, config.AppConfig.SessionDomain, config.AppConfig.SessionSecure)
	rateLimiter := middleware.NewRateLimiter(valkeyClient)

	// Initialize router
	r := mux.NewRouter()

	// Add global rate limiting
	r.Use(rateLimiter.RateLimit(100)) // 100 requests per minute per IP

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Health check endpoint
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Check Valkey health
		if err := valkeyClient.Health(); err != nil {
			http.Error(w, "Cache unavailable", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","cache":"connected"}`))
	}).Methods("GET")

	// Public routes with lighter rate limiting
	public := api.PathPrefix("").Subrouter()
	public.Use(rateLimiter.RateLimit(60)) // 60 requests per minute for public routes
	public.HandleFunc("/posts", handlers.GetPosts).Methods("GET")
	public.HandleFunc("/posts/{id}", handlers.GetPost).Methods("GET")
	public.HandleFunc("/contact", handlers.SubmitContact).Methods("POST")

	// Authentication routes with strict rate limiting
	auth := api.PathPrefix("/auth").Subrouter()
	auth.Use(rateLimiter.AuthRateLimit(10)) // 10 attempts per hour for auth
	auth.HandleFunc("/login", handlers.Login).Methods("POST")
	auth.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		err := sessionManager.DestroySession(w, r)
		if err != nil {
			http.Error(w, "Error logging out", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Logged out successfully"}`))
	}).Methods("POST")

	// Protected routes (require session authentication)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(sessionManager.SessionMiddleware)
	protected.Use(rateLimiter.UserRateLimit(120)) // 120 requests per minute for authenticated users
	protected.HandleFunc("/posts", handlers.CreatePost).Methods("POST")
	protected.HandleFunc("/posts/{id}", handlers.UpdatePost).Methods("PUT")
	protected.HandleFunc("/posts/{id}", handlers.DeletePost).Methods("DELETE")
	protected.HandleFunc("/contacts", handlers.GetContacts).Methods("GET")
	protected.HandleFunc("/contacts/{id}", handlers.GetContact).Methods("GET")
	protected.HandleFunc("/contacts/{id}", handlers.UpdateContactStatus).Methods("PUT")
	protected.HandleFunc("/contacts/{id}/reply", handlers.ReplyToContact).Methods("POST")
	protected.HandleFunc("/contacts/{id}", handlers.DeleteContact).Methods("DELETE")

	// Setup CORS with secure configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   config.AppConfig.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Printf("Server starting on port %s", config.AppConfig.Port)
	log.Printf("Valkey cache connected at %s", config.AppConfig.ValkeyURL)
	log.Fatal(http.ListenAndServe(":"+config.AppConfig.Port, handler))
}
