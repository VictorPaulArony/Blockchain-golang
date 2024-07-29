package pdfgenerator

import (
	"fmt"

	"student-certificate-validation/registration"

	"github.com/jung-kurt/gofpdf"
)

// GeneratePDF creates a PDF for the certificate and returns the file path
func GeneratePDF(certificate registration.Certificate, hash string) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Set font and colors
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(0, 0, 0)       // Black text
	pdf.SetFillColor(255, 255, 255) // White background
	pdf.SetDrawColor(0,0,0)

	pdf.AddPage()
	
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Certificate of Completion")
	pdf.Ln(20)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Registration No: %s", certificate.RegNo))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Course: %s", certificate.Course))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Certificate Hash: %s", hash))

	filePath := fmt.Sprintf("/tmp/certificate_%s.pdf", certificate.RegNo)
	err := pdf.OutputFileAndClose(filePath)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
