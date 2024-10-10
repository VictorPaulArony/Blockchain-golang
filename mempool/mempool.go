package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

type Transaction struct {
	Sender    string
	Receiver  string
	Amount    int
	TimeStamp string
}

type Block struct {
	ID          int
	Transaction []Transaction
	PrevHash    string
	TimeStamp   string
	Hash        string
	Nonce       int
}

type Mempool struct {
	Transaction []Transaction
	mu          sync.Mutex
}

type Blockchain struct {
	Blocks []Block
	mu     sync.Mutex
}

var mempool = CreateMempool()

// Method to create a hashing function for Blocks
func (b *Block) CreateHash() string {
	res := strconv.Itoa(b.ID) + b.PrevHash + b.TimeStamp + strconv.Itoa(b.Nonce)
	for _, tx := range b.Transaction {
		res += tx.Receiver + tx.Sender + tx.TimeStamp + fmt.Sprint(tx.Amount)
	}
	hash := sha256.Sum256([]byte(res))
	return hex.EncodeToString(hash[:])
}

// function to create Genesis Block of the blockchain
func GenesisBlock() Block {
	var transaction []Transaction
	genesis := Block{0, transaction, "", time.Now().String(), "", 0}
	genesis.Hash = genesis.CreateHash()
	return genesis
}

// function to create a new mempool
func CreateMempool() Mempool {
	return Mempool{Transaction: []Transaction{}}
}

// function to create a new transaction to the mempool
func (mp *Mempool) AddTransaction(from, to string, amount int) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	newTransaction := Transaction{
		Sender:    from,
		Receiver:  to,
		Amount:    amount,
		TimeStamp: time.Now().String(),
	}
	mp.Transaction = append(mp.Transaction, newTransaction)
}

// function to add New block to blockchain
func (bc *Blockchain) AddBlock() {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// lock mempool to retrive transaction
	mempool.mu.Lock()
	defer mempool.mu.Unlock()

	// retrieve trasactions from the mempool
	transaction := mempool.Transaction
	if len(transaction) == 0 {
		log.Println("No Transaction to Mine")
		return
	}

	prevBlock := bc.Blocks[len(bc.Blocks)-1]

	newBlock := Block{
		ID:          prevBlock.ID + 1,
		Transaction: transaction,
		PrevHash:    prevBlock.Hash,
		TimeStamp:   time.Now().String(),
	}

	// simple PoW
	for {
		newBlock.Hash = newBlock.CreateHash()
		if IsValidHash(newBlock.Hash) {
			break
		}
		newBlock.Nonce++

	}
	if !bc.IsValidBlock(newBlock) {
		log.Fatalln("Invalid Block: Previous Hash Does not Match")
	}

	bc.Blocks = append(bc.Blocks, newBlock)

	// clear the mempool
	mempool.Transaction = []Transaction{}
}

// function to verify the hash for nonce(PoW)
func IsValidHash(hash string) bool {
	return hash[:5] == "00000"
}

// function to create new blockchain
func NewBlockchain() *Blockchain {
	genesis := GenesisBlock()
	return &Blockchain{Blocks: []Block{genesis}}
}

// function to validate blocks to be added in the blockchain
func (bc *Blockchain) IsValidBlock(newBlock Block) bool {
	if len(bc.Blocks) == 0 {
		return false
	}
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	return newBlock.PrevHash == prevBlock.Hash
}

// Function to display the entire blockchain
func displayBlockchain(bc *Blockchain) {
	for _, block := range bc.Blocks {
		fmt.Printf("Block ID: %d\n", block.ID)
		fmt.Printf("Timestamp: %s\n", block.TimeStamp)
		fmt.Printf("Previous Hash: %s\n", block.PrevHash)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Println("Transactions:")
		for _, tx := range block.Transaction {
			fmt.Printf("\t%s -> %s: %d\n", tx.Sender, tx.Receiver, tx.Amount)
		}
		fmt.Println()
	}
}

// the main function using the CLI
func main() {
	blockchain := NewBlockchain()

	// Adding transaction to the mempool
	mempool.AddTransaction("smally", "pauls", 100)

	blockchain.AddBlock()
	displayBlockchain(blockchain)
}
