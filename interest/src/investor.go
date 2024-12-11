package helpers

import (
	"encoding/json"
	"log"
	"os"
)

// inverstors struct to used in registration
type MoneyMarketInvestor struct {
	ID           string                          `json:"id"`
	Investor     string                          `json:"investor"`      // User who created the fund
	TotalAmount  float64                         `json:"total_amount"`  // Total amount in the fund
	InterestRate float64                         `json:"interest_rate"` // Interest rate for the fund
	Status       string                          `json:"status"`        // Active, Inactive
	Members      map[string][]MoneyMarketDeposit `json:"members"`       // List of users who joined
}

type MoneyMarketInvestorsAccounts struct {
	TotalFounds float64                          `json:"totalfunds"`
	Inverstors  map[string][]MoneyMarketInvestor `json:"inverstors"`
}

const (
	inverstorsFile = "inverstors.json"
)

// function to save the investors info
func SaveInverstors(mmIA MoneyMarketInvestorsAccounts) {
	data, err := json.MarshalIndent(mmIA, "", " ")
	if err != nil {
		log.Println("Error while Marshaling Inverstors info: ", err)
	}

	// save the data to the db
	err = os.WriteFile(inverstorsFile, data, 0o644)
	if err != nil {
		log.Println("Error while Saving investors data in db: ", err)
	}
}

// function ton load the inverstors data from the db
func LoadInvestorData() MoneyMarketInvestorsAccounts {
	data, err := os.ReadFile(inverstorsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return MoneyMarketInvestorsAccounts{Inverstors: make(map[string][]MoneyMarketInvestor)}
		}
	}

	var mmIA MoneyMarketInvestorsAccounts
	if err = json.Unmarshal(data, &mmIA); err != nil {
		log.Println("Error while unmarshaling the inverstors data from the db:", err)
	}
	return mmIA
}

// isEligibleToCreateMMF checks if the user meets the criteria to create an MMF.
func IsEligibleToCreateMMF(user *User) bool {
	totalTransactions := 0.0
	for _, txn := range user.Transactions {
		totalTransactions += txn.Amount
	}

	// Example criteria
	return totalTransactions > 10000 && user.Balance > 500 && len(user.Loans) > 0
}

// function to retrive a user using there emails
func GetUserByEmail(email string) *User {
	var wallet Wallet
	wallet.LoadData()
	var currentUser *User
	for _, user := range wallet.Users {
		if user.Email == email {
			currentUser = user
			break
		}
	}
	return currentUser
}
