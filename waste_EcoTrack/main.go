package main

import (
	"log"
	"net/http"

	"waste_Eco_Track/handlers"
)

func main() {
	file := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", file))

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/resident-register", handlers.ResidentRegisterHandler)
	http.HandleFunc("/resident-login", handlers.ResidentLoginHandler)
	//http.HandleFunc("/resident-dash", handlers.)

	log.Println("server running at : http://localhost:1234")
	http.ListenAndServe(":1234", nil)
}
