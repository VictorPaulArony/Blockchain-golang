package handlers

import (
	"net/http"
	"strconv"
	"time"

	helpers "interest/src"

	"github.com/google/uuid"
)

const (
	transactionThreshold = 500.0  // Example threshold
	baseLoanAmount       = 5000.0 // Minimum loan amount
	maxLoanMultiplier    = 10.0   // Max loan factor for excellent credit
)

// LoanHandler processes loan requests.
func LoanHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve user session from the cookie
	cookie, err := r.Cookie("user_email")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	email := cookie.Value

	// Load wallet data
	var wallet helpers.Wallet
	wallet.LoadData()

	var currentUser *helpers.User
	for _, user := range wallet.Users {
		if user.Email == email {
			currentUser = user
			break
		}
	}

	if currentUser == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check loan eligibility of the user
	now := time.Now()
	joinDate, _ := time.Parse("2006-01-02", currentUser.JoinDate)
	membershipDuration := now.Sub(joinDate).Hours() / (24 * 30) // Membership duration in months

	if membershipDuration < 6 {
		http.Error(w, "You must be a member for at least 6 months to apply for a loan.", http.StatusForbidden)
		return
	}

	totalTransactions := 0.0
	for _, txn := range currentUser.Transactions {
		totalTransactions += txn.Amount
	}
	totalTransactions /= 2.0

	if totalTransactions < transactionThreshold {
		http.Error(w, "You must have transacted more than a certain amount to be eligible.", http.StatusForbidden)
		return
	}

	// Calculate credit score and loan grading
	creditScore := calculateCreditScore(currentUser.Loans)
	var loanGrade string
	var loanMultiplier float64

	switch {
	case creditScore >= 90:
		loanGrade = "A"
		loanMultiplier = maxLoanMultiplier
	case creditScore >= 75:
		loanGrade = "B"
		loanMultiplier = 7.5
	case creditScore >= 50:
		loanGrade = "C"
		loanMultiplier = 5.0
	default:
		loanGrade = "D"
		loanMultiplier = 2.5
	}

	// Calculate maximum loan amount
	maxLoan := baseLoanAmount + (totalTransactions * loanMultiplier / 100)

	// Load investors data
	investors := helpers.LoadInvestorData()

	if r.Method == http.MethodGet {
		// Prepare data for rendering
		data := struct {
			User           helpers.User
			LoanGrade      string
			MaxLoanAmount  float64
			CreditScore    float64
			MembershipTime float64
			Investors      map[string]*helpers.MoneyMarketInvestor
		}{
			User:           *currentUser,
			LoanGrade:      loanGrade,
			MaxLoanAmount:  maxLoan,
			CreditScore:    creditScore,
			MembershipTime: membershipDuration,
			Investors:      investors.Investors,
		}

		// Render loan request form
		tmpl := setupTemplates()
		if err := tmpl.ExecuteTemplate(w, "loan_request.html", data); err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
		}
		return
	}

	// Handle POST: Loan request submission
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing loan request form", http.StatusBadRequest)
		return
	}

	investorEmail := r.FormValue("investor")
	requestedLoan, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil || requestedLoan <= 0 || requestedLoan > maxLoan {
		http.Error(w, "Invalid loan amount.", http.StatusBadRequest)
		return
	}

	// Get the selected investor
	investor, exists := investors.Investors[investorEmail]
	if !exists {
		http.Error(w, "Selected investor not found", http.StatusBadRequest)
		return
	}

	// Check if investor has sufficient funds
	if investor.Amount < requestedLoan {
		http.Error(w, "Insufficient funds from selected investor", http.StatusBadRequest)
		return
	}

	interestRate := investor.InterestRate

	// Calculate total amount to be repaid
	expectedToPay := requestedLoan + (requestedLoan * interestRate)

	// Record the loan request
	newLoan := helpers.LoanRequest{
		MoneyMarketTransaction: helpers.MoneyMarketTransaction{
			Name:     currentUser.Name,
			Wallet:   currentUser.Wallet,
			Joined:   currentUser.JoinDate,
			Maturity: 1,
		},
		ID:        uuid.New().String(),
		Amount:    expectedToPay,
		Status:    "Approved",
		Requested: time.Now().Format("2006-01-02"),
		Grade:     loanGrade,
		DueDate:   time.Now().Add(30 * 24 * time.Hour).Unix(), // Set due date to 30 days from now
	}

	// Deduct the investor's amount after borrowing
	investor.Amount -= requestedLoan
	// investor.LoanMembers[newLoan.Wallet] = newLoan

	// Update investors data
	investors.Investors[investorEmail] = investor
	helpers.SaveInvestors(investors)

	// Add the loan to the user's loans
	currentUser.Loans = append(currentUser.Loans, newLoan)

	// Update the user's balance
	currentUser.Balance += requestedLoan

	// Save wallet data
	if err := wallet.SaveData(); err != nil {
		http.Error(w, "Error saving loan data", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/loan", http.StatusSeeOther)
}

// calculateCreditScore evaluates the credit score based on past loans.
func calculateCreditScore(loans []helpers.LoanRequest) float64 {
	score := 100.0
	for _, loan := range loans {
		if loan.Status == "Defaulted" {
			score -= 30.0
		} else if loan.Status == "Late" {
			score -= 10.0
		}
	}
	if score < 0 {
		return 0
	}
	return score
}
