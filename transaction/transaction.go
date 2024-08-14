package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

// struct Transaction to hold the transaction enetities
type Transaction struct {
	Sender    string
	Receiver  string
	Amount    float64
	TimeStamp string
}

// struct to hold the Blocks entities
type Block struct {
	ID           int
	Transactions []Transaction
	TimeStamp    string
	PrevHash     string
	Hash         string
}

// struct to hold the blocks in the blockchain
type Blockchain struct {
	Blocks []Block
}

var transaction []Transaction

// function to create the hash functionality for the blocks
func (b *Block) CreateHash() string {
	res := strconv.Itoa(b.ID) + b.PrevHash + b.TimeStamp + b.Hash
	for _, tx := range b.Transactions {
		res += tx.Sender + tx.Receiver + tx.TimeStamp + fmt.Sprintf("%f", tx.Amount)
	}
	h := sha256.New()
	h.Write([]byte(res))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// function to create the genesis block of the blockchain
func CreateGenesis() *Block {
	genesis := &Block{0, transaction, time.Now().String(), "", " "}

	genesis.Hash = genesis.CreateHash()
	return genesis
}

// function to create new block
func NewBlock(id int, transaction []Transaction, prevHash string) Block {
	newBlock := Block{
		ID:           id,
		Transactions: transaction,
		TimeStamp:    time.Now().String(),
		PrevHash:     prevHash,
		Hash:         "",
	}
	newBlock.Hash = newBlock.CreateHash()
	return newBlock
}

// Function to create the first block to the blockchain
func NewBlockchain() Blockchain {
	blockchain := CreateGenesis()
	return Blockchain{[]Block{*blockchain}}
}

// function to add a new bolck to the blockchain
func (bc *Blockchain) AddBlock(transaction []Transaction) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(prevBlock.ID+1, transaction, prevBlock.Hash)

	bc.Blocks = append(bc.Blocks, newBlock)
}

// Function main to handle the commputation of the transaction
func main() {
	blockchain := NewBlockchain()

	transaction := []Transaction{
		{
			Sender:    "paul",
			Amount:    30,
			Receiver:  "smally",
			TimeStamp: time.Now().String(),
		},
		{Sender: "Alice", Receiver: "Bob", Amount: 10, TimeStamp: time.Now().String()},
		{Sender: "Bob", Receiver: "Charlie", Amount: 5, TimeStamp: time.Now().String()},
	}

	blockchain.AddBlock(transaction)

	for _, block := range blockchain.Blocks {
		fmt.Printf("ID:  %d \n", block.ID)
		fmt.Printf("TimeStamp: %s \n", block.TimeStamp)
		fmt.Printf("Previous hash:  %s\n", block.PrevHash)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Printf("The Transactions: \n")

		for _, tx := range block.Transactions {
			fmt.Printf("%s --> %f --> %s \n", tx.Sender, tx.Amount, tx.Receiver)
		}
	}
}
