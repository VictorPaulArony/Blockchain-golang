package helpers

import (
	"encoding/json"
	"os"
)

// MoneyMarketDeposit represents a deposit in the money market.
type MoneyMarketDeposit struct {
	MoneyMarketTransaction
	Deposit      float64 `json:"deposit"` // Amount deposited
	Status       string  `json:"status"`  // Status of the deposit (e.g., "Active", "Closed")
	Email        string  `json:"email"`   // Email of the user making the deposit
	MaturityDate int64   `json:"maturity"`
}

// MoneyMarketTransaction is a common struct for shared attributes between deposits and loans.
type MoneyMarketTransaction struct {
	Name     string `json:"name"`
	Wallet   string `json:"wallet"`
	Joined   string `json:"joined"` // Number of months the member has been in the market
	Maturity int    `json:"maturity"`
}

// MoneyMarket struct for managing deposits.
type MoneyMarket struct {
	Total   float64                         `json:"total"`
	Members map[string][]MoneyMarketDeposit `json:"members"` // Map of wallet addresses to deposits
}

// LoadMoneyMarket retrieves money market data from the JSON file.
func LoadMoneyMarket() (MoneyMarket, error) {
	data, err := os.ReadFile(moneyMarketFile)
	if err != nil {
		if os.IsNotExist(err) {
			return MoneyMarket{Members: make(map[string][]MoneyMarketDeposit)}, nil
		}
		return MoneyMarket{}, err
	}

	var mmf MoneyMarket
	if err := json.Unmarshal(data, &mmf); err != nil {
		return MoneyMarket{}, err
	}
	return mmf, nil
}

// SaveMoneyMarket writes updated money market data to the JSON file.
func SaveMoneyMarket(mmf MoneyMarket) error {
	data, err := json.MarshalIndent(mmf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(moneyMarketFile, data, 0o644)
}

// AddMoneyMarketDeposit appends a new deposit to the money market.
func AddMoneyMarketDeposit(deposit MoneyMarketDeposit) error {
	mmf, err := LoadMoneyMarket()
	if err != nil {
		return err
	}

	// Append the deposit to the user's deposits in the money market.
	mmf.Members[deposit.Wallet] = append(mmf.Members[deposit.Wallet], deposit)
	mmf.Total += deposit.Deposit

	// Save the updated money market data.
	if err := SaveMoneyMarket(mmf); err != nil {
		return err
	}
	return nil
}
