package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

type Transaction struct {
	Sender    string
	Receiver  string
	Amount    float64
	TimeStamp string
}

type Block struct {
	ID           int
	Transactions []Transaction
	PrevHash     string
	TimeStamp    string
	Hash         string
}
type Blockchain struct {
	Blocks []Block
}

// function to creat the hashing function
func (bc *Block) CreateHash() string {
	res := strconv.Itoa(bc.ID) + bc.TimeStamp + bc.PrevHash
	for _, tr := range bc.Transactions {
		res += tr.Sender + tr.Receiver + fmt.Sprintf("%f", tr.Amount) + tr.TimeStamp
	}
	h := sha256.New()
	h.Write([]byte(res))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// function to create a new block for the blockchain
func NewBlock(id int, transaction []Transaction, prevHash string) Block {
	newBlock := Block{
		ID:           id,
		Transactions: transaction,
		PrevHash:     prevHash,
		TimeStamp:    time.Now().String(),
		Hash:         "",
	}

	newBlock.Hash = newBlock.CreateHash()
	return newBlock
}

// function to create a new block the genesis blockchain
func CreateGenesis() *Block {
	var transaction []Transaction
	genesis := &Block{0, transaction, "", time.Now().String(), " "}
	genesis.Hash = genesis.CreateHash()
	return genesis
}

// method function to initialixe a new blockchain
func NewBlockchain() *Blockchain {
	genesisBlock := CreateGenesis()
	return &Blockchain{[]Block{*genesisBlock}}
}

// function and a method to add a new block to blockchain
func (bc *Blockchain) AddBlock(transaction []Transaction) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(prevBlock.ID+1, transaction, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func main() {
	blockchain := NewBlockchain()

	// Add a new block with transactions
	transactions := []Transaction{
		{Sender: "Alice", Receiver: "Bob", Amount: 10, TimeStamp: time.Now().String()},
		{Sender: "Bob", Receiver: "Charlie", Amount: 5, TimeStamp: time.Now().String()},
	}

	blockchain.AddBlock(transactions)

	// Print the blockchain
	for _, block := range blockchain.Blocks {
		fmt.Printf("Index: %d\n", block.ID)
		fmt.Printf("TimeStamp: %s\n", block.TimeStamp)
		fmt.Printf("Previous Hash: %s\n", block.PrevHash)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Printf("Transactions:\n")
		for _, tx := range block.Transactions {
			fmt.Printf("  %s -> %s: %f\n", tx.Sender, tx.Receiver, tx.Amount)
		}
		fmt.Println()
	}
}
