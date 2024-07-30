package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
	"waste_Eco_Track/blockchain"
	"waste_Eco_Track/database"
)

var (
	muSync    sync.Mutex
	residents []database.Resident
	staffs    []database.Staff
	requests  []database.Request
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
	if r.URL.Path == "/resident-dashboard" {
		temp := template.Must(template.ParseFiles("templates/resident-dashboard.html"))
		e := temp.Execute(w, nil)
		if e != nil {
			log.Fatalln("Internal server error")
			fmt.Fprint(w, "oops something went wrong")
		}
		return
	}
	if r.URL.Path == "/resident-register" {
		temp := template.Must(template.ParseFiles("templates/resident-login.html"))
		e := temp.Execute(w, nil)
		if e != nil {
			log.Fatalln("Internal server error")
			fmt.Fprint(w, "oops something went wrong")
		}
		return
	}
	if r.URL.Path == "/resident-login" {
		temp := template.Must(template.ParseFiles("templates/resident-login.html"))
		e := temp.Execute(w, nil)
		if e != nil {
			log.Fatalln("Internal server error")
			fmt.Fprint(w, "oops something went wrong")
		}
		return
	}
	if r.URL.Path == "/Company-Dashboard" {
		temp := template.Must(template.ParseFiles("templates/staff-dashboard.html"))
		e := temp.Execute(w, nil)
		if e != nil {
			log.Fatalln("Internal server error")
			fmt.Fprint(w, "oops something went wrong")
		}
		return
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
func StaffDshboardHandler(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("templates/resident-dashboard.html"))
	e := temp.Execute(w, nil)
	if e != nil {
		log.Fatalln("Internal server error")
		fmt.Fprint(w, "oops something went wrong")
	}
}

//function to allow the residents to register to the system
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
		http.ServeFile(w, r, "templates/resident-login.html")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

//function that enable the residents to Login to the system
func ResidentLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp := template.Must(template.ParseFiles("templates/resident-login.html"))
		temp.Execute(w, nil)
		return
	}
	muSync.Lock()
	defer muSync.Unlock()

	userId := r.FormValue("user_id")
	password := database.CreateHash(r.FormValue("password"))

	// Authenticate resident
	for _, resident := range residents {
		if resident.UserId == userId && resident.Password == password {
			http.Redirect(w, r, "/resident-dashboard", http.StatusSeeOther)
			return
		}
	}

	http.Error(w, "Invalid user ID or password", http.StatusUnauthorized)

}

//function that allow the staffs of the company to register to the system
func StaffRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp := template.Must(template.ParseFiles("templates/staff-registration.html"))
		temp.Execute(w, nil)
		return
	}

	muSync.Lock()
	defer muSync.Unlock()

	name := r.FormValue("name")
	staffid := r.FormValue("staffid")
	phone := r.FormValue("phone")
	location := r.FormValue("location")
	password := database.CreateHash(r.FormValue("password"))

	staff := database.Staff{
		Name:     name,
		StaffId:  staffid,
		Phone:    phone,
		Location: location,
		Password: password,
	}

	// Add new staff
	staffs = append(staffs, staff)
	if err := database.SaveStaff(staffs); err != nil {
		http.Error(w, "Failed to save staff", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/company-dashboard", http.StatusSeeOther)
}

//functionthat enable the staffs to login to the system
func StaffLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		temp := template.Must(template.ParseFiles("templates/staff-login.html"))
		temp.Execute(w, nil)
		return
	}
	muSync.Lock()
	defer muSync.Unlock()

	staffid := r.FormValue("staffid")
	password := database.CreateHash(r.FormValue("password"))

	// Authenticate staff
	for _, staff := range staffs {
		if staff.StaffId == staffid && staff.Password == password {
			http.Redirect(w, r, "/company-dashboard", http.StatusSeeOther)
			return
		}
	}

	//
	fmt.Fprint(w, "INVALID USER PASSWORD OR ID")

}

//function that allow the resident to make collection requests
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

//function that allows thestaff to process the request made by the residents
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
