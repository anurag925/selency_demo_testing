# Go Service - Student Report Generator

## Getting Started

### Prerequisites
- Go installed on your system
- curl command-line tool

### Running the Service

1. Navigate to the service directory:
   ```bash
   cd go-service
   ```

2. Start the application:
   ```bash
   go run .
   ```

### Downloading Student Reports

Once the service is running, you can download PDF reports for students using the following command:

```bash
curl -X GET http://localhost:8080/api/v1/students/2/report -o student_report.pdf
```

This will download the report for student ID 2 and save it as `student_report.pdf`.

### API Endpoint
- **GET** `/api/v1/students/{id}/report`
- Downloads a PDF report for the specified student ID
```

This README provides clear instructions on how to run the service and download student reports, with properly formatted code blocks and structured sections.