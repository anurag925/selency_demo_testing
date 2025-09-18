// go-service/main.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Student represents the student data structure
type Student struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	SystemAccess       bool   `json:"systemAccess"`
	Phone              string `json:"phone"`
	Gender             string `json:"gender"`
	Dob                string `json:"dob"`
	Class              string `json:"class"`
	Section            string `json:"section"`
	Roll               int    `json:"roll"`
	FatherName         string `json:"fatherName"`
	FatherPhone        string `json:"fatherPhone"`
	MotherName         string `json:"motherName"`
	MotherPhone        string `json:"motherPhone"`
	GuardianName       string `json:"guardianName"`
	GuardianPhone      string `json:"guardianPhone"`
	RelationOfGuardian string `json:"relationOfGuardian"`
	CurrentAddress     string `json:"currentAddress"`
	PermanentAddress   string `json:"permanentAddress"`
	AdmissionDate      string `json:"admissionDate"`
	ReporterName       string `json:"reporterName"`
}

// APIResponse wraps the student data
type APIResponse struct {
	Data Student `json:"data"`
}

func main() {
	// Get Node.js backend URL from environment variable or use default
	nodeAPIURL := os.Getenv("NODE_API_URL")
	if nodeAPIURL == "" {
		nodeAPIURL = "http://localhost:5007" // Default fallback
	}

	// Create a new HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/students/{id}/report", func(w http.ResponseWriter, r *http.Request) {
		generateStudentReport(w, r, nodeAPIURL)
	})

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting PDF Report Service on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

// generateStudentReport handles the PDF generation request
func generateStudentReport(w http.ResponseWriter, r *http.Request, nodeAPIURL string) {
	// Extract student ID from URL path
	path := r.URL.Path
	// Extract ID from path like "/api/v1/students/2/report"
	parts := strings.Split(path, "/")
	if len(parts) < 6 || parts[5] == "" {
		http.Error(w, "Student ID is required", http.StatusBadRequest)
		return
	}
	id := parts[5]

	// Fetch student data from Node.js API
	student, err := fetchStudentData(nodeAPIURL, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch student: %v", err), http.StatusInternalServerError)
		return
	}

	// Generate PDF content
	pdfContent := generatePDFContent(student)

	// Set response headers for PDF download
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"student_report_%s.pdf\"", id))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfContent)))

	// Write PDF content to response
	_, err = w.Write([]byte(pdfContent))
	if err != nil {
		log.Printf("Error writing PDF to response: %v", err)
	}
}

// fetchStudentData fetches student data from Node.js backend
func fetchStudentData(baseURL, id string) (*Student, error) {
	url := fmt.Sprintf("%s/api/v1/students/%s", baseURL, id)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Node.js API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Node.js API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResponse APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &apiResponse.Data, nil
}

// generatePDFContent creates a basic PDF-like content using text formatting
func generatePDFContent(student *Student) string {
	var content strings.Builder

	// PDF Header (simplified)
	content.WriteString("%PDF-1.4\n")
	content.WriteString("1 0 obj\n")
	content.WriteString("<</Type /Catalog\n")
	content.WriteString("/Pages 2 0 R>>\n")
	content.WriteString("endobj\n")

	// Pages object
	content.WriteString("2 0 obj\n")
	content.WriteString("<</Type /Pages\n")
	content.WriteString("/Kids [3 0 R]\n")
	content.WriteString("/Count 1>>\n")
	content.WriteString("endobj\n")

	// Page object
	content.WriteString("3 0 obj\n")
	content.WriteString("<</Type /Page\n")
	content.WriteString("/Parent 2 0 R\n")
	content.WriteString("/MediaBox [0 0 612 792]\n")
	content.WriteString("/Contents 4 0 R\n")
	content.WriteString("/Resources <<\n")
	content.WriteString("/Font <<\n")
	content.WriteString("/F1 5 0 R\n")
	content.WriteString(">>\n")
	content.WriteString(">>\n")
	content.WriteString(">>\n")
	content.WriteString("endobj\n")

	// Content stream
	content.WriteString("4 0 obj\n")
	content.WriteString("<</Length 1000>>\n")
	content.WriteString("stream\n")
	content.WriteString("BT\n")
	content.WriteString("/F1 12 Tf\n")
	content.WriteString("72 720 Td\n")
	content.WriteString("(STUDENT REPORT) Tj\n")
	content.WriteString("0 -20 Td\n")
	content.WriteString("(Generated on: " + time.Now().Format("02-Jan-2006 15:04:05") + ") Tj\n")
	content.WriteString("0 -30 Td\n")
	content.WriteString("--------------------------------------------------- Tj\n")
	content.WriteString("0 -20 Td\n")
	content.WriteString("(STUDENT INFORMATION) Tj\n")
	content.WriteString("0 -20 Td\n")
	content.WriteString("ID: " + strconv.Itoa(student.ID) + " Tj\n")
	content.WriteString("0 -15 Td\n")
	content.WriteString("Name: " + student.Name + " Tj\n")
	content.WriteString("0 -15 Td\n")
	content.WriteString("Email: " + student.Email + " Tj\n")
	content.WriteString("0 -15 Td\n")
	content.WriteString("Phone: " + student.Phone + " Tj\n")
	content.WriteString("0 -15 Td\n")
	content.WriteString("Gender: " + student.Gender + " Tj\n")
	content.WriteString("0 -15 Td\n")
	content.WriteString("Date of Birth: " + formatDate(student.Dob) + " Tj\n")
	content.WriteString("0 -15 Td\n")
	content.WriteString("Class: " + student.Class + " Tj\n")
	content.WriteString("0 -15 Td\n")
	content.WriteString("Section: " + student.Section + " Tj\n")
	content.WriteString("0 -15 Td\n")
	content.WriteString("Roll Number: " + strconv.Itoa(student.Roll) + " Tj\n")
	content.WriteString("0 -15 Td\n")
	content.WriteString("System Access: " + strconv.FormatBool(student.SystemAccess) + " Tj\n")
	content.WriteString("0 -20 Td\n")
	content.WriteString("--------------------------------------------------- Tj\n")
	content.WriteString("0 -20 Td\n")
	content.WriteString("(GUARDIAN INFORMATION) Tj\n")
	content.WriteString("0 -20 Td\n")
	content.WriteString("Guardian Name: " + student.GuardianName + " Tj\n")
	content.WriteString("0 -15 Td\n")
	content.WriteString("Guardian Phone: " + student.GuardianPhone + " Tj\n")
	content.WriteString("0 -15 Td\n")
	content.WriteString("Relationship: " + student.RelationOfGuardian + " Tj\n")
	content.WriteString("0 -20 Td\n")
	content.WriteString("--------------------------------------------------- Tj\n")
	content.WriteString("0 -20 Td\n")
	content.WriteString("(ADDRESS INFORMATION) Tj\n")
	content.WriteString("0 -20 Td\n")
	content.WriteString("Current Address: " + student.CurrentAddress + " Tj\n")
	content.WriteString("0 -15 Td\n")
	content.WriteString("Permanent Address: " + student.PermanentAddress + " Tj\n")
	content.WriteString("0 -15 Td\n")
	content.WriteString("Admission Date: " + formatDate(student.AdmissionDate) + " Tj\n")
	content.WriteString("0 -20 Td\n")
	content.WriteString("--------------------------------------------------- Tj\n")
	content.WriteString("0 -20 Td\n")
	content.WriteString("(REPORTED BY) Tj\n")
	content.WriteString("0 -20 Td\n")
	content.WriteString("Reporter Name: " + student.ReporterName + " Tj\n")
	content.WriteString("ET\n")
	content.WriteString("endstream\n")
	content.WriteString("endobj\n")

	// Font definition
	content.WriteString("5 0 obj\n")
	content.WriteString("<</Type /Font\n")
	content.WriteString("/Subtype /Type1\n")
	content.WriteString("/BaseFont /Helvetica\n")
	content.WriteString(">>\n")
	content.WriteString("endobj\n")

	// Cross-reference table
	content.WriteString("xref\n")
	content.WriteString("0 6\n")
	content.WriteString("0000000000 65535 f \n")
	content.WriteString("0000000010 00000 n \n")
	content.WriteString("0000000075 00000 n \n")
	content.WriteString("0000000145 00000 n \n")
	content.WriteString("0000000270 00000 n \n")
	content.WriteString("0000000420 00000 n \n")
	content.WriteString("trailer\n")
	content.WriteString("<</Size 6\n")
	content.WriteString("/Root 1 0 R\n")
	content.WriteString(">>\n")
	content.WriteString("startxref\n")
	content.WriteString("480\n")
	content.WriteString("%%EOF\n")

	return content.String()
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
