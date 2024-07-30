package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
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
		temp := template.Must(template.ParseFiles("templates/resident-register.html"))
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

}

//function that enable the residents to Login to the system
func ResidentLoginHandler(w http.ResponseWriter, r *http.Request) {

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
