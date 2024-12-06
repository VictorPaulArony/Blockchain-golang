package helpers

import (
	"encoding/json"
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

// MoneyMarketLoan struct for managing loans
type MoneyMarketLoan struct {
	Total   float64                  `json:"total"`
	Members map[string][]LoanRequest `json:"members"` // List of loan requests
}

const (
	mmlJsonFile = "mmlFile.json"
)

// function to save the mml to the db
func SaveMml(mml MoneyMarketLoan) {
	data, err := json.MarshalIndent(mml, "", " ")
	if err != nil {
		log.Printf("Error while marshalling Mml data: %v", err)
	}

	err = os.WriteFile(mmlJsonFile, data, 0o644)
	if err != nil {
		log.Printf("Error while Saving Mml data: %v", err)
	}
}

// function to load the mml
func Loadmml() MoneyMarketLoan {
	data, err := os.ReadFile(mmlJsonFile)
	if err != nil {
		if os.IsNotExist(err) {
			return MoneyMarketLoan{Members: make(map[string][]LoanRequest)}
		}
		log.Printf("Error while Reading the Mml data: %v", err)
		return MoneyMarketLoan{}
	}

	var mml MoneyMarketLoan

	err = json.Unmarshal(data, &mml)
	if err != nil {
		log.Printf("Error while Unmarshalling Mml data: %v", err)
		return MoneyMarketLoan{}
	}
	return mml
}

// AddMoneyMarketDeposit appends a new deposit to the money market.
func AddMoneyMarketLoan(loan LoanRequest) error {
	mml := Loadmml()

	mml.Members[loan.Wallet] = append(mml.Members[loan.Wallet], loan)
	mml.Total += loan.Amount

	SaveMml(mml)
	return nil
}
