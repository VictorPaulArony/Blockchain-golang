package handlers

import (
	"log"
	"net/http"
	"text/template"
	"time"

	helpers "interest/src"
)

// Template loader
var templates = template.Must(template.New("").ParseFiles(
	"templates/login.html",
	"templates/signup.html",
	"templates/error.html",
	"templates/dashboard.html",
	"templates/index.html",
	"templates/mmf_deposit.html",
	"templates/loan_request.html",
	"templates/service.html",
	"templates/about.html",
	"templates/contact.html",
))

// Registration handles user registration
func Registration(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "signup.html", nil)
		return
	}

	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	var wallet helpers.Wallet
	wallet.LoadData()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	if password != confirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	wallet.CreateAddress(email, name, phone, password)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// LoginHandler handles both displaying the login page and processing login submissions
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var wallet helpers.Wallet
	wallet.LoadData()

	// Handle GET request to display the login form
	if r.Method == http.MethodGet {
		renderTemplate(w, "login.html", nil)
		return
	}

	// Handle POST request for processing login
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		user, exists := wallet.Users[email]
		if !exists || user.Password != helpers.GenerateHash(password) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Set cookie for the logged-in user
		http.SetCookie(w, &http.Cookie{
			Name:    "user_email",
			Value:   email,
			Expires: time.Now().Add(24 * time.Hour),
		})

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	// If the method is not GET or POST, return a method not allowed error
	ErrorHandler(w, http.StatusMethodNotAllowed)
}

// renderTemplate is a helper to render templates
func renderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	err := templates.ExecuteTemplate(w, templateName, data)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		log.Println("Template error:", err)
	}
}

// IndexHandler serves the signup page
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", nil)
}

func ServiceHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "service.html", nil)
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about.html", nil)
}

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "contact.html", nil)
}

// ErrorHandler displays error messages with corresponding status codes
func ErrorHandler(w http.ResponseWriter, code int) {
	renderTemplate(w, "error.html", map[string]int{"Code": code})
}
