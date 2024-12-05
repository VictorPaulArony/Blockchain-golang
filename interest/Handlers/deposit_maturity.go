package handlers

import (
	"net/http"
	"text/template"
	"time"

	helpers "interest/src"
)

// MaturedDepositsHandler handles the request to show matured deposits and calculates profit.
func MaturedDepositsHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve user session from the cookie
	cookie, err := r.Cookie("user_email")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	email := cookie.Value

	// Load wallet data
	var wallet helpers.Wallet
	wallet.LoadData()

	var currentUser *helpers.User
	for _, user := range wallet.Users {
		if user.Email == email {
			currentUser = user
			break
		}
	}

	if currentUser == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Load money market data
	moneyMarket, err := helpers.LoadMoneyMarket()
	if err != nil {
		http.Error(w, "Error loading money market data", http.StatusInternalServerError)
		return
	}

	var maturedDeposits []helpers.MoneyMarketDeposit
	var totalPayout float64

	// Define interest rate
	const interestRate = 0.05 // 5% annual interest

	// Check for matured deposits in the user's MMFs
	for i, deposit := range currentUser.Mmfs {
		if deposit.Status == "Active" && deposit.MaturityDate <= time.Now().Unix() {
			// Calculate profit
			profit := deposit.Deposit * interestRate
			totalPayout += deposit.Deposit + profit

			// Update deposit status in the user's MMFs
			currentUser.Mmfs[i].Status = "Matured"
			// currentUser.Mmfs[i].Profit = profit

			// Update the same deposit in the money market
			for j, marketDeposit := range moneyMarket.Members[currentUser.Wallet] {
				if marketDeposit.Wallet == deposit.Wallet {
					moneyMarket.Members[currentUser.Wallet][j].Status = "Matured"
					break
				}
			}

			// Add to matured deposits list for display
			maturedDeposits = append(maturedDeposits, currentUser.Mmfs[i])
		}
	}

	// Update user's balance
	currentUser.Balance += totalPayout

	// Deduct total payout from the money market total
	moneyMarket.Total -= totalPayout

	// Save updated user data
	if err := wallet.SaveData(); err != nil {
		http.Error(w, "Error updating user data", http.StatusInternalServerError)
		return
	}

	// Save updated money market data
	if err := helpers.SaveMoneyMarket(moneyMarket); err != nil {
		http.Error(w, "Error updating money market data", http.StatusInternalServerError)
		return
	}

	// Prepare data for the template
	data := struct {
		User            helpers.User
		MaturedDeposits []helpers.MoneyMarketDeposit
		TotalProfit     float64
	}{
		User:            *currentUser,
		MaturedDeposits: maturedDeposits,
		TotalProfit:     totalPayout, //- float64(len(maturedDeposits)), // Assuming total profit includes deposits
	}

	// Render the matured deposits template
	tmpl := setupTemplates()
	if err := tmpl.ExecuteTemplate(w, "matured_deposits.html", data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// setupTemplates parses templates with helper functions.
func setupTemplates() *template.Template {
	funcMap := template.FuncMap{
		"unixToTime": helpers.UnixToTime,
	}

	tmpl := template.New("").Funcs(funcMap)
	tmpl, _ = tmpl.ParseGlob("templates/*.html")
	return tmpl
}
