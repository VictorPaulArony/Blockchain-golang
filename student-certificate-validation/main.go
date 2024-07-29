package main

import (
	"fmt"
	"net/http"

	"student-certificate-validation/handler"
	"student-certificate-validation/registration"
)

func main() {
	fs := http.FileServer(http.Dir("templates"))
	http.Handle("/templates/", http.StripPrefix("/templates/", fs))

	// Load students from the file
	if err := registration.LoadStudents(); err != nil {
		fmt.Println("ERROR LOADING STUDENTS:", err)
		return
	}

	// Load certificate requests from the file
	if err := registration.LoadRequests(); err != nil {
		fmt.Println("ERROR LOADING REQUESTS:", err)
		return
	}

	// Load admins from the file
	if err := registration.LoadAdmins(); err != nil {
		fmt.Println("ERROR LOADING ADMINS:", err)
		return
	}

	http.HandleFunc("/", handler.HomeHandler)
	http.HandleFunc("/register", handler.RegisterStudentHandler)
	http.HandleFunc("/login", handler.LoginStudent)
	http.HandleFunc("/student-dashboard", handler.StudentDashboardHandler)
	http.HandleFunc("/request-certificate", handler.StudentCertificateRequestHandler)
	http.HandleFunc("/view-download", handler.CertificateHandler)
	http.HandleFunc("/view-request", handler.ViewRequestHandler)
	http.HandleFunc("/admin-dashboard", handler.AdminDashboardHandler)
	http.HandleFunc("/admin-login", handler.AdminLoginHandler)
	http.HandleFunc("/admin-registration", handler.AdminRegistrationHandler)
	http.HandleFunc("/process-certificate", handler.AdminProcessCertificateHandler)
	http.HandleFunc("/download-certificate", handler.DownloadCertificateHandler)

	fmt.Println("Server running at http://localhost:1234")
	http.ListenAndServe(":1234", nil)
}
