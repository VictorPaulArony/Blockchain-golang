package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Account struct {
	Balance float64
}

func main() {
	account := &Account{Balance: 120.00}

	input := bufio.NewScanner(os.Stdin)

	for {

		fmt.Println("\n ------- Welcome to your account-----")
		fmt.Println("1. Check Balance")
		fmt.Println("2. Withdraw cash")
		fmt.Println("3. Deposite Cash")
		fmt.Println("4. Exit")
		fmt.Println("choose an option")

		input.Scan()

		options := strings.TrimSpace(input.Text())

		switch options {
		case "1":
			account.CheckBalance()
		case "2":
			fmt.Println("Enter amount to withdraw")
			input.Scan()
			amount, _ := strconv.ParseFloat(strings.TrimSpace(input.Text()), 64)
			account.Withdraw(amount)
		case "3":
			fmt.Println("Enter amaount to Deposite")
			input.Scan()
			amount, _ := strconv.ParseFloat(strings.TrimSpace(input.Text()), 64)
			account.Deposite(amount)
		case "4":
			fmt.Println("Thank you, come back again")
			return
		default:
			fmt.Println("Invalid option, Try again")
		}

	}
}

// function main using fmt.Scanln functionality
// func main() {
// 	account := &Account{Balance: 120}

// 	for {
// 		fmt.Println("\n ------- Welcome to your account-----")
// 		fmt.Println("1. Check Balance")
// 		fmt.Println("2. Withdraw cash")
// 		fmt.Println("3. Deposite Cash")
// 		fmt.Println("4. Exit")
// 		fmt.Println("choose an option")

// 		var option string

// 		fmt.Scanln(&option)
// 		switch option {
// 		case "1":
// 			fmt.Println()
// 			account.CheckBalance()
// 		case "2":
// 			fmt.Println("Enter amount to withdraw")
// 			var amount string
// 			fmt.Scanln(&amount)
// 			am, _ := strconv.ParseFloat(amount, 64)
// 			account.Withdraw(am)
// 		case "3":
// 			fmt.Println("Enter amaount to Deposite")
// 			var amount string
// 			fmt.Scanln(&amount)
// 			am, _ := strconv.ParseFloat(amount, 64)
// 			account.Deposite(am)
// 		case "4":
// 			fmt.Println("Thank you, come back again")
// 			return
// 		default:
// 			fmt.Println("Invalid option, Try again")

// 		}
// 	}
// }

// function to check the balance in the account
func (acc *Account) CheckBalance() {
	fmt.Printf("Your Balance: %v\n", acc.Balance)
}

// function to withdraw amount from the account
func (acc *Account) Withdraw(amount float64) {
	if amount > acc.Balance {
		fmt.Printf("Insufficient account balance, your account balance is %v\n", acc.Balance)
		return
	}
	acc.Balance -= amount
	fmt.Printf("Successfully withdrew: %.2f\n", amount)
}

// function to Deposite the amount in the account
func (acc *Account) Deposite(amount float64) {
	acc.Balance += amount
	fmt.Printf("Successfully Deposited: %.2f\n", amount)
}
