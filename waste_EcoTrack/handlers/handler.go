package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
	"waste_Eco_Track/blockchain"
	"waste_Eco_Track/database"
)

var (
	requests  []database.Request
	residents []database.Resident
	//resident  database.Resident
	muSync sync.Mutex
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	var templateName string

	switch r.URL.Path {
	case "/", "/home":
		templateName = "templates/home.html"
	case "/resident-dashboard":
		templateName = "templates/resident-dashboard.html"
	case "/resident-register":
		templateName = "templates/resident-register.html"
	case "/resident-login":
		templateName = "templates/resident-login.html"
	case "/company-dashboard":
		templateName = "templates/staff-dashboard.html"
	default:
		http.NotFound(w, r)
		return
	}

	temp := template.Must(template.ParseFiles(templateName))
	if err := temp.Execute(w, nil); err != nil {
		log.Println("Internal server error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// ResidentresidentHandler allows residents to resident to the system
func ResidentRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		phone := r.FormValue("phone")
		location := r.FormValue("location")
		userID := r.FormValue("user_id")
		password := database.CreateHash(r.FormValue("password"))
	
		// Create new resident
		resident := database.Resident{
			Name:     name,
			Phone:    phone,
			UserId:   userID,
			Location: location,
			Password: password,
		}
		
		// Load existing residents
		var residents []database.Resident
		data, err := os.ReadFile(database.FileName)
		if err == nil {
			json.Unmarshal(data, &residents)
		}

		// Check if user already exists
		for _, reg := range residents {
			if reg.UserId == resident.UserId || reg.Phone == resident.Phone {
				http.Error(w, "Resident registration or phone number already exists", http.StatusConflict)
				return
			}
		}

		// Add new resident
		residents = append(residents, resident)

		// Save residents
		if err := database.SaveResident(residents); err != nil {
			http.Error(w, "Failed to save resident", http.StatusInternalServerError)
			return
		}

		// Redirect to login
		http.Redirect(w, r, "/resident-login", http.StatusSeeOther)
	} else if r.Method == http.MethodGet {
		// Render the registration form
		http.ServeFile(w, r, "templates/resident-register.html")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// ResidentLoginHandler enables residents to login to the system
func ResidentLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		userId := r.FormValue("user_id") // Changed from phone to user_id
		password := database.CreateHash(r.FormValue("password"))

		// Load existing residents
		var residents []database.Resident
		data, err := os.ReadFile(database.FileName)
		if err != nil {
			http.Error(w, "Failed to load residents", http.StatusInternalServerError)
			return
		}
		json.Unmarshal(data, &residents)

		// Authenticate resident
		for _, resident := range residents {
			if resident.UserId == userId && resident.Password == password {
				http.Redirect(w, r, "/resident-dashboard", http.StatusSeeOther)
				return
			}
		}

		http.Error(w, "Invalid user ID or password", http.StatusUnauthorized)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// StaffRegistrationHandler allows staff of the company to resident to the system
func StaffRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		phone := r.FormValue("phone")
		location := r.FormValue("location")
		password := database.CreateHash(r.FormValue("password"))

		staff := database.Staff{
			Name:     name,
			Phone:    phone,
			Location: location,
			Password: password,
		}

		// Load existing staff
		var staffs []database.Staff
		data, err := os.ReadFile(database.StaffFile)
		if err == nil {
			json.Unmarshal(data, &staffs)
		}

		// Add new staff
		staffs = append(staffs, staff)

		// Save staff
		if err := database.SaveStaff(staffs); err != nil {
			http.Error(w, "Failed to save staff", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/company-dashboard", http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// StaffLoginHandler enables staff to login to the system
func StaffLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		phone := r.FormValue("phone")
		password := database.CreateHash(r.FormValue("password"))

		// Load existing staff
		var staffs []database.Staff
		data, err := os.ReadFile(database.StaffFile)
		if err != nil {
			http.Error(w, "Failed to load staff", http.StatusInternalServerError)
			return
		}
		json.Unmarshal(data, &staffs)

		// Authenticate staff
		for _, staff := range staffs {
			if staff.Phone == phone && staff.Password == password {
				http.Redirect(w, r, "/company-dashboard", http.StatusSeeOther)
				return
			}
		}

		http.Error(w, "Invalid phone or password", http.StatusUnauthorized)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// ResidentRequestHandler allows residents to make collection requests
func ResidentRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		userId := r.FormValue("user_id")
		nature := r.FormValue("nature")
		location := r.FormValue("location")

		request := database.Request{
			ID:        len(requests) + 1,
			UserId:    userId,
			Nature:    nature,
			Location:  location,
			CreatedAt: time.Now().String(),
			Status:    "Pending",
		}

		muSync.Lock()
		requests = append(requests, request)
		muSync.Unlock()

		if err := database.SaveRequest(); err != nil {
			http.Error(w, "Failed to save request", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/resident-dashboard", http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// StaffProcessRequestHandler allows staff to process requests made by residents
func StaffProcessRequestHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	muSync.Lock()
	defer muSync.Unlock()
	var request *database.Request
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
	request.Status = "Completed"

	bc := blockchain.Blockchain{}
	if err := bc.LoadBlock(); err != nil {
		http.Error(w, fmt.Sprintf("Error loading blockchain: %v", err), http.StatusInternalServerError)
		return
	}

	// Add request to blockchain
	data := fmt.Sprintf("Request ID: %d, UserID: %s, Status: %s", request.ID, request.UserId, request.Status)
	if hash := bc.AddBlock(data); hash == "" {
		http.Error(w, "Failed to add block to blockchain", http.StatusInternalServerError)
		return
	}

	if err := bc.SaveBlock(); err != nil {
		http.Error(w, "Failed to save blockchain", http.StatusInternalServerError)
		return
	}

	if err := database.SaveRequest(); err != nil {
		http.Error(w, "Failed to save requests", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Request processed successfully")
}

// ViewRequestHandler allows staff to view requested collections
func ViewRequestHandler(w http.ResponseWriter, r *http.Request) {
	muSync.Lock()
	defer muSync.Unlock()

	temp := template.Must(template.ParseFiles("templates/view-requests.html"))
	if err := temp.Execute(w, requests); err != nil {
		log.Fatalln("Internal server error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
