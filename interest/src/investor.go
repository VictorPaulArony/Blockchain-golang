package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// LoanRequest using the common struct with a status field
type LoanRequest struct {
	MoneyMarketTransaction
	ID        string  `json:"id"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"` // Paid, Approved, Late, Defaulted
	Requested string  `json:"requested"`
	Grade     string  `json:"grade"`
	DueDate   int64   `json:"due_date"` // Add a due date for the loan repayment
}

// MoneyMarketInvestor struct for investor details
type MoneyMarketInvestor struct {
	ID             string                          `json:"id"`
	Name           string                          `json:"name"`
	Email          string                          `json:"manager_email"`
	Phone          string                          `json:"phone"`
	Wallet         string                          `json:"wallet"`
	Amount         float64                         `json:"amount"`
	InvestmentType string                          `json:"investment_type"`
	Description    string                          `json:"description"`
	InterestRate   float64                         `json:"interest_rate"` // Interest rate for the fund
	Status         string                          `json:"status"`        // Active, Inactive
	LoanMembers    map[string][]LoanRequest        `json:"loan_members"`  // All loans managed by the investor
	MmfMembers     map[string][]MoneyMarketDeposit `json:"mmf_members"`
}

// MoneyMarketInvestorsAccounts struct to manage investor accounts
type MoneyMarketInvestorsAccounts struct {
	TotalFunds float64                         `json:"total_funds"`
	Investors  map[string]*MoneyMarketInvestor `json:"investors"`
}

const (
	investorsFile = "investors.json"
)

// SaveInvestors saves the investors information to a file
func SaveInvestors(mmIA MoneyMarketInvestorsAccounts) error {
	if mmIA.Investors == nil {
		mmIA.Investors = make(map[string]*MoneyMarketInvestor)
	}

	data, err := json.MarshalIndent(mmIA, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling investors info: %w", err)
	}

	if err = os.WriteFile(investorsFile, data, 0o644); err != nil {
		return fmt.Errorf("error saving investors data: %w", err)
	}

	return nil
}

// LoadInvestorData loads investors data from the file
func LoadInvestorData() MoneyMarketInvestorsAccounts {
	data, err := os.ReadFile(investorsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return MoneyMarketInvestorsAccounts{
				TotalFunds: 0,
				Investors:  make(map[string]*MoneyMarketInvestor),
			}
		}
		log.Printf("Error reading investors file: %v", err)
		return MoneyMarketInvestorsAccounts{
			Investors: make(map[string]*MoneyMarketInvestor),
		}
	}

	var mmIA MoneyMarketInvestorsAccounts
	if err = json.Unmarshal(data, &mmIA); err != nil {
		log.Printf("Error unmarshaling investors data: %v", err)
		return MoneyMarketInvestorsAccounts{
			Investors: make(map[string]*MoneyMarketInvestor),
		}
	}

	if mmIA.Investors == nil {
		mmIA.Investors = make(map[string]*MoneyMarketInvestor)
	}

	return mmIA
}

// AddLoanToInvestor adds a loan directly to the investor's record
func AddLoanToInvestor(investorEmail string, loan LoanRequest) error {
	mmIA := LoadInvestorData()

	investor, exists := mmIA.Investors[investorEmail]
	if !exists {
		return fmt.Errorf("investor with email %s not found", investorEmail)
	}

	if investor.LoanMembers == nil {
		investor.LoanMembers = make(map[string][]LoanRequest)
	}

	// Add the loan to the investor's loan members
	investor.LoanMembers[loan.Wallet] = append(investor.LoanMembers[loan.Wallet], loan)
	mmIA.TotalFunds += loan.Amount

	// Save updated data
	if err := SaveInvestors(mmIA); err != nil {
		return fmt.Errorf("error saving updated investor data: %w", err)
	}

	return nil
}

// GetLoansByInvestor retrieves all loans for a specific investor
func GetLoansByInvestor(investorEmail string) ([]LoanRequest, error) {
	mmIA := LoadInvestorData()

	investor, exists := mmIA.Investors[investorEmail]
	if !exists {
		return nil, fmt.Errorf("investor with email %s not found", investorEmail)
	}

	var allLoans []LoanRequest
	for _, loans := range investor.LoanMembers {
		allLoans = append(allLoans, loans...)
	}

	return allLoans, nil
}

// IsEligibleToCreateMMF checks if the user meets the criteria to create an MMF
func IsEligibleToCreateMMF(user *User) bool {
	// Defensive check for nil user
	if user == nil {
		return false
	}

	totalTransactions := 0.0
	for _, txn := range user.Transactions {
		totalTransactions += txn.Amount
	}

	// Example criteria: minimum total transactions, balance, and existing loans
	return totalTransactions > 1000 && user.Balance > 1000 && len(user.Loans) > 0
}

// GetUserByEmail retrieves a user by their email
func GetUserByEmail(email string) *User {
	var wallet Wallet
	wallet.LoadData()

	for _, user := range wallet.Users {
		if user.Email == email {
			return user
		}
	}

	return nil
}