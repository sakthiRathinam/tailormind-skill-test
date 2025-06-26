package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	// Create reference PDF for testing
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set title
	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(0, 15, "Student Report")
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
	pdf.Cell(0, 8, "John Doe")
	pdf.Ln(10)

	// Email
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Email:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, "john.doe@example.com")
	pdf.Ln(10)

	// Phone
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Phone:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, "+1234567890")
	pdf.Ln(10)

	// Address
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Address:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, "123 Main St, City, State")
	pdf.Ln(10)

	// Class Information
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Class:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, "10th Grade - Section A")
	pdf.Ln(10)

	// Admission Date
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Admission Date:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, "2023-01-15")
	pdf.Ln(10)

	// Status
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Status:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, "Active")
	pdf.Ln(20)

	// Footer
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 8, "Generated on: 2023-12-01 10:00:00")

	// Save the reference PDF
	outputPath := filepath.Join("testdata", "reference_student_report.pdf")
	
	// Ensure directory exists
	if err := os.MkdirAll("testdata", 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	if err := pdf.OutputFileAndClose(outputPath); err != nil {
		fmt.Printf("Error creating reference PDF: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Reference PDF created successfully: %s\n", outputPath)
} 