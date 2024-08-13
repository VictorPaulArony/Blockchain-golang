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
	Adresses map[string]*Address
}
type Address struct {
	Privatekey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Balance    float64
}

// function to create a new wallet
func CreateWallet() *Wallet {
	return &Wallet{Adresses: make(map[string]*Address)}
}

// method to create a new Address for the wallet
func (w *Wallet) CreateAddress() string {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalln(privateKey)
	}

	publicKey := privateKey.PublicKey
	publicKeyBytes := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
	add := sha256.Sum256(publicKeyBytes)
	address := hex.EncodeToString(add[:])

	w.Adresses[address] = &Address{
		Privatekey: privateKey,
		PublicKey:  &publicKey,
		Balance:    0,
	}
	return address
}

// function to enable transfer of funds between adresses
func (w *Wallet) Transfer(to, from string, amonut float64) error {
	addFrom, exist := w.Adresses[from]
	if !exist {
		log.Fatalln(from)
	}

	addTo, exist := w.Adresses[to]
	if !exist {
		log.Fatalln(to)
	}
	addFrom.Balance += amonut
	addTo.Balance -= amonut

	return nil
}

// function to get the balance of the wallet addresses
func (w *Wallet) GetBalance(adress string) error {
	add, exist := w.Adresses[adress]
	if !exist {
		log.Fatalln(add)
	}
	return nil
}

// function main for the commputations
func main() {
	wallet := CreateWallet()

	addr1 := wallet.CreateAddress()
	addr2 := wallet.CreateAddress()

	fmt.Printf("THE FIRST ADDRESS: %s \n", addr1)
	fmt.Printf("THE SECOND ADDRESS: %s\n", addr2)
	wallet.Adresses[addr1].Balance = 136

	fmt.Printf("INITIAL BALANCE FOR FIRST ADDRESS: %f\n", wallet.Adresses[addr1].Balance)
	fmt.Printf("INITIAL BALANCS FOR SECOND ADDRESS: %f\n", wallet.Adresses[addr2].Balance)

	err := wallet.Transfer(addr1, addr2, 25)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("CURRENT BALANCE FOR FIRST ADDRESS: %f\n", wallet.Adresses[addr1].Balance)
	fmt.Printf("CURRENT BALANCS FOR SECOND ADDRESS: %f\n", wallet.Adresses[addr2].Balance)
}
