package main

import (
	"log"
	"net/http"
	"time"

	blockchains "money-market/blockchain"
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

	blockchains.InitializeBlockchain()

	go func() {
		for {
			time.Sleep(1 * time.Minute) // Run every 1 minute
			helpers.CalculateInterest()
			helpers.UpdateMarketTrends() //for market trends
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
	http.HandleFunc("/market-trends", handlers.MarketTrendsHandler) 

	log.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
