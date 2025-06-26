package v1

import (
	"go-service/internal/config"
	"github.com/gorilla/mux"
)

// RegisterV1Routes registers all v1 API routes to the given router
func RegisterV1Routes(router *mux.Router, cfg *config.Config) {
	// Create PDF handler
	pdfHandler := NewPDFHandler(cfg)
	
	// Create v1 subrouter
	v1Router := router.PathPrefix("/api/v1").Subrouter()
	
	// Register all v1 routes
	v1Router.HandleFunc("/students/{id}/report", pdfHandler.GenerateStudentReport).Methods("GET")
	v1Router.HandleFunc("/health", HealthCheck).Methods("GET")
}
