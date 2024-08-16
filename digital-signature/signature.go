package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"
)

type Transaction struct {
	Sender    string
	Receiver  string
	Amount    float64
	TimeStamp string
	Signature string
}

type Block struct {
	ID          int
	Transaction []Transaction
	TimeStamp   string
	PrevHash    string
	Hash        string
}
type Blockchain struct {
	Blocks []Block
}

type Address struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Balance    float64
}

type Wallet struct {
	Addresses map[string]*Address
}

var transaction []Transaction

// function to create a hash function for the blockchain
func (b *Block) CreateHash() string {
	res := strconv.Itoa(b.ID) + b.TimeStamp + b.PrevHash + b.Hash

	for _, tx := range b.Transaction {
		res += tx.Receiver + tx.Sender + tx.TimeStamp + fmt.Sprintf("%f", tx.Amount)
	}
	h := sha256.New()
	h.Write([]byte(res))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// function to create the newBlock
func NewBlock(id int, transaction []Transaction, prevHash string) Block {
	newBlock := Block{
		ID:          id,
		Transaction: transaction,
		PrevHash:    prevHash,
		TimeStamp:   time.Now().String(),
		Hash:        "",
	}
	newBlock.Hash = newBlock.CreateHash()
	return newBlock
}

// function used to add new blocks to the blockchain
func (bc *Blockchain) AddBlock(transaction []Transaction) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(prevBlock.ID+1, transaction, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

// function to create the genesis block of the blockchain
func CreateGenesis() Block {
	genesis := NewBlock(0, transaction, "")
	genesis.Hash = genesis.CreateHash()
	return genesis
}

// Function to create the blockchain with the genesis block as the first block
func NewBlockchain() Blockchain {
	blockchain := CreateGenesis()
	return Blockchain{[]Block{blockchain}}
}

// function to create a new Wallet
func NewWallet() *Wallet {
	return &Wallet{Addresses: make(map[string]*Address)}
}

// function to create the address entities
func (W *Wallet) CreateAddress() string {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalln(err)
	}

	publicKey := privateKey.PublicKey
	publicKeyBytes := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
	add := sha256.Sum256(publicKeyBytes)
	address := hex.EncodeToString(add[:])

	W.Addresses[address] = &Address{
		PrivateKey: privateKey,
		PublicKey:  &publicKey,
		Balance:    0,
	}

	return address
}

// function to get the balance of the Addresses
func (W *Wallet) GetBalance(address string) (float64, error) {
	addrStr, exist := W.Addresses[address]
	if !exist {
		return 0, fmt.Errorf(address)
	}
	return addrStr.Balance, nil
}

// function to sign Transaction using sender's privatekey
func (w *Wallet) SignTransaction(tx *Transaction) error {
	addr, exist := w.Addresses[tx.Sender]
	if !exist {
		log.Fatalln(tx.Sender)
	}

	hash := sha256.Sum256([]byte(tx.Receiver + fmt.Sprintf("%f", tx.Amount) + tx.Sender + tx.TimeStamp))
	r, s, err := ecdsa.Sign(rand.Reader, addr.PrivateKey, hash[:])
	if err != nil {
		log.Fatalln(err)
	}

	signature := append(r.Bytes(), s.Bytes()...)

	tx.Signature = hex.EncodeToString(signature)

	return nil
}

// VerifyTransaction verifies the signature of a transaction

func (w *Wallet) VerifyTransaction(tx *Transaction) bool {
	addr, exist := w.Addresses[tx.Sender]
	if !exist {
		log.Fatalln("INVALID ADDRESS")
		return false
	}

	hash := sha256.Sum256([]byte(tx.Receiver + tx.Sender + tx.TimeStamp + fmt.Sprintf("%f", tx.Amount)))
	signature, err := hex.DecodeString(tx.Signature)
	if err != nil {
		log.Fatalln(err)
	}

	r := big.Int{}
	s := big.Int{}

	signLen := len(signature)

	r.SetBytes(signature[:(signLen / 2)])
	s.SetBytes(signature[(signLen / 2):])

	return ecdsa.Verify(&addr.PrivateKey.PublicKey, hash[:], &r, &s)
}

// function to use for transfaring the funds between the address
func (w *Wallet) Transfer(from, to string, amount float64) error {
	// var transaction []Transaction
	addrfrom, exist := w.Addresses[from]
	if !exist {
		log.Fatalln(from)
	}

	addrTo, exist := w.Addresses[to]
	if !exist {
		log.Fatalln(to)
	}

	if addrfrom.Balance < amount {
		log.Fatalln("INSAFICIANT FUNDS IN THE ACCOUNT")
	}

	transaction := Transaction{
		Sender:    from,
		Receiver:  to,
		Amount:    amount,
		TimeStamp: time.Now().String(),
	}

	if !w.VerifyTransaction(&transaction) {
		log.Fatalln("INVALID TRANSACTION")
	}

	err := w.SignTransaction(&transaction)
	if err != nil{
		log.Fatalln()
	}
	addrfrom.Balance -= amount
	addrTo.Balance += amount

	return nil
}
