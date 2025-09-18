// go-service/main.go
package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
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

	// Create a new HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/students/{id}/report-old", func(w http.ResponseWriter, r *http.Request) {
		generateStudentReportOLD(w, r)
	})

	mux.HandleFunc("/api/v1/students/{id}/report", generateStudentReport)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting PDF Report Service on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
