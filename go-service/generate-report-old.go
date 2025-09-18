package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// generateStudentReport handles the PDF generation request
func generateStudentReportOLD(w http.ResponseWriter, r *http.Request) {

	// Get Node.js backend URL from environment variable or use default
	nodeAPIURL := os.Getenv("NODE_API_URL")
	if nodeAPIURL == "" {
		nodeAPIURL = "http://localhost:5007" // Default fallback
	}

	// Extract student ID from URL path
	id := r.PathValue("id")

	// Fetch student data from Node.js API
	student, err := fetchStudentData(nodeAPIURL, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch student: %v", err), http.StatusInternalServerError)
		slog.Error("Failed to fetch student", "error", err)
		return
	}

	// Generate PDF content
	pdfContent := generatePDFContent(student)

	// Set response headers for PDF download
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"student_report_%s.pdf\"", id))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfContent)))

	// // Write PDF content to response
	// _, err = w.Write([]byte(pdfContent))
	// if err != nil {
	// 	log.Printf("Error writing PDF to response: %v", err)
	// }
	http.ServeContent(w, r, "", time.Now(), strings.NewReader(pdfContent))
}

// fetchStudentData fetches student data from Node.js backend
func fetchStudentData(baseURL, id string) (*Student, error) {
	url := fmt.Sprintf("%s/api/v1/students/%s", baseURL, id)
	slog.Info("Fetching student data from", "url", url)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Get tokens from environment variables
	accessToken := os.Getenv("ACCESS_TOKEN")
	csrfToken := os.Getenv("CSRF_TOKEN")
	refreshToken := os.Getenv("REFRESH_TOKEN")
	slog.Info("token", accessToken, csrfToken, refreshToken)

	// Set cookies if tokens are available
	if accessToken != "" {
		req.AddCookie(&http.Cookie{
			Name:  "accessToken",
			Value: accessToken,
		})
	}

	if csrfToken != "" {
		req.AddCookie(&http.Cookie{
			Name:  "csrfToken",
			Value: csrfToken,
		})
	}

	if refreshToken != "" {
		req.AddCookie(&http.Cookie{
			Name:  "refreshToken",
			Value: refreshToken,
		})
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Node.js API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	slog.Info("Response body", "body", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api returned status code: %d", resp.StatusCode)
	}
	var apiResponse Student
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &apiResponse, nil
}

// generatePDFContent creates a basic PDF-like content using text formatting
func generatePDFContent(student *Student) string {
	var content strings.Builder
	content.WriteString("STUDENT REPORT\n")
	content.WriteString("Generated on: " + time.Now().Format("02-Jan-2006 15:04:05") + "\n")
	content.WriteString("---------------------------------------------------\n")
	content.WriteString("STUDENT INFORMATION\n")
	content.WriteString("ID: " + strconv.Itoa(student.ID) + "\n")
	content.WriteString("Name: " + student.Name + "\n")
	content.WriteString("Email: " + student.Email + "\n")
	content.WriteString("Phone: " + student.Phone + "\n")
	content.WriteString("Class: " + student.Class + "\n")
	content.WriteString("Section: " + student.Section + "\n")
	content.WriteString("Roll Number: " + strconv.Itoa(student.Roll) + "\n")
	content.WriteString("System Access: " + strconv.FormatBool(student.SystemAccess) + "\n")
	content.WriteString("---------------------------------------------------\n")
	content.WriteString("GUARDIAN INFORMATION\n")
	content.WriteString("Guardian Name: " + student.GuardianName + "\n")
	content.WriteString("Guardian Phone: " + student.GuardianPhone + "\n")
	content.WriteString("Relationship: " + student.RelationOfGuardian + "\n")
	content.WriteString("---------------------------------------------------\n")
	content.WriteString("ADDRESS INFORMATION\n")
	content.WriteString("Current Address: " + student.CurrentAddress + "\n")
	content.WriteString("Permanent Address: " + student.PermanentAddress + "\n")
	content.WriteString("Admission Date: " + formatDate(student.AdmissionDate) + "\n")
	content.WriteString("---------------------------------------------------\n")
	content.WriteString("REPORTED BY\n")
	content.WriteString("Reporter Name: " + student.ReporterName + "\n")
	return content.String()
}
