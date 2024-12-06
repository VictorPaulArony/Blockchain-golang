package handlers

import (
	"net/http"
	"strconv"
	"time"

	helpers "interest/src"
)

// RepayLoanHandler processes the repayment of loans.
func RepayLoanHandler(w http.ResponseWriter, r *http.Request) {
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

	// Load money market loans
	moneyMarket := helpers.Loadmml()

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing repayment form", http.StatusBadRequest)
			return
		}

		repaymentAmount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
		if err != nil || repaymentAmount <= 0 {
			http.Error(w, "Invalid repayment amount", http.StatusBadRequest)
			return
		}

		// Process repayments for the user's loans
		for i := range moneyMarket.Members[currentUser.Wallet] {
			loan := &moneyMarket.Members[currentUser.Wallet][i]
			if loan.Status == "Approved" && time.Now().Unix() <= loan.DueDate {
				if repaymentAmount >= loan.Amount {
					// Full repayment
					currentUser.Balance -= loan.Amount
					loan.Status = "Paid"
					repaymentAmount -= loan.Amount
				} else {
					// Partial repayment
					currentUser.Balance -= repaymentAmount
					loan.Amount -= repaymentAmount
					repaymentAmount = 0
					if loan.Amount == 0 {
						loan.Status = "Paid"
					}
				}
				break // Exit after processing one loan
			} else if loan.Status == "Approved" && time.Now().Unix() > loan.DueDate {
				loan.Status = "Defaulted" // If the loan is overdue
			}
		}

		// Save updated user and money market data
		wallet.SaveData()
		helpers.SaveMml(moneyMarket)

		http.Redirect(w, r, "/loan", http.StatusSeeOther)
		return
	}

	// Render repayment form
	data := struct {
		User helpers.User
	}{
		User: *currentUser,
	}

	tmpl := setupTemplates()
	if err := tmpl.ExecuteTemplate(w, "repay_loan.html", data); err != nil {
		http.Error(w, "Error rendering repayment template", http.StatusInternalServerError)
	}
}
