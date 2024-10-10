/*
package main

import (
	"errors"
	"fmt"
)

// Token represents the smart contract state for our simple token.
type Token struct {
	name     string
	symbol   string
	decimals uint8
	owner    string
	balances map[string]uint
}

// NewToken creates a new token with a name, symbol, decimals, and initializes the contract.
func NewToken(name, symbol string, decimals uint8, owner string) *Token {
	return &Token{
		name:     name,
		symbol:   symbol,
		decimals: decimals,
		owner:    owner,
		balances: make(map[string]uint),
	}
}

// Mint allows the owner to create new tokens and assign them to an account.
func (t *Token) Mint(to string, amount uint) error {
	if t.owner != to {
		return errors.New("only the owner can mint tokens")
	}
	t.balances[to] += amount
	fmt.Printf("Minted %d tokens to %s\n", amount, to)
	return nil
}

// Transfer allows a user to transfer tokens from one account to another.
func (t *Token) Transfer(from, to string, amount uint) error {
	if t.balances[from] < amount {
		return errors.New("insufficient balance")
	}
	t.balances[from] -= amount
	t.balances[to] += amount
	fmt.Printf("Transferred %d tokens from %s to %s\n", amount, from, to)
	return nil
}

// BalanceOf returns the balance of the specified account.
func (t *Token) BalanceOf(account string) uint {
	return t.balances[account]
}

// TransferOwnership allows the owner to transfer contract ownership to another account.
func (t *Token) TransferOwnership(newOwner string) error {
	if newOwner == "" {
		return errors.New("new owner cannot be empty")
	}
	fmt.Printf("Ownership transferred from %s to %s\n", t.owner, newOwner)
	t.owner = newOwner
	return nil
}

// GetOwner returns the current owner of the token contract.
func (t *Token) GetOwner() string {
	return t.owner
}

// DisplayContractInfo displays the token contract details.
func (t *Token) DisplayContractInfo() {
	fmt.Printf("Token Name: %s\nSymbol: %s\nDecimals: %d\nOwner: %s\n", t.name, t.symbol, t.decimals, t.owner)
}

func main() {
	// Create a new token with an initial owner.
	token := NewToken("MyToken", "MTK", 18, "alice")

	// Display contract information.
	token.DisplayContractInfo()

	// Mint some tokens to the owner.
	err := token.Mint("alice", 1000)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Check balance of alice.
	fmt.Println("Alice's balance:", token.BalanceOf("alice"))

	// Transfer tokens from Alice to Bob.
	err = token.Transfer("alice", "bob", 500)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Check balances of Alice and Bob.
	fmt.Println("Alice's balance:", token.BalanceOf("alice"))
	fmt.Println("Bob's balance:", token.BalanceOf("bob"))

	// Transfer ownership to Bob.
	err = token.TransferOwnership("bob")
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Try to mint tokens from Bob's account (since he is the new owner).
	err = token.Mint("bob", 500)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Check Bob's balance again.
	fmt.Println("Bob's balance after minting:", token.BalanceOf("bob"))

	// Display the final state of the contract.
	token.DisplayContractInfo()
}

*/
