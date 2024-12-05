package handlers

import (
	"net/http"
	"strconv"
	"time"

	helpers "interest/src"
)


func MoneyMarketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Retrieve user session from the cookie
		cookie, err := r.Cookie("user_email")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		email := cookie.Value

		var wallet helpers.Wallet
		wallet.LoadData()
		var userEmail *helpers.User

		for _, user := range wallet.Users {
			if user.Email == email {
				userEmail = user
				break
			}
		}

		if userEmail == nil {
			println("User not found")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Load money market data
		moneyMarket, err := helpers.LoadMoneyMarket()
		if err != nil {
			http.Error(w, "Error loading money market data", http.StatusInternalServerError)
			return
		}

		var userDeposits []helpers.MoneyMarketDeposit
		var totalBalance float64

		// Filter user's deposits and calculate the total balance
		if deposits, exists := moneyMarket.Members[userEmail.Wallet]; exists {
			userDeposits = deposits
			for _, deposit := range deposits {
				if deposit.Status == "Active" {
					totalBalance += deposit.Deposit
				}
			}
		}

		data := struct {
			User         helpers.User
			UserDeposits []helpers.MoneyMarketDeposit
			TotalBalance float64
			AllDeposits  map[string][]helpers.MoneyMarketDeposit
			MarketTotal  float64
		}{
			User:         *userEmail,
			UserDeposits: userDeposits,
			TotalBalance: totalBalance,
			AllDeposits:  moneyMarket.Members,
			MarketTotal:  moneyMarket.Total,
		}

		// Render the money market template
		templates.ExecuteTemplate(w, "mmf_deposit.html", data)
		return
	}

	// Handle POST: Join the money market
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("user_email")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := cookie.Value

	var wallet helpers.Wallet
	wallet.LoadData()

	var currentUser *helpers.User
	for _, user := range wallet.Users {
		if user.Email == email {
			currentUser = user
			break
		}
	}

	// Validate the wallet and deposit information
	depositStr := r.FormValue("deposit")
	deposit, err := strconv.ParseFloat(depositStr, 64)
	if err != nil {
		http.Error(w, "Invalid deposit amount. Minimum deposit is 100.", http.StatusBadRequest)
		return
	}

	if deposit > currentUser.Balance {
		http.Error(w, "Insufficient balance for the deposit.", http.StatusBadRequest)
		return
	}

	// Create and save the new deposit
	newDeposit := helpers.MoneyMarketDeposit{
		MoneyMarketTransaction: helpers.MoneyMarketTransaction{
			Name:     currentUser.Name,
			Wallet:   currentUser.Wallet,
			Joined:   currentUser.JoinDate,
			Maturity: 1,
		},
		Deposit:      deposit,
		Status:       "Active",
		Email:        currentUser.Email,      // Add the user's email to the deposit
		MaturityDate: time.Now().Unix() + 60, // 1 mint
	}

	currentUser.Balance -= deposit
	currentUser.Mmfs = append(currentUser.Mmfs, newDeposit)
	wallet.SaveData()

	if err := helpers.AddMoneyMarketDeposit(newDeposit); err != nil {
		http.Error(w, "Error adding money market deposit", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/money-market", http.StatusSeeOther)
}
