package service

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"go-service/internal/config"
	"go-service/internal/models"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	// Set log level to error to reduce noise during tests
	logrus.SetLevel(logrus.ErrorLevel)

	// Load environment variables
	err := godotenv.Load("../../config.env")
	if err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}
}

// TestFetchStudentData tests the student data fetching functionality
func TestFetchStudentData(t *testing.T) {

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	expectedStudent := &models.Student{
		ID:                 1,
		Name:               "John Doe",
		Email:              "john.doe@school.com",
		Phone:              "+1234567890",
		Gender:             "Male",
		DOB:                "2005-01-15T00:00:00.000Z",
		Class:              "10th Grade",
		Section:            "A",
		Roll:               101,
		FatherName:         "Robert Doe",
		FatherPhone:        "+1234567891",
		MotherName:         "Mary Doe",
		MotherPhone:        "+1234567892",
		GuardianName:       "Robert Doe",
		GuardianPhone:      "+1234567891",
		RelationOfGuardian: "Father",
		CurrentAddress:     "123 Main St, City, State 12345",
		PermanentAddress:   "123 Main St, City, State 12345",
		AdmissionDate:      "2023-01-15T00:00:00.000Z",
		SystemAccess:       true,
		ReporterName:       "Admin User",
	}

	// Create PDF service
	service := NewPDFService(cfg)

	// Test successful fetch
	t.Run("SuccessfulFetch", func(t *testing.T) {
		student, err := service.FetchStudentData(1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if !reflect.DeepEqual(student, expectedStudent) {
			t.Errorf("Expected student data to match, got:\nExpected: %+v\nActual: %+v", expectedStudent, student)
		}

	})

	// Test not found
	t.Run("StudentNotFound", func(t *testing.T) {
		_, err := service.FetchStudentData(9999999)
		if err == nil {
			t.Fatal("Expected error for non-existent student, got nil")
		}

		// Verify the error message contains expected content
		if !strings.Contains(err.Error(), "404") && !strings.Contains(err.Error(), "not found") {
			t.Errorf("Expected error to contain 404 or 'not found', got: %v", err)
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

func TestEntireWorkflow(t *testing.T) {
	tempDir := t.TempDir()
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	cfg.PDF.OutputDir = tempDir
	service := NewPDFService(cfg)

	// Test the entire workflow: fetch student data and generate PDF
	t.Run("FetchAndGeneratePDF", func(t *testing.T) {
		studentID := 1

		// Fetch student data
		student, err := service.FetchStudentData(studentID)
		if err != nil {
			t.Fatalf("Failed to fetch student data: %v", err)
		}

		// Generate PDF report
		filePath, err := service.GeneratePDFReport(student)
		if err != nil {
			t.Fatalf("Failed to generate PDF report: %v", err)
		}

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Fatalf("Expected PDF file to exist at %s", filePath)
		}

		// Read the generated PDF
		generatedPDF, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read generated PDF: %v", err)
		}

		// Read the expected PDF from testdata
		expectedPDFPath := "../../testdata/dummy_pdf_report.pdf"
		expectedPDF, err := os.ReadFile(expectedPDFPath)
		if err != nil {
			t.Fatalf("Failed to read expected PDF: %v", err)
		}

		// Compare PDFs with tolerance for dynamic content like timestamps
		comparePDFsWithTolerance(t, generatedPDF, expectedPDF, 0.05) // Allow 5% byte differences

		t.Logf("Successfully completed workflow: fetched student %d and generated matching PDF", studentID)
	})
}

// comparePDFsWithTolerance compares two PDFs allowing for some byte differences
// This accounts for dynamic content like timestamps while still validating structure
func comparePDFsWithTolerance(t *testing.T, generated, expected []byte, tolerance float64) {
	// Log file sizes for debugging
	t.Logf("Generated PDF size: %d bytes", len(generated))
	t.Logf("Expected PDF size: %d bytes", len(expected))

	// Check if files are roughly the same size (within 10% difference)
	sizeDiff := float64(abs(len(generated)-len(expected))) / float64(max(len(generated), len(expected)))
	if sizeDiff > 0.1 {
		t.Errorf("PDF sizes differ too much: generated=%d, expected=%d (%.2f%% difference)",
			len(generated), len(expected), sizeDiff*100)
		return
	}

	// Compare byte by byte and count differences
	minLen := min(len(generated), len(expected))
	maxLen := max(len(generated), len(expected))
	differences := 0
	firstDiff := -1

	// Count differences in overlapping portion
	for i := 0; i < minLen; i++ {
		if generated[i] != expected[i] {
			differences++
			if firstDiff == -1 {
				firstDiff = i
			}
		}
	}

	// Add length difference as additional differences
	differences += maxLen - minLen

	// Calculate difference percentage
	diffPercent := float64(differences) / float64(maxLen)

	t.Logf("Byte differences: %d out of %d bytes (%.2f%%)", differences, maxLen, diffPercent*100)

	if firstDiff != -1 {
		t.Logf("First difference at byte %d: generated=0x%02x, expected=0x%02x",
			firstDiff, generated[firstDiff], expected[firstDiff])
	}

	// Check if within tolerance
	if diffPercent > tolerance {
		t.Errorf("PDF differences exceed tolerance: %.2f%% > %.2f%%", diffPercent*100, tolerance*100)
	} else {
		t.Logf("PDF comparison passed within tolerance: %.2f%% <= %.2f%%", diffPercent*100, tolerance*100)
	}
}

// Helper functions for min/max/abs
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
