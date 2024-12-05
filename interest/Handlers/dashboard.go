package handlers

import (
	"net/http"

	helpers "interest/src"
	// Assume there's a package for blockchain handling
	blockchains "interest/blockchain"
)

// DashboardHandler serves the user dashboard
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Get the logged-in user's email from the cookie
	cookie, err := r.Cookie("user_email")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	email := cookie.Value

	// Load user data
	var wallet helpers.Wallet
	wallet.LoadData()

	// Fetch the current user's details
	currentUser, exists := wallet.Users[email]
	if !exists {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Load all transactions
	allTransactions := helpers.LoadTransactions()

	// Determine the filter based on the query parameter
	viewAll := r.URL.Query().Get("view") == "all"
	var filteredTransactions []helpers.Transaction

	if viewAll {
		// Show all transactions
		filteredTransactions = allTransactions
	} else {
		// Show only transactions related to the user's wallet address
		for _, tx := range allTransactions {
			if tx.Sender == currentUser.Wallet || tx.Receiver == currentUser.Wallet {
				filteredTransactions = append(filteredTransactions, tx)
			}
		}
	}

	// Limit the displayed transactions to the 10 most recent
	if len(filteredTransactions) > 10 {
		filteredTransactions = filteredTransactions[len(filteredTransactions)-10:]
	}

	// Load the blockchain data
	blockchain := blockchains.LoadBlockchain()

	// Prepare the data for the dashboard template
	data := struct {
		User         helpers.User
		Blockchain   blockchains.Blockchain
		MempoolSize  int
		Transactions []helpers.Transaction
		ViewingAll   bool
	}{
		User:         *currentUser,
		Blockchain:   blockchain,
		MempoolSize:  len(blockchain.Mempool),
		Transactions: filteredTransactions,
		ViewingAll:   viewAll,
	}

	// Render the dashboard template
	renderTemplate(w, "dashboard.html", data)
}
