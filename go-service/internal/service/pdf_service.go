package service

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go-service/internal/config"

	"github.com/go-resty/resty/v2"
	"github.com/jung-kurt/gofpdf"
	"github.com/sirupsen/logrus"
)

// Student represents the student data structure from the Node.js API
type Student struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	ClassID   int    `json:"class_id"`
	ClassName string `json:"class_name"`
	Section   string `json:"section"`
	AdmissionDate string `json:"admission_date"`
	Status    string `json:"status"`
}

type PDFService struct {
	client *resty.Client
	config *config.Config
}

// NewPDFService creates a new PDF service instance
func NewPDFService(cfg *config.Config) *PDFService {
	client := resty.New()
	client.SetTimeout(cfg.NodeJS.Timeout)
	client.SetBaseURL(cfg.NodeJS.BaseURL)

	return &PDFService{
		client: client,
		config: cfg,
	}
}

// FetchStudentData fetches student data from the Node.js API
func (s *PDFService) FetchStudentData(studentID int) (*Student, error) {
	logrus.Infof("Fetching student data for ID: %d", studentID)

	resp, err := s.client.R().
		SetResult(&Student{}).
		SetError(map[string]interface{}{}).
		Get(fmt.Sprintf("/api/v1/students/%d", studentID))

	if err != nil {
		logrus.WithError(err).Error("Failed to fetch student data")
		return nil, fmt.Errorf("failed to fetch student data: %w", err)
	}

	if resp.StatusCode() != 200 {
		logrus.WithFields(logrus.Fields{
			"status_code": resp.StatusCode(),
			"response":    string(resp.Body()),
		}).Error("Non-200 response from API")
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	student := resp.Result().(*Student)
	logrus.Infof("Successfully fetched data for student: %s %s", student.FirstName, student.LastName)
	
	return student, nil
}

// GeneratePDFReport generates a PDF report for a student
func (s *PDFService) GeneratePDFReport(student *Student) (string, error) {
	logrus.Infof("Generating PDF report for student: %s %s", student.FirstName, student.LastName)

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set title
	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(0, 15, s.config.PDF.Title)
	pdf.Ln(20)

	// Student Information Header
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Student Information")
	pdf.Ln(15)

	// Student Details
	pdf.SetFont("Arial", "", 12)
	
	// Name
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Name:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("%s %s", student.FirstName, student.LastName))
	pdf.Ln(10)

	// Email
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Email:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, student.Email)
	pdf.Ln(10)

	// Phone
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Phone:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, student.Phone)
	pdf.Ln(10)

	// Address
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Address:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, student.Address)
	pdf.Ln(10)

	// Class Information
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Class:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, fmt.Sprintf("%s - Section %s", student.ClassName, student.Section))
	pdf.Ln(10)

	// Admission Date
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Admission Date:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, student.AdmissionDate)
	pdf.Ln(10)

	// Status
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Status:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, student.Status)
	pdf.Ln(20)

	// Footer
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 8, fmt.Sprintf("Generated on: %s", time.Now().Format("2006-01-02 15:04:05")))

	// Ensure output directory exists
	if err := os.MkdirAll(s.config.PDF.OutputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate filename
	filename := fmt.Sprintf("student_%d_report_%s.pdf", student.ID, time.Now().Format("20060102_150405"))
	filepath := filepath.Join(s.config.PDF.OutputDir, filename)

	// Save PDF
	if err := pdf.OutputFileAndClose(filepath); err != nil {
		logrus.WithError(err).Error("Failed to save PDF")
		return "", fmt.Errorf("failed to save PDF: %w", err)
	}

	logrus.Infof("PDF report generated successfully: %s", filepath)
	return filepath, nil
}

// GenerateStudentReport is the main function to generate a complete student report
func (s *PDFService) GenerateStudentReport(studentID int) (string, error) {
	// Fetch student data
	student, err := s.FetchStudentData(studentID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch student data: %w", err)
	}

	// Generate PDF report
	filepath, err := s.GeneratePDFReport(student)
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF report: %w", err)
	}

	return filepath, nil
} 