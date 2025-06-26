package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go-service/internal/config"
	"go-service/internal/models"

	"github.com/jung-kurt/gofpdf"
	"github.com/sirupsen/logrus"
)

func init() {
	// Set log level to error to reduce noise during tests
	logrus.SetLevel(logrus.ErrorLevel)
}

// TestPDFService_FetchStudentData tests the student data fetching functionality
func TestPDFService_FetchStudentData(t *testing.T) {
	// Create mock student data using the new schema
	mockStudent := &models.Student{
		ID:                 1,
		Name:               "John Doe",
		Email:              "john.doe@example.com",
		Phone:              "+1234567890",
		Gender:             "Male",
		DOB:                "2005-01-15",
		Class:              "10th Grade",
		Section:            "A",
		Roll:               123,
		CurrentAddress:     "123 Main St, City, State",
		PermanentAddress:   "123 Main St, City, State",
		AdmissionDate:      "2023-01-15",
		SystemAccess:       true,
		FatherName:         "John Sr",
		FatherPhone:        "+1234567891",
		MotherName:         "Jane Doe",
		MotherPhone:        "+1234567892",
		GuardianName:       "John Sr",
		GuardianPhone:      "+1234567891",
		RelationOfGuardian: "Father",
		ReporterName:       "Admin User",
	}

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/students/1" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockStudent)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create test config
	cfg := &config.Config{
		NodeJS: config.NodeJSConfig{
			BaseURL: server.URL,
			Timeout: 30 * time.Second,
		},
	}

	// Create PDF service
	service := NewPDFService(cfg)

	// Test successful fetch
	t.Run("SuccessfulFetch", func(t *testing.T) {
		student, err := service.FetchStudentData(1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if student == nil {
			t.Fatal("Expected student data, got nil")
		}

		if student.Name != "John Doe" {
			t.Errorf("Expected John Doe, got %s", student.Name)
		}

		if student.Email != "john.doe@example.com" {
			t.Errorf("Expected john.doe@example.com, got %s", student.Email)
		}

		if student.Roll != 123 {
			t.Errorf("Expected roll 123, got %d", student.Roll)
		}

		if student.Class != "10th Grade" {
			t.Errorf("Expected 10th Grade, got %s", student.Class)
		}
	})

	// Test not found
	t.Run("StudentNotFound", func(t *testing.T) {
		_, err := service.FetchStudentData(999)
		if err == nil {
			t.Fatal("Expected error for non-existent student, got nil")
		}
	})
}

// TestPDFService_GeneratePDFReport tests the PDF generation functionality
func TestPDFService_GeneratePDFReport(t *testing.T) {
	// Create test config
	tempDir := t.TempDir()
	cfg := &config.Config{
		PDF: config.PDFConfig{
			OutputDir: tempDir,
			Title:     "Test Student Report",
		},
	}

	// Create test student data using the new schema
	student := &models.Student{
		ID:                 1,
		Name:               "Jane Smith",
		Email:              "jane.smith@example.com",
		Phone:              "+1987654321",
		Gender:             "Female",
		DOB:                "2006-03-20",
		Class:              "11th Grade",
		Section:            "B",
		Roll:               456,
		CurrentAddress:     "456 Oak Ave, Town, State",
		PermanentAddress:   "789 Pine St, Village, State",
		AdmissionDate:      "2022-08-20",
		SystemAccess:       false,
		FatherName:         "Robert Smith",
		FatherPhone:        "+1987654322",
		MotherName:         "Mary Smith",
		MotherPhone:        "+1987654323",
		GuardianName:       "Robert Smith",
		GuardianPhone:      "+1987654322",
		RelationOfGuardian: "Father",
		ReporterName:       "School Admin",
	}

	// Create PDF service
	service := NewPDFService(cfg)

	// Test PDF generation
	t.Run("GeneratePDF", func(t *testing.T) {
		filePath, err := service.GeneratePDFReport(student)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Fatalf("Expected PDF file to exist at %s", filePath)
		}

		// Check file size (should be > 0)
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			t.Fatalf("Failed to get file info: %v", err)
		}

		if fileInfo.Size() == 0 {
			t.Fatal("Expected PDF file to have content, but it's empty")
		}

		t.Logf("Generated PDF file: %s (size: %d bytes)", filePath, fileInfo.Size())
	})

	// Test PDF with minimal data
	t.Run("GeneratePDFMinimalData", func(t *testing.T) {
		minimalStudent := &models.Student{
			ID:             2,
			Name:           "Minimal Student",
			Email:          "minimal@example.com",
			Phone:          "+1111111111",
			Class:          "9th Grade",
			Section:        "C",
			Roll:           1,
			CurrentAddress: "Address",
			AdmissionDate:  "2023-01-01",
			SystemAccess:   false,
		}

		filePath, err := service.GeneratePDFReport(minimalStudent)
		if err != nil {
			t.Fatalf("Expected no error with minimal data, got %v", err)
		}

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Fatalf("Expected PDF file to exist at %s", filePath)
		}

		t.Logf("Generated minimal PDF file: %s", filePath)
	})
}

// TestPDFComparison tests PDF file comparison functionality
func TestPDFComparison(t *testing.T) {
	tempDir := t.TempDir()

	// Create two identical PDFs
	pdf1Path := createTestPDF(t, tempDir, "test1.pdf", "Test Student 1")
	pdf2Path := createTestPDF(t, tempDir, "test2.pdf", "Test Student 1") // Same content
	pdf3Path := createTestPDF(t, tempDir, "test3.pdf", "Test Student 2") // Different content

	t.Run("IdenticalPDFs", func(t *testing.T) {
		equal, err := ComparePDFFiles(pdf1Path, pdf2Path)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !equal {
			t.Error("Expected PDFs to be identical")
		}
	})

	t.Run("DifferentPDFs", func(t *testing.T) {
		equal, err := ComparePDFFiles(pdf1Path, pdf3Path)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if equal {
			t.Error("Expected PDFs to be different")
		}
	})

	t.Run("NonExistentFile", func(t *testing.T) {
		_, err := ComparePDFFiles(pdf1Path, "nonexistent.pdf")
		if err == nil {
			t.Fatal("Expected error for non-existent file")
		}
	})
}

// TestPDFQuality tests the quality of generated PDFs
func TestPDFQuality(t *testing.T) {
	tempDir := t.TempDir()
	cfg := &config.Config{
		PDF: config.PDFConfig{
			OutputDir: tempDir,
			Title:     "Quality Test Report",
		},
	}

	// Create comprehensive student data using the new schema
	student := &models.Student{
		ID:                 123,
		Name:               "Quality Test Student",
		Email:              "quality.test@example.com",
		Phone:              "+1555123456",
		Gender:             "Non-binary",
		DOB:                "2004-07-15",
		Class:              "12th Grade Advanced",
		Section:            "C",
		Roll:               789,
		CurrentAddress:     "789 Quality Lane, Test City, QT 12345",
		PermanentAddress:   "456 Testing Blvd, QA Town, QT 54321",
		AdmissionDate:      "2021-09-01",
		SystemAccess:       true,
		FatherName:         "Quality Father",
		FatherPhone:        "+1555123457",
		MotherName:         "Quality Mother",
		MotherPhone:        "+1555123458",
		GuardianName:       "Quality Guardian",
		GuardianPhone:      "+1555123459",
		RelationOfGuardian: "Uncle",
		ReporterName:       "Quality Assurance Admin",
	}

	service := NewPDFService(cfg)

	t.Run("PDFQualityMetrics", func(t *testing.T) {
		filePath, err := service.GeneratePDFReport(student)
		if err != nil {
			t.Fatalf("Failed to generate PDF: %v", err)
		}

		// Test file size (should be reasonable)
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			t.Fatalf("Failed to get file info: %v", err)
		}

		size := fileInfo.Size()
		if size < 1000 { // Less than 1KB might indicate missing content
			t.Errorf("PDF file too small (%d bytes), might be missing content", size)
		}
		if size > 1000000 { // More than 1MB might indicate inefficiency
			t.Errorf("PDF file too large (%d bytes), might be inefficient", size)
		}

		t.Logf("PDF file size: %d bytes (acceptable range)", size)

		// Test PDF structure by attempting to read it
		err = ValidatePDFStructure(filePath)
		if err != nil {
			t.Errorf("PDF structure validation failed: %v", err)
		}
	})

	t.Run("PDFContentValidation", func(t *testing.T) {
		filePath, err := service.GeneratePDFReport(student)
		if err != nil {
			t.Fatalf("Failed to generate PDF: %v", err)
		}

		// Validate that the file was created with expected naming pattern
		expectedPattern := fmt.Sprintf("student_%d_report_", student.ID)
		filename := filepath.Base(filePath)
		if !strings.Contains(filename, expectedPattern) {
			t.Errorf("Expected filename to contain %s, got %s", expectedPattern, filename)
		}

		// Validate file extension
		if filepath.Ext(filename) != ".pdf" {
			t.Errorf("Expected .pdf extension, got %s", filepath.Ext(filename))
		}
	})
}

// TestPDFService_GenerateStudentReport tests the complete workflow
func TestPDFService_GenerateStudentReport(t *testing.T) {
	// Create test student data
	mockStudent := &models.Student{
		ID:               42,
		Name:             "Integration Test Student",
		Email:            "integration@test.com",
		Phone:            "+1999888777",
		Gender:           "Female",
		DOB:              "2005-12-25",
		Class:            "12th Grade",
		Section:          "A",
		Roll:             999,
		CurrentAddress:   "123 Integration St",
		PermanentAddress: "456 Test Ave",
		AdmissionDate:    "2020-01-01",
		SystemAccess:     true,
		FatherName:       "Test Father",
		FatherPhone:      "+1999888778",
		MotherName:       "Test Mother",
		MotherPhone:      "+1999888779",
		ReporterName:     "Integration Admin",
	}

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/students/42" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockStudent)
		} else {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Student not found"})
		}
	}))
	defer server.Close()

	// Create test config
	tempDir := t.TempDir()
	cfg := &config.Config{
		NodeJS: config.NodeJSConfig{
			BaseURL: server.URL,
			Timeout: 30 * time.Second,
		},
		PDF: config.PDFConfig{
			OutputDir: tempDir,
			Title:     "Integration Test Report",
		},
	}

	// Create PDF service
	service := NewPDFService(cfg)

	t.Run("CompleteWorkflow", func(t *testing.T) {
		filePath, err := service.GenerateStudentReport(42)
		if err != nil {
			t.Fatalf("Expected no error in complete workflow, got %v", err)
		}

		// Validate the generated file
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Fatalf("Expected PDF file to exist at %s", filePath)
		}

		// Validate PDF structure
		err = ValidatePDFStructure(filePath)
		if err != nil {
			t.Errorf("PDF structure validation failed: %v", err)
		}

		t.Logf("Complete workflow successful: %s", filePath)
	})

	t.Run("StudentNotFound", func(t *testing.T) {
		_, err := service.GenerateStudentReport(404)
		if err == nil {
			t.Fatal("Expected error for non-existent student, got nil")
		}

		t.Logf("Expected error received: %v", err)
	})
}

// createTestPDF creates a test PDF file for comparison testing
func createTestPDF(t *testing.T, dir, filename, content string) string {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, content)

	filepath := filepath.Join(dir, filename)
	err := pdf.OutputFileAndClose(filepath)
	if err != nil {
		t.Fatalf("Failed to create test PDF: %v", err)
	}

	return filepath
}

// ComparePDFFiles compares two PDF files for equality
func ComparePDFFiles(file1, file2 string) (bool, error) {
	// Read first file
	data1, err := os.ReadFile(file1)
	if err != nil {
		return false, fmt.Errorf("failed to read first file: %w", err)
	}

	// Read second file
	data2, err := os.ReadFile(file2)
	if err != nil {
		return false, fmt.Errorf("failed to read second file: %w", err)
	}

	// Compare byte-by-byte
	return bytes.Equal(data1, data2), nil
}

// ValidatePDFStructure validates that a PDF file has proper structure
func ValidatePDFStructure(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open PDF file: %w", err)
	}
	defer file.Close()

	// Read first few bytes to check PDF header
	header := make([]byte, 4)
	_, err = file.Read(header)
	if err != nil {
		return fmt.Errorf("failed to read PDF header: %w", err)
	}

	// Check PDF magic number
	if string(header) != "%PDF" {
		return fmt.Errorf("invalid PDF header, expected %%PDF, got %s", string(header))
	}

	// Read entire file to check for PDF trailer
	file.Seek(0, 0)
	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read PDF content: %w", err)
	}

	// Check for EOF marker
	if !bytes.Contains(content, []byte("%%EOF")) {
		return fmt.Errorf("PDF file missing EOF marker")
	}

	return nil
}
