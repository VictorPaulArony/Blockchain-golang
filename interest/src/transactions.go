package helpers

import (
	"encoding/json"
	"log"
	"os"
)

// SaveTransaction saves a new transaction to the JSON file
func SaveTransaction(tx Transaction) error {
	transactions := LoadTransactions()
	transactions = append(transactions, tx) 

	data, err := json.MarshalIndent(transactions, "", "  ")
	if err != nil {
		return err 
	}

	err = os.WriteFile("transactions.json", data, 0o644)
	if err != nil {
		return err 
	}

	return nil 
}

// LoadTransactions loads all transactions from a JSON file
func LoadTransactions() []Transaction {
	var transactions []Transaction

	data, err := os.ReadFile("transactions.json")
	if err != nil {
		log.Printf("Error reading transactions file: %v", err)
		return transactions 
	}

	if err := json.Unmarshal(data, &transactions); err != nil {
		log.Printf("Error unmarshalling transactions data: %v", err)
		return transactions
	}

	return transactions
}
