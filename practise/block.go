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

type Address struct {
	PublicKey  *ecdsa.PublicKey
	PrivateKey *ecdsa.PrivateKey
	Balance    float64
}

type Wallet struct {
	Addresses map[string]*Address
}

// function to create new wallet for the addresses
func CreateWallet() Wallet {
	return Wallet{Addresses: make(map[string]*Address)}
}

// function to create the address public and private key(s)
func (w *Wallet) CreateAddress() string {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalln("ERROR WHILE GENERATING PUBLICKEY", err)
	}
	publicKey := privateKey.PublicKey
	pubblicKeyByte := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
	addr := sha256.Sum256([]byte(pubblicKeyByte))
	addrStr := hex.EncodeToString(addr[:])

	w.Addresses[addrStr] = &Address{
		PrivateKey: privateKey,
		PublicKey:  &publicKey,
		Balance:    0.0,
	}

	return addrStr
}

// function to allow transfer of funds between addresses
func (w *Wallet) Transfer(from string, to string, amount float64) error {
	addrTo, exist := w.Addresses[to]
	if !exist {
		log.Fatalln(addrTo)
	}

	addrFrom, exist := w.Addresses[from]
	if !exist {
		log.Fatalln(addrFrom)
	}

	addrFrom.Balance -= amount
	addrTo.Balance += amount

	return nil
}

// function main to handle all the function as one
func main() {
	wallet := CreateWallet()

	addr1 := wallet.CreateAddress()
	addr2 := wallet.CreateAddress()

	wallet.Addresses[addr1].Balance = 120.0
	wallet.Addresses[addr2].Balance = 100.0

	fmt.Printf("Address one: %s\n", addr1)
	fmt.Printf("Address two: %s\n", addr2)
	fmt.Println()


	fmt.Printf("Initial Balance in Address One: %f\n", wallet.Addresses[addr1].Balance)
	fmt.Printf("Initial Balance in Address two: %f\n", wallet.Addresses[addr2].Balance)
	fmt.Println()

	err := wallet.Transfer(addr1, addr2, 25)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Current Balance in Address One: %f\n", wallet.Addresses[addr1].Balance)
	fmt.Printf("Current Balance in Address two: %f\n", wallet.Addresses[addr2].Balance)
}
