package handlers

import (
	"net/http"
	"regexp"
	"strconv"
	"text/template"
	"time"

	blockchains "money-market/blockchain"
	helpers "money-market/utils"

	"github.com/google/uuid"
)

var blockchain blockchains.Blockchain

// Templates for the pages to be used in the program
var templates = template.Must(template.New("").Funcs(template.FuncMap{
	"add":        helpers.Add,
	"mul":        helpers.Mul,
	"formatDate": helpers.Format,
	"toJSON":     helpers.ToJson,
}).ParseFiles(
	"templates/index.html",
	"templates/dashboard.html",
	"templates/money_market.html",
	"templates/market_trends.html",
))

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

	if !isValidEmail(email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	users := helpers.LoadUsers()
	for _, user := range users {
		if user.Email == email {
			http.Error(w, "Email already registered", http.StatusBadRequest)
			return
		}
	}

	hashedPassword := helpers.GenerateHash(password)
	user := helpers.User{
		ID:       uuid.New().String(),
		Email:    email,
		Name:     name,
		Phone:    phone,
		Password: hashedPassword,
		Wallet:   helpers.GenerateWallet(email),
		JoinDate: time.Now().Format(time.RFC3339),
		Balance:  1000, // Default balance given to every new account
	}

	users = append(users, user)
	helpers.SaveUsers(users)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// IndexHandler serves the index.html page
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

// function to vallidate an email of the user as in the formart
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
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
	hashedPassword := helpers.GenerateHash(password)

	users := helpers.LoadUsers()
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

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Get the logged-in user's email from the cookie
	cookie, err := r.Cookie("user_email")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	email := cookie.Value

	// Fetch the current user's details
	users := helpers.LoadUsers()
	var currentUser helpers.User
	for _, user := range users {
		if user.Email == email {
			currentUser = user
			break
		}
	}

	// Load all transactions
	allTransactions := helpers.LoadTransactions()

	// Determine the filter based on the query parameter
	viewAll := r.URL.Query().Get("view") == "all"
	var filteredTransactions []helpers.Transaction

	if viewAll {
		// Show all transactions
		filteredTransactions = allTransactions
	} else {
		// Show only transactions related to the user's wallet address
		for _, tx := range allTransactions {
			if tx.Sender == currentUser.Wallet || tx.Receiver == currentUser.Wallet {
				filteredTransactions = append(filteredTransactions, tx)
			}
		}
	}

	// Limit the displayed transactions to the 10 most recent
	if len(filteredTransactions) > 10 {
		filteredTransactions = filteredTransactions[len(filteredTransactions)-10:]
	}

	// Load the blockchain
	blockchain := blockchains.LoadBlockchain()

	// Prepare the data for the dashboard template
	data := struct {
		User         helpers.User
		Blockchain   blockchains.Blockchain
		MempoolSize  int
		Transactions []helpers.Transaction
		ViewingAll   bool
	}{
		User:         currentUser,
		Blockchain:   blockchain,
		MempoolSize:  len(blockchain.Mempool),
		Transactions: filteredTransactions,
		ViewingAll:   viewAll,
	}

	// Render the dashboard template
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

	users := helpers.LoadUsers()
	var senderUser, receiverUser *helpers.User
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

	// Deduct and update balances
	senderUser.Balance -= amount
	receiverUser.Balance += amount
	helpers.SaveUsers(users)

	// Add transaction to the mempool
	transaction := helpers.Transaction{
		ID:        uuid.New().String(),
		Sender:    sender,
		Receiver:  receiver,
		Amount:    amount,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	blockchains.AddTransactionToMempool(transaction)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// function to handle the money market page
func MoneyMarketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Get the logged-in user's email from the cookie
		cookie, err := r.Cookie("user_email")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		email := cookie.Value

		// Load users and find the current user
		users := helpers.LoadUsers()
		var currentUser helpers.User
		for _, user := range users {
			if user.Email == email {
				currentUser = user
				break
			}
		}

		// Load all money market accounts
		allAccounts := helpers.LoadMoneyMarketAccounts()

		// Filter accounts for the current user
		var userAccounts []helpers.MoneyMarketAccount
		var fixedBalance, nonFixedBalance float64
		for i, account := range allAccounts {
			if account.Wallet == currentUser.Wallet {
				// Calculate maturity countdown for fixed accounts
				if account.AccountType == "fixed" {
					fixedEndDate, _ := time.Parse(time.RFC3339, account.FixedEndDate)
					remainingTime := fixedEndDate.Sub(time.Now())
					if remainingTime > 0 {
						allAccounts[i].FixedEndDate = helpers.FormatCountdown(remainingTime)
					} else {
						allAccounts[i].FixedEndDate = "Matured"
					}
				}

				// Update balances
				userAccounts = append(userAccounts, account)
				if account.AccountType == "fixed" {
					fixedBalance += account.Deposit
				} else if account.AccountType == "non-fixed" {
					nonFixedBalance += account.Deposit
				}
			}
		}

		// Prepare data for the template
		data := struct {
			User            helpers.User
			UserAccounts    []helpers.MoneyMarketAccount
			FixedBalance    float64
			NonFixedBalance float64
			AllAccounts     []helpers.MoneyMarketAccount
		}{
			User:            currentUser,
			UserAccounts:    userAccounts,
			FixedBalance:    fixedBalance,
			NonFixedBalance: nonFixedBalance,
			AllAccounts:     allAccounts,
		}

		// Render the money market template
		templates.ExecuteTemplate(w, "money_market.html", data)
		return
	}

	// Handle POST for joining the money market
	r.ParseForm()

	// Only allow the current user to join the market using their wallet
	cookie, err := r.Cookie("user_email")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	email := cookie.Value

	// Load users and find the current user
	users := helpers.LoadUsers()
	var currentUser *helpers.User
	for i := range users {
		if users[i].Email == email {
			currentUser = &users[i]
			break
		}
	}

	// Ensure user is registered
	if currentUser == nil {
		http.Error(w, "User not found. Only registered users can join the money market.", http.StatusUnauthorized)
		return
	}

	// Ensure the user is only using their wallet
	if r.FormValue("wallet") != currentUser.Wallet {
		http.Error(w, "Invalid wallet address. You can only use your own wallet to join the money market.", http.StatusBadRequest)
		return
	}

	// Validate deposit amount
	depositStr := r.FormValue("deposit")
	accountType := r.FormValue("account_type")
	deposit, err := strconv.ParseFloat(depositStr, 64)
	if err != nil || deposit < 100 { // Minimum deposit is 100
		http.Error(w, "Invalid deposit amount. Minimum is 100.", http.StatusBadRequest)
		return
	}

	// Ensure the deposit is less than or equal to the user's balance
	if deposit > currentUser.Balance {
		http.Error(w, "Insufficient balance for the deposit.", http.StatusBadRequest)
		return
	}

	// Deduct the deposit from the user's balance
	currentUser.Balance -= deposit
	helpers.SaveUsers(users) // Save updated users

	// Add the account to the money market
	account := helpers.MoneyMarketAccount{
		Wallet:      currentUser.Wallet,
		AccountType: accountType,
		Deposit:     deposit,
		JoinDate:    time.Now().Format(time.RFC3339),
	}
	helpers.AddMoneyMarketAccount(account)

	// Redirect back to the money market page
	http.Redirect(w, r, "/money-market", http.StatusSeeOther)
}
func MarketTrendsHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user's email from the cookie
	cookie, err := r.Cookie("user_email")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	email := cookie.Value

	// Fetch the current user's details
	users := helpers.LoadUsers()
	var currentUser helpers.User
	for _, user := range users {
		if user.Email == email {
			currentUser = user
			break
		}
	}

	// Load market trends data
	trends := helpers.LoadMarketTrends()

	// Prepare data for the template
	data := struct {
		User   helpers.User
		Trends []helpers.MoneyMarketTrend
	}{
		User:   currentUser,
		Trends: trends,
	}

	// Render the market trends template
	templates.ExecuteTemplate(w, "market_trends.html", data)
}