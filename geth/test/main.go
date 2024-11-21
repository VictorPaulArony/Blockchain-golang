package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Transaction struct {
	Type      string  `json:"type"`
	Amount    float64 `json:"amount"`
	Timestamp string  `json:"timestamp"`
}

type User struct {
	Username     string        `json:"username"`
	Email        string        `json:"email"`
	Phone        string        `json:"phone"`
	Address      string        `json:"address"`
	PrivateKey   string        `json:"privateKey"`
	Balance      float64       `json:"balance"`
	Transactions []Transaction `json:"transactions"`
}

type Wallet struct {
	Users map[string]*User
	mu    sync.Mutex
}

const dataFile = "users.json"

// CreateWallet initializes a new Wallet instance
func CreateWallet() *Wallet {
	w := &Wallet{Users: make(map[string]*User)}
	w.LoadData()
	return w
}

// SaveData saves the wallet details to a JSON file
func (w *Wallet) SaveData() error {


	data, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		return fmt.Errorf("error while marshalling data: %v", err)
	}

	err = os.WriteFile(dataFile, data, 0o644)
	if err != nil {
		return fmt.Errorf("error while writing to JSON: %v", err)
	}
	return nil
}

// LoadData loads the data from the JSON file
func (w *Wallet) LoadData() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return
	}

	data, err := os.ReadFile(dataFile)
	if err != nil {
		log.Printf("Error while reading user file: %v", err)
		return
	}

	if err = json.Unmarshal(data, &w); err != nil {
		log.Printf("Error while unmarshalling the user data: %v", err)
	}
}

// Helper function to encode private key to hex
func encodePrivateKeyToHex(privateKey *ecdsa.PrivateKey) string {
	return hex.EncodeToString(privateKey.D.Bytes())
}

// CreateAddress generates a new address for the user
func (w *Wallet) CreateAddress(username, email, phone string) (string, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, exists := w.Users[email]; exists {
		return "", fmt.Errorf("user with email %s already exists", email)
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", fmt.Errorf("error generating private key: %v", err)
	}

	publicKey := privateKey.PublicKey
	publicKeyBytes := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
	hashAddr := sha256.Sum256(publicKeyBytes)
	address := hex.EncodeToString(hashAddr[:])

	privateKeyHex := encodePrivateKeyToHex(privateKey)

	user := &User{
		Username:     username,
		Email:        email,
		Phone:        phone,
		Address:      address,
		PrivateKey:   privateKeyHex,
		Balance:      100.0,
		Transactions: []Transaction{},
	}

	w.Users[email] = user
	return address, w.SaveData()
}

// RegisterHandler handles user registration
func (w *Wallet) RegisterHandler(wrt http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(wrt, r, "index.html")
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		phone := r.FormValue("phone")

		address, err := w.CreateAddress(username, email, phone)
		if err != nil {
			http.Error(wrt, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Fprintf(wrt, "User registered successfully! Wallet Address: %s", address)
		return
	}
	http.Error(wrt, "Invalid request method", http.StatusMethodNotAllowed)
}

// TransactionHandler handles adding a transaction to a user's account
func (w *Wallet) TransactionHandler(wrt http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		txType := r.FormValue("type") // "deposit" or "withdrawal"
		amountStr := r.FormValue("amount")

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			http.Error(wrt, "Invalid amount", http.StatusBadRequest)
			return
		}

		w.mu.Lock()
		user, exists := w.Users[email]
		if !exists {
			w.mu.Unlock()
			http.Error(wrt, "User not found", http.StatusBadRequest)
			return
		}

		// Adjust the balance based on transaction type
		if txType == "withdrawal" && user.Balance < amount {
			w.mu.Unlock()
			http.Error(wrt, "Insufficient funds", http.StatusBadRequest)
			return
		}

		if txType == "deposit" {
			user.Balance += amount
		} else if txType == "withdrawal" {
			user.Balance -= amount
		} else {
			w.mu.Unlock()
			http.Error(wrt, "Invalid transaction type", http.StatusBadRequest)
			return
		}

		// Add transaction record
		transaction := Transaction{
			Type:      txType,
			Amount:    amount,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		user.Transactions = append(user.Transactions, transaction)

		w.mu.Unlock()
		w.SaveData()

		fmt.Fprintf(wrt, "Transaction successful! New balance for %s: %.2f", email, user.Balance)
		return
	}
	http.Error(wrt, "Invalid request method", http.StatusMethodNotAllowed)
}

func main() {
	wallet := CreateWallet()

	http.HandleFunc("/", wallet.RegisterHandler)
	http.HandleFunc("/register", wallet.RegisterHandler)
	http.HandleFunc("/transaction", wallet.TransactionHandler)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
