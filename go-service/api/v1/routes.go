package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"go-service/internal/config"
	"go-service/internal/service"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// PDFHandler handles PDF-related requests
type PDFHandler struct {
	pdfService *service.PDFService
	config     *config.Config
}

// NewPDFHandler creates a new PDF handler
func NewPDFHandler(cfg *config.Config) *PDFHandler {
	return &PDFHandler{
		pdfService: service.NewPDFService(cfg),
		config:     cfg,
	}
}



// GenerateStudentReport generates a PDF report for a student
func (h *PDFHandler) GenerateStudentReport(w http.ResponseWriter, r *http.Request) {
	// Extract student ID from URL parameters
	vars := mux.Vars(r)
	studentIDStr, exists := vars["id"]
	if !exists {
		http.Error(w, "Student ID is required", http.StatusBadRequest)
		return
	}

	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		http.Error(w, "Invalid student ID format", http.StatusBadRequest)
		return
	}

	logrus.Infof("Processing PDF report request for student ID: %d", studentID)

	// Generate the PDF report
	filePath, err := h.pdfService.GenerateStudentReport(studentID)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to generate PDF report for student %d", studentID)
		
		// Check if it's a student not found error (API returned 404)
		if contains(err.Error(), "status 404") || contains(err.Error(), "not found") {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}
		
		http.Error(w, "Failed to generate PDF report", http.StatusInternalServerError)
		return
	}

	// Check if download query parameter is present
	download := r.URL.Query().Get("download")
	
	if download == "true" {
		// Serve the file for download
		h.serveFileDownload(w, r, filePath)
	} else {
		// Return JSON response with file information
		h.returnFileInfo(w, filePath, studentID)
	}
}

// serveFileDownload serves the PDF file for download
func (h *PDFHandler) serveFileDownload(w http.ResponseWriter, r *http.Request, filePath string) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		logrus.WithError(err).Error("Failed to open PDF file")
		http.Error(w, "Failed to open PDF file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		logrus.WithError(err).Error("Failed to get file info")
		http.Error(w, "Failed to get file info", http.StatusInternalServerError)
		return
	}

	// Set headers for file download
	filename := filepath.Base(filePath)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// Copy file content to response
	http.ServeFile(w, r, filePath)
	
	logrus.Infof("PDF file served for download: %s", filename)
}

// returnFileInfo returns JSON information about the generated file
func (h *PDFHandler) returnFileInfo(w http.ResponseWriter, filePath string, studentID int) {
	// Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		logrus.WithError(err).Error("Failed to get file info")
		http.Error(w, "Failed to get file info", http.StatusInternalServerError)
		return
	}

	// Create response
	response := map[string]interface{}{
		"success":     true,
		"message":     "PDF report generated successfully",
		"student_id":  studentID,
		"file_path":   filePath,
		"file_name":   filepath.Base(filePath),
		"file_size":   fileInfo.Size(),
		"generated_at": fileInfo.ModTime().Format("2006-01-02 15:04:05"),
		"download_url": fmt.Sprintf("/api/v1/students/%d/report?download=true", studentID),
	}

	// Set content type and return JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.WithError(err).Error("Failed to encode JSON response")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logrus.Infof("PDF report info returned for student %d", studentID)
}

// HealthCheck provides a simple health check endpoint
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"service":   "go-pdf-service",
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		"version":   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// contains checks if a string contains a substring (case-insensitive helper)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    len(s) > len(substr) && 
		    (s[:len(substr)] == substr || 
		     s[len(s)-len(substr):] == substr || 
		     containsMiddle(s, substr)))
}

// containsMiddle is a helper function for substring checking
func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
