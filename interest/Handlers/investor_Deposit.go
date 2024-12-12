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
		if currentUser == nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// load all the investors data from the json db
		investorData := helpers.LoadInvestorData()

		// Find and update the corresponding investor
		investor, exists := investorData.Investors[email]
		if !exists {
			http.Error(w, "Investor not found", http.StatusNotFound)
			return
		}

		// for _, investors := range investorData.Investors {
		// for _, investor := range investorData.Investors {
		// 	if investor.Name == email {
		// 		currentUser.Email = investor.Name

		// 		newInvestor := helpers.MoneyMarketInvestor{
		// 			ID:           investor.ID,
		// 			Name:         investor.Name,
		// 			Amount:       investor.Amount + amount,
		// 			InterestRate: investor.InterestRate,
		// 			Status:       investor.Status,
		// 			// Members:      investor.Members,
		// 		}
		// 		// investorData.Investors[email] = newInvestor
		// 		// investorData.Investors[email].Email = append(investorData.Investors[email].Email, newInvestor)
		// 		break

		// 	}
		// }

		investor.Amount += amount

		currentUser.Balance -= amount
		wallet.SaveData()

		// save updated investor data
		helpers.SaveInvestors(investorData)

		http.Redirect(w, r, "/investor", http.StatusSeeOther)
		return

	}
}
