package models

import (
	"time"
)

// Student represents the student data structure from the Node.js API
type Student struct {
	ID                    int    `json:"id"`
	Name                  string `json:"name"`
	Email                 string `json:"email"`
	SystemAccess          bool   `json:"systemAccess"`
	Phone                 string `json:"phone"`
	Gender                string `json:"gender"`
	DOB                   string `json:"dob"`
	Class                 string `json:"class"`
	Section               string `json:"section"`
	Roll                  int    `json:"roll"`
	FatherName            string `json:"fatherName"`
	FatherPhone           string `json:"fatherPhone"`
	MotherName            string `json:"motherName"`
	MotherPhone           string `json:"motherPhone"`
	GuardianName          string `json:"guardianName"`
	GuardianPhone         string `json:"guardianPhone"`
	RelationOfGuardian    string `json:"relationOfGuardian"`
	CurrentAddress        string `json:"currentAddress"`
	PermanentAddress      string `json:"permanentAddress"`
	AdmissionDate         string `json:"admissionDate"`
	ReporterName          string `json:"reporterName"`
}

// StudentResponse represents the API response wrapper
type StudentResponse struct {
	Student Student `json:"student,omitempty"`
	Message string  `json:"message,omitempty"`
	Error   string  `json:"error,omitempty"`
}

// PDFReportRequest represents a request to generate a PDF report
type PDFReportRequest struct {
	StudentID string            `json:"student_id"`
	Options   PDFReportOptions  `json:"options,omitempty"`
}

// PDFReportOptions represents options for PDF generation
type PDFReportOptions struct {
	Title       string `json:"title,omitempty"`
	IncludeLogo bool   `json:"include_logo,omitempty"`
	Template    string `json:"template,omitempty"`
}

// PDFReportResponse represents the response for PDF generation
type PDFReportResponse struct {
	FileName  string    `json:"file_name"`
	Size      int       `json:"size"`
	Generated time.Time `json:"generated"`
	StudentID string    `json:"student_id"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
} 