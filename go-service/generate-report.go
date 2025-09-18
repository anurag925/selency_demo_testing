package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// generateStudentReport handles the PDF generation request
func generateStudentReport(w http.ResponseWriter, r *http.Request) {
	nodeAPIURL := getEnv("NODE_API_URL", "http://localhost:5007")
	id := r.PathValue("id")

	student, err := fetchStudentDataInternal(nodeAPIURL, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch student: %v", err), http.StatusInternalServerError)
		slog.Error("Failed to fetch student", "error", err)
		return
	}

	pdf := generatePDF(student)

	// Set headers for PDF download
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"student_report_%s.pdf\"", id))

	err = pdf.Output(w)
	if err != nil {
		slog.Error("Failed to write PDF output", "error", err)
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
	}
}

// getEnv returns env variable or default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// fetchStudentData fetches student data from Node.js backend
func fetchStudentDataInternal(baseURL, id string) (*Student, error) {
	url := fmt.Sprintf("%s/api/v1/internals/students/%s", baseURL, id)
	slog.Info("Fetching student data", "url", url)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("x-service-token", os.Getenv("SERVICE_TOKEN"))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Node.js API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	var apiResponse Student
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &apiResponse, nil
}

// // addAuthCookies adds authentication cookies to the request
// func addAuthCookies(req *http.Request) {
// 	accessToken := os.Getenv("ACCESS_TOKEN")
// 	csrfToken := os.Getenv("CSRF_TOKEN")
// 	refreshToken := os.Getenv("REFRESH_TOKEN")

// 	if accessToken != "" {
// 		req.AddCookie(&http.Cookie{Name: "accessToken", Value: accessToken})
// 	}
// 	if csrfToken != "" {
// 		req.AddCookie(&http.Cookie{Name: "csrfToken", Value: csrfToken})
// 	}
// 	if refreshToken != "" {
// 		req.AddCookie(&http.Cookie{Name: "refreshToken", Value: refreshToken})
// 	}
// }

// generatePDF generates a PDF report for the student
func generatePDF(student *Student) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "STUDENT REPORT")
	pdf.Ln(12)

	// Timestamp
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, "Generated on: "+time.Now().Format("02-Jan-2006 15:04:05"))
	pdf.Ln(10)

	// Divider
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "---------------------------------------------------")
	pdf.Ln(8)

	// Student Info
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "STUDENT INFORMATION")
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("ID: %d", student.ID))
	pdf.Ln(6)
	pdf.Cell(40, 10, "Name: "+student.Name)
	pdf.Ln(6)
	pdf.Cell(40, 10, "Email: "+student.Email)
	pdf.Ln(6)
	pdf.Cell(40, 10, "Phone: "+student.Phone)
	pdf.Ln(6)
	pdf.Cell(40, 10, "Class: "+student.Class)
	pdf.Ln(6)
	pdf.Cell(40, 10, "Section: "+student.Section)
	pdf.Ln(6)
	pdf.Cell(40, 10, "Roll Number: "+strconv.Itoa(student.Roll))
	pdf.Ln(6)
	pdf.Cell(40, 10, "System Access: "+strconv.FormatBool(student.SystemAccess))
	pdf.Ln(8)

	// Guardian Info
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "GUARDIAN INFORMATION")
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, "Guardian Name: "+student.GuardianName)
	pdf.Ln(6)
	pdf.Cell(40, 10, "Guardian Phone: "+student.GuardianPhone)
	pdf.Ln(6)
	pdf.Cell(40, 10, "Relationship: "+student.RelationOfGuardian)
	pdf.Ln(8)

	// Address Info
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "ADDRESS INFORMATION")
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, "Current Address: "+student.CurrentAddress)
	pdf.Ln(6)
	pdf.Cell(40, 10, "Permanent Address: "+student.PermanentAddress)
	pdf.Ln(6)
	pdf.Cell(40, 10, "Admission Date: "+formatDate(student.AdmissionDate))
	pdf.Ln(8)

	// Reporter Info
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "REPORTED BY")
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, "Reporter Name: "+student.ReporterName)

	return pdf
}

// formatDate converts ISO date string to readable format
func formatDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr // Return original if parsing fails
	}
	return t.Format("02-Jan-2006")
}
