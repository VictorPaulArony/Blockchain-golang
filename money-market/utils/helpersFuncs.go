package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// User represents a registered user
type User struct {
	ID       string  `json:"id"`
	Email    string  `json:"email"`
	Name     string  `json:"name"`
	Phone    string  `json:"phone"`
	Password string  `json:"password"`
	Wallet   string  `json:"wallet"`
	JoinDate string  `json:"join_date"`
	Balance  float64 `json:"balance"`
}

// Transaction represents a transaction between users
type Transaction struct {
	ID        string  `json:"id"`
	Sender    string  `json:"sender"`
	Receiver  string  `json:"receiver"`
	Amount    float64 `json:"amount"`
	Timestamp string  `json:"timestamp"`
}

type MoneyMarketAccount struct {
	Wallet       string  `json:"wallet"`
	AccountType  string  `json:"account_type"` // "fixed" or "non-fixed"
	Deposit      float64 `json:"deposit"`
	InterestRate float64 `json:"interest_rate"` // Interest rate (e.g., 5% = 0.05)
	JoinDate     string  `json:"join_date"`
	FixedEndDate string  `json:"fixed_end_date"` // For fixed accounts only
	LastInterest string  `json:"last_interest"`  // Last interest calculation date
}

const MoneyMarketFile = "money_market.json"

// File paths
const (
	UserFile        = "users.json"
	TransactionFile = "transactions.json"
)

// function to add the two accounts
func Add(a, b float64) float64 {
	return a + b
}

// function to hashes the password using SHA-256
func GenerateHash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

// GenerateWallet generates a unique wallet address
func GenerateWallet(email string) string {
	hash := sha256.Sum256([]byte(email + time.Now().String()))
	return hex.EncodeToString(hash[:])
}

// LoadUsers loads users from the JSON file
func LoadUsers() []User {
	data, err := os.ReadFile(UserFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []User{}
		}
		log.Fatalf("Failed to read user file: %v", err)
	}
	var users []User
	json.Unmarshal(data, &users)
	return users
}

// SaveUsers saves users to the JSON file
func SaveUsers(users []User) {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		log.Fatalf("Failed to save user data: %v", err)
	}
	os.WriteFile(UserFile, data, 0o644)
}

// LoadTransactions loads transactions from the JSON file
func LoadTransactions() []Transaction {
	data, err := os.ReadFile(TransactionFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Transaction{}
		}
		log.Fatalf("Failed to read transactions file: %v", err)
	}
	var transactions []Transaction
	json.Unmarshal(data, &transactions)
	return transactions
}

// SaveTransactions saves transactions to the JSON file
func SaveTransactions(transactions []Transaction) {
	data, err := json.MarshalIndent(transactions, "", "  ")
	if err != nil {
		log.Fatalf("Failed to save transactions: %v", err)
	}
	os.WriteFile(TransactionFile, data, 0o644)
}

// function to save the data to the money market json db
func SaveMoneyMarketAccounts(accounts []MoneyMarketAccount) {
	data, err := json.MarshalIndent(accounts, "", " ")
	if err != nil {
		log.Printf("failed to save money market account: %v", err)
	}

	os.WriteFile(MoneyMarketFile, data, 0o644)
}

// function to load the saved money market data to the money market page
func LoadMoneyMarketAccounts() []MoneyMarketAccount {
	data, err := os.ReadFile(MoneyMarketFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []MoneyMarketAccount{}
		}
		log.Printf("Failed to read money market file: %v", err)
	}
	var accounts []MoneyMarketAccount
	json.Unmarshal(data, &accounts)
	return accounts
}

// function to add the accounts to the money market accounts
func AddMoneyMarketAccount(account MoneyMarketAccount) {
	accounts := LoadMoneyMarketAccounts()

	// If account type is fixed, calculate the FixedEndDate
	if account.AccountType == "fixed" {
		joinDate, _ := time.Parse(time.RFC3339, account.JoinDate)
		account.FixedEndDate = joinDate.AddDate(1, 0, 0).Format(time.RFC3339) // Fixed term of 1 year
	}

	accounts = append(accounts, account)
	SaveMoneyMarketAccounts(accounts)
}

// function to calculate the intrest rate for the fixed and non-fixed accounts in money market
func CalculateInterest() {
	accounts := LoadMoneyMarketAccounts()
	updated := false

	for i, account := range accounts {
		now := time.Now()

		// Parse last interest date
		lastInterestDate, _ := time.Parse(time.RFC3339, account.LastInterest)

		// Non-Fixed Accounts: Apply monthly interest
		if account.AccountType == "non-fixed" {
			if now.Sub(lastInterestDate).Hours() >= 720 { // 30 days
				interest := account.Deposit * account.InterestRate
				account.Deposit += interest
				account.LastInterest = now.Format(time.RFC3339)
				updated = true
			}
		}

		// Fixed Accounts: Apply interest at maturity
		if account.AccountType == "fixed" {
			fixedEndDate, _ := time.Parse(time.RFC3339, account.FixedEndDate)
			if now.After(fixedEndDate) && account.LastInterest == "" {
				interest := account.Deposit * account.InterestRate
				account.Deposit += interest
				account.LastInterest = now.Format(time.RFC3339) // Mark interest as applied
				updated = true
			}
		}

		accounts[i] = account
	}

	// Save updated accounts if interest was calculated
	if updated {
		SaveMoneyMarketAccounts(accounts)
		log.Println("Interest calculated and accounts updated.")
	}
}

// function to to convert time to read formart
func FormatCountdown(duration time.Duration) string {
	days := duration / (24 * time.Hour)
	hours := (duration % (24 * time.Hour)) / time.Hour
	minutes := (duration % time.Hour) / time.Minute

	return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
}
