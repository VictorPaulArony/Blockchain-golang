package handlers

import (
	"net/http"
	"strconv"
	"strings"

	helpers "interest/src"

	"github.com/google/uuid"
)

// InvestorsHandler handles investor-related operations
func InvestorsHandler(w http.ResponseWriter, r *http.Request) {
	// GET method handler
	if r.Method == http.MethodGet {
		cookie, err := r.Cookie("user_email")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		email := cookie.Value

		// Load investor data
		investorData := helpers.LoadInvestorData()
		var currentInvestor *helpers.MoneyMarketInvestor

		// Find current investor
		for _, investor := range investorData.Investors {
			if investor.Email == email {
				currentInvestor = investor
				break
			}
		}

		// Prepare data for template
		data := struct {
			User            *helpers.User
			CurrentInvestor *helpers.MoneyMarketInvestor
			InvestorData    helpers.MoneyMarketInvestorsAccounts
			Investor        *helpers.User
		}{
			User:            helpers.GetUserByEmail(email),
			CurrentInvestor: currentInvestor,
			InvestorData:    investorData,
			Investor:        helpers.GetUserByEmail(email),
		}
		

		// Execute template
		tmpl := setupTemplates()
		if err := tmpl.ExecuteTemplate(w, "investors.html", data); err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
		}
		return
	}

	// POST method handler
	if r.Method == http.MethodPost {
		// Check user authentication
		cookie, err := r.Cookie("user_email")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		loggedInEmail := cookie.Value

		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		// Parse and validate input
		name := r.Form.Get("name")
		investorEmail := r.Form.Get("email")
		phone := r.Form.Get("phone")
		wallet := r.Form.Get("wallet")
		amountStr := r.Form.Get("amount")
		investmentType := r.Form.Get("investment-type")
		description := r.Form.Get("description")

		// Validate that the logged-in user is registering only themselves
		if strings.TrimSpace(investorEmail) != loggedInEmail {
			http.Error(w, "You can only register yourself as an investor", http.StatusForbidden)
			return
		}

		// Load investor data
		investorData := helpers.LoadInvestorData()

		// Get user details
		investorUserDetails := helpers.GetUserByEmail(loggedInEmail)
		if investorUserDetails == nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		if _, exists := investorData.Investors[name]; exists {
			http.Error(w, "An investor is already registered with that name", http.StatusConflict)
			return
		}

		// Check investor eligibility
		if !helpers.IsEligibleToCreateMMF(investorUserDetails) {
			http.Error(w, "You are not eligible to become an investor", http.StatusForbidden)
			return
		}

		// Check if investor already exists
		if _, exists := investorData.Investors[loggedInEmail]; exists {
			http.Error(w, "You are already registered as an investor", http.StatusConflict)
			return
		}

		// Validate investment amount
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil || amount < 1000 {
			http.Error(w, "Invalid investment amount. Minimum is $1000", http.StatusBadRequest)
			return
		}

		// Ensure Investors map is initialized
		if investorData.Investors == nil {
			investorData.Investors = make(map[string]*helpers.MoneyMarketInvestor)
		}

		// Create new investor account
		newInvestor := &helpers.MoneyMarketInvestor{
			ID:             uuid.New().String(),
			Name:           investorUserDetails.Name, // Use logged-in user's name
			Email:          loggedInEmail,            // Use logged-in user's email
			Phone:          phone,
			Wallet:         wallet,
			Amount:         amount,
			InvestmentType: investmentType,
			Description:    description,
			InterestRate:   calculateInterestRate(amount),
			Status:         "Active",
			LoanMembers:    make(map[string][]helpers.LoanRequest),
			MmfMembers:     make(map[string][]helpers.MoneyMarketDeposit),
		}

		// Add new investor to the data
		investorData.Investors[loggedInEmail] = newInvestor
		investorData.TotalFunds += newInvestor.Amount

		// update user registered as investors
		var UserWallet helpers.Wallet
		UserWallet.LoadData()

		var userR *helpers.User

		for _, user := range UserWallet.Users {
			if user.Wallet == wallet {
				userR = user
			}
		}

		if userR.Balance < amount {
			http.Error(w, "Insufficient balance", http.StatusBadRequest)
			return
		}

		// user balance is reduced after joining mmf as investor
		userR.Balance -= amount
		UserWallet.SaveData()

		// Save investor data
		if err := helpers.SaveInvestors(investorData); err != nil {
			http.Error(w, "Failed to save investor data", http.StatusInternalServerError)
			return
		}

		// Redirect to investors page
		http.Redirect(w, r, "/investors", http.StatusSeeOther)
		return
	}
}

// calculateInterestRate determines interest rate based on investment amount
func calculateInterestRate(amount float64) float64 {
	switch {
	case amount < 5000:
		return 0.5
	case amount < 25000:
		return 0.75
	case amount < 100000:
		return 1.0
	default:
		return 1.5
	}
}
