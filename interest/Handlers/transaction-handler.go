package handlers

import (
	"net/http"
	"strconv"
	"time"

	blockchains "interest/blockchain"
	helpers "interest/src"

	"github.com/google/uuid"
)

// function to handle the transactions usage in a handler
func TransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		senderWallet := r.FormValue("sender_wallet")
		receiverWallet := r.FormValue("receiver_wallet")
		amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)

		// Load the wallets of the users
		var wallet helpers.Wallet
		wallet.LoadData()

		// Finding the sender and receiver by wallet address
		var sender *helpers.User
		var receiver *helpers.User

		for _, user := range wallet.Users {
			if user.Wallet == senderWallet {
				sender = user
			}
			if user.Wallet == receiverWallet {
				receiver = user
			}
		}

		// Check if sender and receiver were found in db
		if sender == nil || receiver == nil {
			http.Error(w, "Sender or receiver does not exist", http.StatusBadRequest)
			return
		}

		// Check if sender has enough balance
		if sender.Balance < amount {
			http.Error(w, "Insufficient balance", http.StatusBadRequest)
			return
		}

		// Create the transaction
		transaction := helpers.Transaction{
			ID:        uuid.New().String(),
			Sender:    senderWallet,
			Receiver:  receiverWallet,
			Amount:    amount,
			Timestamp: time.Now().Format(time.RFC3339),
		}

		// Add transaction to the mempool
		blockchains.AddTransactionToMempool(transaction)

		// Deduct the amount from sender's balance and add to receiver's balance
		sender.Balance -= amount
		receiver.Balance += amount

		// Append the transaction to the users' transaction slices
		sender.Transactions = append(sender.Transactions, transaction)
		receiver.Transactions = append(receiver.Transactions, transaction)

		// Save updated wallet data
		if err := wallet.SaveData(); err != nil {
			http.Error(w, "Error saving wallet data", http.StatusInternalServerError)
			return
		}

		// Save transaction to the transaction database
		if err := helpers.SaveTransaction(transaction); err != nil {
			http.Error(w, "Error saving transaction", http.StatusInternalServerError)
			return
		}

		// Mine a new block to include the transaction
		blockchains.MineBlock()

		// Redirect to the dashboard
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}
