package router

import (
	"go-service/internal/config"
	"go-service/internal/handlers"
	"go-service/internal/middleware"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// SetupRouter configures and returns the main router with all routes and middleware
func SetupRouter(cfg *config.Config) *mux.Router {
	// Create main router
	r := mux.NewRouter()

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(cfg)
	studentHandler := handlers.NewStudentHandler(cfg)

	// Add custom middleware
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.RecoveryMiddleware)

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins: cfg.CORS.AllowedOrigins,
		AllowedMethods: cfg.CORS.AllowedMethods,
		AllowedHeaders: cfg.CORS.AllowedHeaders,
		AllowCredentials: true,
	})

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Health check endpoint
	api.HandleFunc("/health", healthHandler.HealthCheck).Methods("GET")

	// Student endpoints
	api.HandleFunc("/student/{id}", studentHandler.GetStudentData).Methods("GET")
	api.HandleFunc("/student/{id}/pdf-report", studentHandler.GenerateStudentPDFReport).Methods("GET")
	api.HandleFunc("/student/{id}/save-report", studentHandler.SaveStudentPDFReport).Methods("POST")

	// Apply CORS middleware and return the handler
	return r
} 