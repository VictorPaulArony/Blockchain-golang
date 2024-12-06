package helpers

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           string               `json:"id"`
	Email        string               `json:"email"`
	Name         string               `json:"name"`
	Phone        string               `json:"phone"`
	Address      string               `json:"address"`
	PrivateKey   string               `json:"privateKey"`
	Wallet       string               `json:"wallet"`
	JoinDate     string               `json:"join_date"`
	Balance      float64              `json:"balance"`
	Password     string               `json:"password"`
	Transactions []Transaction        `json:"transactions"`
	Loans        []LoanRequest        `jsons:"loans"`
	Mmfs         []MoneyMarketDeposit `json:"mmfs"`
}
type Wallet struct {
	Users map[string]*User
	mu    sync.Mutex
}

type Transaction struct {
	ID        string  `json:"id"`
	Sender    string  `json:"sender"`
	Receiver  string  `json:"receiver"`
	Amount    float64 `json:"amount"`
	Timestamp string  `json:"timestamp"`
}



// constant variables for file names
const (
	userJsonFile    = "users.json"
	loanJsonFile    = "loan.json"
	moneyMarketFile = "moneyMarketFile.json"
)

// SaveData saves the wallet details to a JSON db
func (w *Wallet) SaveData() error {
	data, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		return fmt.Errorf("error while marshalling data: %v", err)
	}

	err = os.WriteFile(userJsonFile, data, 0o644)
	if err != nil {
		return fmt.Errorf("error while writing to JSON: %v", err)
	}
	return nil
}

// function to load the user data fron the json db file
func (w *Wallet) LoadData() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, err := os.Stat(userJsonFile); os.IsNotExist(err) {
		return
	}

	data, err := os.ReadFile(userJsonFile)
	if err != nil {
		log.Printf("Error while reading user file: %v", err)
		return
	}

	if err = json.Unmarshal(data, &w); err != nil {
		log.Printf("Error while unmarshalling the user data: %v", err)
	}
}

func (wallet *Wallet) CreateAddress(email, name, phone, password string) {
	wallet.mu.Lock()
	defer wallet.mu.Unlock()

	if wallet.Users == nil {
		wallet.Users = make(map[string]*User)
	}

	if _, exists := wallet.Users[email]; exists {
		log.Printf("user with email %s already exists", email)
		return
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Printf("error generating private key: %v", err)
		return
	}

	publicKey := privateKey.PublicKey
	publicKeyBytes := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
	hashAddr := sha256.Sum256(publicKeyBytes)
	address := hex.EncodeToString(hashAddr[:])

	hashedPassword := GenerateHash(password)
	privateKeyHex := encodePrivateKeyToHex(privateKey)

	user := &User{
		ID:           uuid.New().String(),
		Email:        email,
		Name:         name,
		Phone:        phone,
		Address:      address,
		PrivateKey:   privateKeyHex,
		Wallet:       GenerateWallet(email),
		JoinDate:     time.Now().Format("2006-01-02 15:04:05"),
		Balance:      1000, // Default balance given to every new account
		Password:     hashedPassword,
		Transactions: []Transaction{},
		Loans:        []LoanRequest{},
		Mmfs:         []MoneyMarketDeposit{},
	}
	wallet.Users[email] = user
	wallet.SaveData()
}

// Helper function to encode private key to hex
func encodePrivateKeyToHex(privateKey *ecdsa.PrivateKey) string {
	return hex.EncodeToString(privateKey.D.Bytes())
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
