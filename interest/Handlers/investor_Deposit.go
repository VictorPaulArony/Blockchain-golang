package handlers

import (
	"log"
	"net/http"
	"strconv"

	helpers "interest/src"
)

// function to enable the investor to transfer funds to the money market
func TransferHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		cookie, err := r.Cookie("user_email")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}

		email := cookie.Value

		aoumtStr := r.FormValue("amount")
		amount, err := strconv.ParseFloat(aoumtStr, 64)
		if err != nil {
			log.Println("Error W hile Parsing the Investor Deposit amount")
		}

		// load the users wallets

		var wallet helpers.Wallet
		wallet.LoadData()

		var currentUser *helpers.User

		for _, user := range wallet.Users {
			if user.Email == email {
				currentUser = user
				break
			}
		}

		// check if the investor has enough amount inthe account to deposite to the money market
		if currentUser.Balance < amount {
			http.Error(w, "Insufficient funds", http.StatusForbidden)
			return
		}

		// load all the investors data from the json db
		investorData := helpers.LoadInvestorData()

		for _, investors := range investorData.Inverstors {
			for _, investor := range investors {
				if investor.Investor == email {
					currentUser.Email = investor.Investor

					newInvestor := helpers.MoneyMarketInvestor{
						ID:           investor.ID,
						Investor:     investor.Investor,
						TotalAmount:  investor.TotalAmount + amount,
						InterestRate: investor.InterestRate,
						Status:       investor.Status,
						Members:      investor.Members,
					}
					investorData.Inverstors[email] = append(investorData.Inverstors[email], newInvestor)
					break

				}
			}
		}

		currentUser.Balance -= amount
		wallet.SaveData()

		// save updated investor data
		helpers.SaveInverstors(investorData)

		http.Redirect(w, r, "/investors", http.StatusSeeOther)
		return

	}
}
