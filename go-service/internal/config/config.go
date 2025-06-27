package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Config holds all configuration for the application
type Config struct {
	Server  ServerConfig
	NodeJS  NodeJSConfig
	PDF     PDFConfig
	Logging LoggingConfig
	CORS    CORSConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
	Host string
}

// NodeJSConfig holds Node.js API configuration
type NodeJSConfig struct {
	BaseURL   string
	Timeout   time.Duration
	AuthToken string
}

// PDFConfig holds PDF generation configuration
type PDFConfig struct {
	OutputDir string
	Title     string
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string
	Format string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

var AppConfig *Config

// LoadConfig loads configuration from environment variables and config files
func LoadConfig() (*Config, error) {
	// Load from config.env file if it exists
	if _, err := os.Stat("config.env"); err == nil {
		if err := godotenv.Load("config.env"); err != nil {
			logrus.Warn("Error loading config.env file:", err)
		}
	}

	// Load from .env file if it exists
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			logrus.Warn("Error loading .env file:", err)
		}
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnvWithDefault("PORT", "8080"),
			Host: getEnvWithDefault("HOST", "localhost"),
		},
		NodeJS: NodeJSConfig{
			BaseURL:   getEnvWithDefault("NODEJS_API_URL", "http://localhost:3000"),
			Timeout:   time.Duration(getEnvAsInt("NODEJS_API_TIMEOUT", 30)) * time.Second,
			AuthToken: getEnvWithDefault("AUTH_TOKEN", ""),
		},
		PDF: PDFConfig{
			OutputDir: getEnvWithDefault("PDF_OUTPUT_DIR", "./reports"),
			Title:     getEnvWithDefault("PDF_TITLE", "Student Report"),
		},
		Logging: LoggingConfig{
			Level:  getEnvWithDefault("LOG_LEVEL", "info"),
			Format: getEnvWithDefault("LOG_FORMAT", "json"),
		},
		CORS: CORSConfig{
			AllowedOrigins: strings.Split(getEnvWithDefault("CORS_ALLOWED_ORIGINS", "*"), ","),
			AllowedMethods: strings.Split(getEnvWithDefault("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"), ","),
			AllowedHeaders: strings.Split(getEnvWithDefault("CORS_ALLOWED_HEADERS", "*"), ","),
		},
	}

	// Set global config
	AppConfig = config

	// Configure logging
	configureLogging(config.Logging)

	return config, nil
}

// getEnvWithDefault gets an environment variable with a default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as integer with default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// configureLogging configures the logging system
func configureLogging(config LoggingConfig) {
	// Set log level
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		logrus.Warn("Invalid log level, using info")
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// Set log format
	if config.Format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}
}

// GetConfig returns the global configuration
func GetConfig() *Config {
	if AppConfig == nil {
		panic("Configuration not loaded. Call LoadConfig() first.")
	}
	return AppConfig
}
