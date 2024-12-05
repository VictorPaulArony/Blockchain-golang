package main

import (
	"log"
	"net/http"

	handlers "interest/Handlers"
	blockchains "interest/blockchain"
)

func main() {
	blockchains.InitializeBlockchain()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/signup", handlers.Registration)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/dashboard", handlers.DashboardHandler)
	http.HandleFunc("/transactions", handlers.TransactionHandler)
	http.HandleFunc("/money-market", handlers.MoneyMarketHandler)
	http.HandleFunc("/matured-deposits", handlers.MaturedDepositsHandler)
	http.HandleFunc("/loan", handlers.LoanHandler)

	log.Println("http://localhost:1234")

	http.ListenAndServe(":1234", nil)
}
