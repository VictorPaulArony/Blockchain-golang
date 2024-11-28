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

func Mul(a, b float64) float64 {
	return a * b
}

func Format(date string) string {
	parsed, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return date
	}
	return parsed.Format("02 Jan 2006 15:04")
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
// func CalculateInterest() {
// 	accounts := LoadMoneyMarketAccounts()
// 	updated := false

// 	for i, account := range accounts {
// 		now := time.Now()

// 		// Parse last interest date
// 		lastInterestDate, _ := time.Parse(time.RFC3339, account.LastInterest)

// 		// Non-Fixed Accounts: Apply monthly interest
// 		if account.AccountType == "non-fixed" {
// 			if now.Sub(lastInterestDate).Minutes() >= 720 { // 30 days
// 				interest := account.Deposit * account.InterestRate
// 				account.Deposit += interest
// 				account.LastInterest = now.Format(time.RFC3339)
// 				updated = true
// 			}
// 		}

// 		// Fixed Accounts: Apply interest at maturity
// 		if account.AccountType == "fixed" {
// 			fixedEndDate, _ := time.Parse(time.RFC3339, account.FixedEndDate)
// 			if now.After(fixedEndDate) && account.LastInterest == "" {
// 				interest := account.Deposit * account.InterestRate
// 				account.Deposit += interest
// 				account.LastInterest = now.Format(time.RFC3339) // Mark interest as applied
// 				updated = true
// 			}
// 		}

// 		accounts[i] = account
// 	}

// 	// Save updated accounts if interest was calculated
// 	if updated {
// 		SaveMoneyMarketAccounts(accounts)
// 		log.Println("Interest calculated and accounts updated.")
// 	}
// }

func CalculateInterest() {
	accounts := LoadMoneyMarketAccounts()
	users := LoadUsers()
	totalMoneyMarket := 0.0
	updated := false

	for i, account := range accounts {
		now := time.Now()

		// Parse last interest date or initialize it to the join date
		lastInterestDate, _ := time.Parse(time.RFC3339, account.LastInterest)
		if account.LastInterest == "" {
			lastInterestDate, _ = time.Parse(time.RFC3339, account.JoinDate)
		}

		// Find the user associated with this account
		var user *User
		for j := range users {
			if users[j].Wallet == account.Wallet {
				user = &users[j]
				break
			}
		}

		// Skip if user not found (data integrity issue)
		if user == nil {
			log.Printf("User not found for wallet %s. Skipping account.", account.Wallet)
			continue
		}

		// Non-Fixed Accounts: Apply interest every 1 minute (testing)
		if account.AccountType == "non-fixed" {
			if now.Sub(lastInterestDate).Minutes() >= 1 { // Testing interval: 1 minute
				interest := account.Deposit * account.InterestRate
				user.Balance += interest    // Add interest to user's main balance
				account.Deposit -= interest // Subtract from money market total
				account.LastInterest = now.Format(time.RFC3339)
				updated = true
				// log.Printf("Interest of %.2f added to user %s's main balance for non-fixed account.", interest, user.Email)
			}
		}

		// Fixed Accounts: Apply interest after maturity (testing: 1 minute)
		if account.AccountType == "fixed" {
			fixedEndDate, _ := time.Parse(time.RFC3339, account.FixedEndDate)
			if now.After(fixedEndDate) && account.LastInterest == "" {
				// interest := account.Deposit * account.InterestRate
				interest := account.Deposit * 10
				user.Balance += account.Deposit + interest      // Add deposit + interest to main balance
				account.Deposit = 0                             // Clear deposit after maturity
				account.LastInterest = now.Format(time.RFC3339) // Mark interest as applied
				updated = true
				// log.Printf("Matured interest of %.2f added to user %s's main balance for fixed account.", interest, user.Email)
			}
		}

		// Update the account in the list
		accounts[i] = account
	}

	// Calculate total money market balance after interest calculations
	for _, account := range accounts {
		totalMoneyMarket += account.Deposit
	}

	// Save updated accounts and users if any interest was calculated
	if updated {
		SaveMoneyMarketAccounts(accounts)
		SaveUsers(users)
		// log.Printf("Interest calculated, accounts updated, and money market total is now %.2f.", totalMoneyMarket)
	}
}

// function to to convert time to read formart
func FormatCountdown(duration time.Duration) string {
	days := duration / (24 * time.Hour)
	hours := (duration % (24 * time.Hour)) / time.Hour
	minutes := (duration % time.Hour) / time.Minute

	return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
}

type MoneyMarketTrend struct {
	Timestamp   string  `json:"timestamp"`
	TotalAmount float64 `json:"total_amount"`
	UserCount   int     `json:"user_count"`
}

const MarketTrendsFile = "market_trends.json"

// LoadMarketTrends loads market trends from the JSON file
func LoadMarketTrends() []MoneyMarketTrend {
	data, err := os.ReadFile(MarketTrendsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []MoneyMarketTrend{}
		}
		log.Fatalf("Failed to load market trends: %v", err)
	}
	var trends []MoneyMarketTrend
	json.Unmarshal(data, &trends)
	return trends
}

// SaveMarketTrends saves market trends to the JSON file
func SaveMarketTrends(trends []MoneyMarketTrend) {
	data, err := json.MarshalIndent(trends, "", "  ")
	if err != nil {
		log.Fatalf("Failed to save market trends: %v", err)
	}
	os.WriteFile(MarketTrendsFile, data, 0o644)
}

// UpdateMarketTrends appends the latest trends to the file
func UpdateMarketTrends() {
	accounts := LoadMoneyMarketAccounts()
	totalAmount := 0.0
	userWallets := make(map[string]bool)

	for _, account := range accounts {
		totalAmount += account.Deposit
		userWallets[account.Wallet] = true
	}

	trends := LoadMarketTrends()
	newTrend := MoneyMarketTrend{
		Timestamp:   time.Now().Format(time.RFC3339),
		TotalAmount: totalAmount,
		UserCount:   len(userWallets),
	}
	trends = append(trends, newTrend)
	SaveMarketTrends(trends)
}

func ToJson(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("Error serializing data to JSON: %v", err)
		return "[]"
	}
	return string(data)
}
