package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
)

type Wallet struct {
	Address map[string]*Adrress
}

type Adrress struct {
	PublicKey  *ecdsa.PublicKey
	PrivateKey *ecdsa.PrivateKey
	Balance    float64
}

// function to create a wallet
func CreateWallet() Wallet {
	return Wallet{Address: make(map[string]*Adrress)}
}

// function to create the Address for the wallet
func (w *Wallet) CreateAddress() string {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalln("Error Generating Key", err)
		return ""
	}

	publicKey := privateKey.PublicKey
	publicKeyBytes := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
	add := sha256.Sum256(publicKeyBytes)
	address := hex.EncodeToString(add[:])

	w.Address[address] = &Adrress{
		PublicKey:  &publicKey,
		PrivateKey: privateKey,
		Balance:    0,
	}
	return address
}

// function to create Transfer of funds
func (w *Wallet) Transfer(from, to string, amount float64) {
	sender, exist := w.Address[from]
	if !exist {
		log.Fatalln("The Address does not exist: ", exist)
	}

	reciever, exist := w.Address[to]
	if !exist {
		log.Fatalln("The Address does not exist: ", exist)
	}

	sender.Balance -= amount
	reciever.Balance += amount
}

// function to get the balance of in the Address
func (w *Wallet) GetBalance(address string) {
	addstr, exist := w.Address[address]
	if !exist {
		log.Fatalln("The Address does not exist: ", addstr)
	}
}

//function main 
func main(){
	wallet := CreateWallet()


	address1 := wallet.CreateAddress()
	address2 := wallet.CreateAddress()


	wallet.Address[address1].Balance = 120
	wallet.Address[address2].Balance = 300

	fmt.Printf("The Adress1 Initial Blance: %f\n", wallet.Address[address1].Balance)
	fmt.Printf("The Address2 Initial Balance:  %f\n", wallet.Address[address2].Balance)

	wallet.Transfer(address1, address2, 110.0)

	fmt.Printf("Address 1 Balance: %f\n", wallet.Address[address1].Balance)
	fmt.Printf("Address 2 Balance: %f\n",wallet.Address[address2].Balance)

}