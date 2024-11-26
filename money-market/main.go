package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
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

// File paths
const (
	UserFile        = "users.json"
	TransactionFile = "transactions.json"
)

// Templates
var templates = template.Must(template.ParseFiles(
	"templates/index.html",     // Registration and login
	"templates/dashboard.html", // User dashboard
))

// HashPassword hashes the password using SHA-256
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
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

// IndexHandler serves the index.html page
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	email := r.FormValue("email")
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	if password != confirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	users := LoadUsers()
	for _, user := range users {
		if user.Email == email {
			http.Error(w, "Email already registered", http.StatusBadRequest)
			return
		}
	}

	hashedPassword := HashPassword(password)
	user := User{
		ID:       uuid.New().String(),
		Email:    email,
		Name:     name,
		Phone:    phone,
		Password: hashedPassword,
		Wallet:   GenerateWallet(email),
		JoinDate: time.Now().Format(time.RFC3339),
		Balance:  1000, // Default balance
	}

	users = append(users, user)
	SaveUsers(users)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")
	hashedPassword := HashPassword(password)

	users := LoadUsers()
	for _, user := range users {
		if user.Email == email && user.Password == hashedPassword {
			// Set user session using a cookie
			http.SetCookie(w, &http.Cookie{
				Name:    "user_email",
				Value:   email,
				Expires: time.Now().Add(24 * time.Hour),
			})

			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
	}

	http.Error(w, "Invalid email or password", http.StatusUnauthorized)
}

// DashboardHandler displays the user dashboard
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Get the logged-in user's email from the session (cookie)
	cookie, err := r.Cookie("user_email")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	email := cookie.Value

	users := LoadUsers()
	var currentUser User
	for _, user := range users {
		if user.Email == email {
			currentUser = user
			break
		}
	}

	transactions := LoadTransactions()
	data := struct {
		User         User
		Transactions []Transaction
	}{
		User:         currentUser,
		Transactions: transactions,
	}

	templates.ExecuteTemplate(w, "dashboard.html", data)
}

// TransactionHandler processes transactions
func TransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	sender := r.FormValue("sender_wallet")
	receiver := r.FormValue("receiver_wallet")
	amountStr := r.FormValue("amount")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	users := LoadUsers()
	var senderUser, receiverUser *User
	for i := range users {
		if users[i].Wallet == sender {
			senderUser = &users[i]
		} else if users[i].Wallet == receiver {
			receiverUser = &users[i]
		}
	}

	if senderUser == nil || receiverUser == nil {
		http.Error(w, "Invalid wallet address", http.StatusBadRequest)
		return
	}
	if senderUser.Balance < amount {
		http.Error(w, "Insufficient balance", http.StatusBadRequest)
		return
	}

	// Update balances
	senderUser.Balance -= amount
	receiverUser.Balance += amount
	SaveUsers(users)

	// Record transaction
	transaction := Transaction{
		ID:        uuid.New().String(),
		Sender:    sender,
		Receiver:  receiver,
		Amount:    amount,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	transactions := LoadTransactions()
	transactions = append(transactions, transaction)
	SaveTransactions(transactions)

	// Redirect back to the dashboard
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/dashboard", DashboardHandler)
	http.HandleFunc("/transactions", TransactionHandler)

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
