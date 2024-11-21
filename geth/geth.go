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
	"sync"
)

type User struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	PrivateKey string `json:"privatekey"`
}

type Wallet struct {
	User map[string]*User
	mu   sync.Mutex
}

// CreateWallet initializes a new Wallet instance
func CreateWallet() *Wallet {
	return &Wallet{User: make(map[string]*User)}
}

// SaveInfo saves the wallet details to a JSON file
func (w *Wallet) SaveInfo() error {
	// w.mu.Lock()
	// defer w.mu.Unlock()

	data, err := json.MarshalIndent(w.User, "", " ")
	if err != nil {
		log.Println("Error while marshalling the data:", err)
		return nil
	}
	err = os.WriteFile("users.json", data, 0o644)
	if err != nil {
		log.Println("Error while writing to JSON:", err)
	}
	return nil
}

// LoadWallet loads the data from the JSON file
func (w *Wallet) LoadWallet() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, err := os.Stat("users.json"); os.IsNotExist(err) {
		return
	}

	data, err := os.ReadFile("users.json")
	if err != nil {
		log.Println("Error while reading user file:", err)
		return
	}
	if err = json.Unmarshal(data, &w.User); err != nil {
		log.Println("Error while unmarshaling the user data:", err)
	}
}

// Helper function to encode private key
func encodePrivateKeyToHex(privateKey *ecdsa.PrivateKey) string {
	return hex.EncodeToString(privateKey.D.Bytes())
}

// CreateAddress generates a new address for the user
func (w *Wallet) CreateAddress(username, email, phone string) (string, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, exists := w.User[email]; exists {
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
		Username:   username,
		Email:      email,
		Phone:      phone,
		Address:    address,
		PrivateKey: privateKeyHex,
	}

	w.User[email] = user
	err = w.SaveInfo()
	if err != nil {
		log.Println("Error while saving: ", err)
	}

	return address, nil
}

// RegistrationHandler handles the user registration request
func (w *Wallet) RegistrationHandler(wrt http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(wrt, r, "index.html")
		return
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(wrt, "Unable to parse form", http.StatusBadRequest)
			return
		}
		username := r.FormValue("username")
		email := r.FormValue("email")
		phone := r.FormValue("phone")

		address, err := w.CreateAddress(username, email, phone)
		if err != nil {
			http.Error(wrt, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Fprintf(wrt, "User registered successfully! Address: %s", address)
		return
	}
	http.Error(wrt, "Invalid request method", http.StatusMethodNotAllowed)
}

func main() {
	wallet := CreateWallet()
	wallet.LoadWallet()

	http.HandleFunc("/", wallet.RegistrationHandler)

	fmt.Println("Server is running on port: http://localhost:1234")
	log.Fatal(http.ListenAndServe(":1234", nil))
}
