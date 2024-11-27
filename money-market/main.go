package main

import (
	"log"
	"net/http"
	"time"

	"money-market/handlers"
	helpers "money-market/utils"
)

func main() {
	// staticDir := os.Getenv("STATIC_DIR")
	// if staticDir == "" {
	// 	staticDir = "static"
	// }
	// fs := http.FileServer(http.Dir(staticDir))
	// http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Periodically calculate interest
	go func() {
		for {
			time.Sleep(24 * time.Hour) // Run every 24 hours
			helpers.CalculateInterest()
		}
	}()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/dashboard", handlers.DashboardHandler)
	http.HandleFunc("/transactions", handlers.TransactionHandler)
	http.HandleFunc("/money-market", handlers.MoneyMarketHandler)

	log.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
