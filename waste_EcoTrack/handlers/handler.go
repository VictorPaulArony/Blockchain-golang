package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"waste_Eco_Track/database"
)

var (
	muSync    sync.Mutex
	residents []database.Resident
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Path == "/home" {
		temp := template.Must(template.ParseFiles("templates/home.html"))
		e := temp.Execute(w, nil)
		if e != nil {
			log.Fatalln("Internal server error")
			fmt.Fprint(w, "oops something went wrong")
		}
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
		http.NotFound(w, r)
		return
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

}

//functionthat enable the staffs to login to the system
func StaffLoginHandler(w http.ResponseWriter, r *http.Request) {

}

//function that allow the resident to make collection requests
func ResidentRequestHandler(w http.ResponseWriter, r *http.Request) {

}

//function that allows thestaff to process the request made by the residents
func StaffProcessRequestHandler(w http.ResponseWriter, r *http.Request) {

}

//function that allow the staff to view the requested collections
func ViewRequestHandler(w http.ResponseWriter, r *http.Request) {

}
