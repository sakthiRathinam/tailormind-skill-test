package main

import (
	"context"
	"fmt"
	"go-service/internal/config"
	"go-service/internal/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load configuration")
	}

	logrus.WithFields(logrus.Fields{
		"port":           cfg.Server.Port,
		"host":           cfg.Server.Host,
		"nodejs_api_url": cfg.NodeJS.BaseURL,
		"log_level":      cfg.Logging.Level,
	}).Info("Starting Go PDF Service")

	// Setup router with all routes and middleware
	r := router.SetupRouter(cfg)

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   cfg.CORS.AllowedMethods,
		AllowedHeaders:   cfg.CORS.AllowedHeaders,
		AllowCredentials: true,
	})

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      c.Handler(r),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logrus.WithField("address", server.Addr).Info("Starting HTTP server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Print available endpoints
	printEndpoints(cfg)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := server.Shutdown(ctx); err != nil {
		logrus.WithError(err).Error("Server forced to shutdown")
	} else {
		logrus.Info("Server exited gracefully")
	}
}

// printEndpoints prints all available API endpoints
func printEndpoints(cfg *config.Config) {
	baseURL := fmt.Sprintf("http://%s:%s", cfg.Server.Host, cfg.Server.Port)
	logrus.Info("--------------------------------")
	logrus.Info("Available endpoints:")

	logrus.Infof("  • Health Check:      GET  %s/api/v1/health", baseURL)
	logrus.Infof("  • Student Report:    GET  %s/api/v1/students/{id}/report", baseURL)
	logrus.Infof("  • Download PDF:      GET  %s/api/v1/students/{id}/report?download=true", baseURL)
	logrus.Info("")
	logrus.Info("Example usage:")
	logrus.Infof("  curl %s/api/v1/health", baseURL)
	logrus.Infof("  curl %s/api/v1/students/1/report", baseURL)
	logrus.Infof("  curl -o student_report.pdf \"%s/api/v1/students/1/report?download=true\"", baseURL)
	logrus.Info("")
	logrus.Info("Note: The service fetches student data from Node.js API configured at:")
	logrus.Infof("  %s/api/v1/students/{id}", cfg.NodeJS.BaseURL)
}
