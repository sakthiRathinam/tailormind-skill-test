package service

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go-service/internal/config"
	"go-service/internal/models"

	"github.com/go-resty/resty/v2"
	"github.com/jung-kurt/gofpdf"
	"github.com/sirupsen/logrus"
)

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
func (s *PDFService) FetchStudentData(studentID int) (*models.Student, error) {
	logrus.Infof("Fetching student data for ID: %d", studentID)

	// add headers
	s.client.SetHeader("x-auth-token", s.config.NodeJS.AuthToken)
	s.client.SetHeader("Content-Type", "application/json")
	s.client.SetHeader("internal-service", "true")
	logrus.Infof("auth token: %s", s.config.NodeJS.AuthToken)
	resp, err := s.client.R().
		SetResult(&models.StudentResponse{}).
		SetError(map[string]interface{}{}).
		Get(fmt.Sprintf("/api/v1/students/%d", studentID))
	logrus.Infof("response: %s", resp.Body())
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

	student := resp.Result().(*models.StudentResponse).Student
	logrus.Infof("Successfully fetched data for student: %s", student.Name)

	return &student, nil
}

// createTableRow creates a properly sized table row with borders in the PDF
func createTableRow(pdf *gofpdf.Fpdf, label, value string, isHeader bool) {
	// Larger row height to fill page better
	rowHeight := 7.0
	labelWidth := 75.0
	valueWidth := 115.0

	if isHeader {
		pdf.SetFillColor(52, 73, 94)    // Same color as header/footer
		pdf.SetTextColor(255, 255, 255) // White text
		pdf.SetFont("Arial", "B", 11)
	} else {
		pdf.SetFillColor(248, 248, 248) // Very light gray background for data rows
		pdf.SetTextColor(0, 0, 0)       // Black text
		pdf.SetFont("Arial", "", 10)
	}

	// Draw label cell with border
	pdf.CellFormat(labelWidth, rowHeight, label, "1", 0, "L", true, 0, "")

	// Set font for value (always regular for readability)
	if isHeader {
		pdf.SetFont("Arial", "B", 11)
	} else {
		pdf.SetFont("Arial", "", 10)
	}

	// Draw value cell with border
	pdf.CellFormat(valueWidth, rowHeight, value, "1", 1, "L", true, 0, "")
}

// GeneratePDFReport generates a PDF report for a student
func (s *PDFService) GeneratePDFReport(student *models.Student) (string, error) {
	logrus.Infof("Generating PDF report for student: %s", student.Name)

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Header Section
	pdf.SetFillColor(52, 73, 94) // Same color as footer
	pdf.Rect(0, 0, 210, 25, "F") // Full width header background

	pdf.SetTextColor(255, 255, 255) // White text
	pdf.SetFont("Arial", "B", 18)
	pdf.SetXY(10, 8)
	pdf.Cell(0, 10, "TAILORMIND SCHOOL MANAGEMENT SYSTEM")

	pdf.SetFont("Arial", "", 12)
	pdf.SetXY(10, 18)
	pdf.Cell(0, 5, "Student Detail Report")

	// Main content area
	pdf.SetY(32)
	pdf.SetTextColor(0, 0, 0)

	// Single comprehensive table
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 8, "COMPLETE STUDENT INFORMATION")
	pdf.Ln(12)

	// Create single table with all student details
	createTableRow(pdf, "FIELD", "INFORMATION", true)

	// Personal Information
	createTableRow(pdf, "Student ID", fmt.Sprintf("%d", student.ID), false)
	createTableRow(pdf, "Full Name", student.Name, false)
	createTableRow(pdf, "Email Address", student.Email, false)
	createTableRow(pdf, "Phone Number", student.Phone, false)
	createTableRow(pdf, "Gender", student.Gender, false)
	createTableRow(pdf, "Date of Birth", student.DOB, false)
	createTableRow(pdf, "Admission Date", student.AdmissionDate, false)

	// Academic Information
	createTableRow(pdf, "Class", student.Class, false)
	createTableRow(pdf, "Section", student.Section, false)
	createTableRow(pdf, "Roll Number", fmt.Sprintf("%d", student.Roll), false)
	createTableRow(pdf, "System Access", func() string {
		if student.SystemAccess {
			return "Enabled"
		}
		return "Disabled"
	}(), false)

	// Address Information
	createTableRow(pdf, "Current Address", student.CurrentAddress, false)
	if student.PermanentAddress != "" && student.PermanentAddress != student.CurrentAddress {
		createTableRow(pdf, "Permanent Address", student.PermanentAddress, false)
	}

	// Family Information
	createTableRow(pdf, "Father's Name", student.FatherName, false)
	createTableRow(pdf, "Father's Phone", student.FatherPhone, false)
	createTableRow(pdf, "Mother's Name", student.MotherName, false)
	createTableRow(pdf, "Mother's Phone", student.MotherPhone, false)

	// Guardian Information (if different from parents)
	if student.GuardianName != "" && student.GuardianName != student.FatherName && student.GuardianName != student.MotherName {
		createTableRow(pdf, "Guardian Name", student.GuardianName, false)
		createTableRow(pdf, "Guardian Relation", student.RelationOfGuardian, false)
		createTableRow(pdf, "Guardian Phone", student.GuardianPhone, false)
	}

	// Reporter Information (if available)
	if student.ReporterName != "" {
		createTableRow(pdf, "Reporter Name", student.ReporterName, false)
	}

	// Footer Section - compact footer after table
	pdf.Ln(5) // Small gap after table

	// Ensure footer fits on same page (A4 is 297mm tall)
	footerStart := 280.0
	pdf.SetFillColor(52, 73, 94)           // Dark blue-gray for footer
	pdf.Rect(0, footerStart, 210, 20, "F") // Compact footer background

	// Footer right side - page info
	pdf.SetXY(170, footerStart+5)
	pdf.SetTextColor(255, 255, 255) // Set text color to white
	pdf.Text(170, footerStart+5, "Page 1 of 1")

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
