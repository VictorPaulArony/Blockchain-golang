package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("templates/home.html"))
	e := temp.Execute(w, nil)
	if e != nil {
		log.Fatalln("Internal server error")
		fmt.Fprint(w, "oops something went wrong")
	}
}
