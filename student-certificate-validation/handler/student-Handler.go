package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	//"student-certificate-validation/blockchain"

	"student-certificate-validation/blockchain"
	"student-certificate-validation/pdfgenerator"
	"student-certificate-validation/registration"
)

var (
	students     []registration.Register
	certificates []registration.Certificate
	muSync       sync.Mutex
	student      registration.Register
)

// HomeHandler renders the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("templates/home.html"))
	temp.Execute(w, nil)
}

// func StudentDashboardHandler(w http.ResponseWriter, r *http.Request) {
// 	temp := template.Must(template.ParseFiles("templates/student.html"))
// 	temp.Execute(w, nil)
// }

// LoginStudent handles student login
func LoginStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp := template.Must(template.ParseFiles("templates/login.html"))
		temp.Execute(w, nil)
		return
	}

	muSync.Lock()
	defer muSync.Unlock()

	reg := r.FormValue("regno")
	pass := registration.HashPassword(r.FormValue("pass"))

	for _, std := range students {
		if std.RegNo == reg && std.Password == pass {
			http.Redirect(w, r, "/request-certificate", http.StatusSeeOther)
			return
		}
	}
	http.Error(w, "INVALID REGISTRATION NUMBER or PASSWORD", http.StatusUnauthorized) // Corrected error message text
}

// RegisterStudentHandler handles student registration
func RegisterStudentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp := template.Must(template.ParseFiles("templates/register.html"))
		temp.Execute(w, nil)
		return
	}

	muSync.Lock()
	defer muSync.Unlock()

	student.Name = r.FormValue("name")
	student.RegNo = r.FormValue("regno")
	student.Course = r.FormValue("course")
	student.Email = r.FormValue("email")
	student.Phone = r.FormValue("phone")
	student.Password = registration.HashPassword(r.FormValue("pass"))

	// Check if student already exists
	for _, reg := range students {
		if reg.RegNo == student.RegNo || reg.Email == student.Email {
			http.Error(w, "Student Registration or Email Already Exist", http.StatusConflict)
			return
		}
	}

	// Add new student to the database
	students = append(students, student)

	if err := registration.AddStudent(students); err != nil { // Added students as a parameter
		http.Error(w, "Error Saving the Student", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

var (
	admins     []registration.Admin
	adminMutex sync.Mutex
	admin      registration.Admin
)

func AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("templates/admin-home.html"))
	temp.Execute(w, nil)
}

func AdminRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp := template.Must(template.ParseFiles("templates/admin-registration.html"))
		temp.Execute(w, nil)
		return
	}
	adminMutex.Lock()
	defer adminMutex.Unlock()

	admin.AdminId = r.FormValue("adminid")
	admin.Name = r.FormValue("name")
	admin.Department = r.FormValue("department")
	admin.Phone = r.FormValue("phone")
	admin.Email = r.FormValue("email")
	admin.Password = registration.HashPassword(r.FormValue("pass"))

	for _, ad := range admins {
		if ad.AdminId == admin.AdminId || ad.Email == admin.Email {
			http.Error(w, "Admin Registration number or Email already exsist", http.StatusConflict)
			return
		}
	}

	admins = append(admins, admin)
	if err := registration.AddAdmin(admins); err != nil {
		http.Error(w, "Error saving the Admin", http.StatusConflict)
		return
	}
	http.Redirect(w, r, "/admin-dashboard", http.StatusSeeOther)
}

func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp := template.Must(template.ParseFiles("templates/admin-login.html"))
		temp.Execute(w, nil)
		return
	}
	adminMutex.Lock()
	defer adminMutex.Unlock()

	adminId := r.FormValue("adminid")
	password := registration.HashPassword(r.FormValue("pass"))

	for _, ad := range admins {
		if ad.AdminId == adminId && ad.Password == password {
			http.Redirect(w, r, "/admin-dashboard", http.StatusSeeOther)
			return
		}
	}
	http.Error(w, "INVALID ADMIN ID OR PASSWORD", http.StatusUnauthorized)
}

var (
	requests     []registration.CertificateRequest
	requestMutex sync.Mutex
)

// StudentCertificateRequestHandler handles certificate request submissions by students
func StudentCertificateRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp := template.Must(template.ParseFiles("templates/student.html"))
		temp.Execute(w, nil)
		return
	} else if r.Method == http.MethodPost {
		requestMutex.Lock()
		defer requestMutex.Unlock()
		r.ParseForm()
		request := registration.CertificateRequest{
			ID:        len(requests) + 1,
			Name:      r.FormValue("name"),
			RegNo:     r.FormValue("regno"),
			Course:    r.FormValue("course"),
			CreatedAt: time.Now().String(),
			Status:    "Pending",
		}

		requests = append(requests, request)

		if err := registration.SaveRequests(requests); err != nil {
			http.Error(w, "Error Saving the Request", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/student-dashboard", http.StatusSeeOther)
	}
}

// function to enable the admin to view the certificates to aprove there processing
func ViewRequestHandler(w http.ResponseWriter, r *http.Request) {
	muSync.Lock()
	defer muSync.Unlock()
	temp := template.Must(template.ParseFiles("templates/view-request.html"))
	if err := temp.Execute(w, requests); err != nil {
		http.Error(w, "ERROR RENDERING TEMPLATE: ", http.StatusInternalServerError)
	}
}

// function to allow the download certificate after the request is made
func CertificateHandler(w http.ResponseWriter, r *http.Request) {
	muSync.Lock()
	defer muSync.Unlock()
	temp := template.Must(template.ParseFiles("templates/download.html"))
	if err := temp.Execute(w, requests); err != nil {
		http.Error(w, "ERROR RENDERING FILE: ", http.StatusInternalServerError)
	}
}

// AdminProcessCertificateHandler handles processing of certificate requests by admin
func AdminProcessCertificateHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	requestMutex.Lock()
	defer requestMutex.Unlock()

	var request *registration.CertificateRequest
	for i := range requests {
		if requests[i].ID == id {
			request = &requests[i]
			break
		}
	}
	if request == nil {
		http.Error(w, "Request not found", http.StatusNotFound)
		return
	}

	if request.Status != "Pending" {
		http.Error(w, "Request already processed", http.StatusBadRequest)
		return
	}

	// Update status and generate certificate
	request.Status = "Completed"
	certificate := registration.Certificate{
		ID:        len(certificates) + 1,
		Name:      request.Name,
		RegNo:     request.RegNo,
		Course:    request.Course,
		CreatedAt: time.Now().String(),
	}

	// Add the certificate to the blockchain
	bc := blockchain.Blockchain{}
	if err := bc.LoadBlockchain(); err != nil {
		http.Error(w, fmt.Sprintf("Error loading blockchain: %v", err), http.StatusInternalServerError)
		return
	}

	bc.AddBlock(certificate.Name)
	newBlock := bc.Certificates[len(bc.Certificates)-1]
	certificate.Hash = newBlock.Hash

	// Append the certificate to the certificates list
	certificates = append(certificates, certificate)

	// Generate the PDF
	filePath, err := pdfgenerator.GeneratePDF(certificate, certificate.Hash)
	if err != nil {
		http.Error(w, "Error generating PDF", http.StatusInternalServerError)
		return
	}

	// Save the updated requests and certificates
	if err := registration.SaveRequests(requests); err != nil {
		http.Error(w, "Error saving requests", http.StatusInternalServerError)
		return
	}

	if err := registration.SaveCertificates(certificates); err != nil {
		http.Error(w, "Error saving certificates", http.StatusInternalServerError)
		return
	}

	if err := bc.SaveBlocks(); err != nil {
		http.Error(w, fmt.Sprintf("Error saving blockchain: %v", err), http.StatusInternalServerError)
		return
	}

	// Send the PDF back to the student (for simplicity, we'll provide a download link)
	http.Redirect(w, r, "/download?file="+filePath, http.StatusSeeOther)
}

// StudentDashboardHandler renders the student dashboard page.
func StudentDashboardHandler(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("templates/student.html"))

	// data := map[string]interface{}{
	// 	"Certificates": certificates,
	// }

	err := temp.Execute(w, nil)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
		return
	}
}

// DownloadCertificateHandler serves the PDF certificate file.
func DownloadCertificateHandler(w http.ResponseWriter, r *http.Request) {
	certIDStr := r.URL.Query().Get("id")
	certID, err := strconv.Atoi(certIDStr)
	if err != nil {
		http.Error(w, "Invalid certificate ID", http.StatusBadRequest)
		return
	}

	var certificate registration.Certificate
	for _, cert := range certificates {
		if cert.ID == certID {
			certificate = cert
			break
		}
	}

	if (certificate == registration.Certificate{}) {
		http.Error(w, "Certificate not found", http.StatusNotFound)
		return
	}

	filePath, err := pdfgenerator.GeneratePDF(certificate, certificate.Hash)
	if err != nil {
		http.Error(w, "Error generating PDF", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, filePath)
}
