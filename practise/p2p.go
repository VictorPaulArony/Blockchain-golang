package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type Block struct {
	ID          int
	Transaction []Transaction
	TimeStamp   string
	MerkleRoot  string
	PrevHash    string
	Hash        string
	Nonce       int
}

type Transaction struct {
	Sender   string
	Receiver string
	Amount   float64
}

type Blockchain struct {
	Blocks []Block
	mu     sync.Mutex
}

// Function to generate the a salt
//(to prevent collision attack to the blocks during hashing)
func Salt() (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

// fuction to create do the hashing
func Hash(data string, salt string) string {
	hash := sha256.Sum256([]byte(data + salt))
	return hex.EncodeToString(hash[:])
}

// function to hash the transaction block
func (b *Block) HashBlock() string {
	data := strconv.Itoa(b.ID) + b.TimeStamp + b.MerkleRoot + b.PrevHash + b.Hash + strconv.Itoa(b.Nonce)
	for _, tx := range b.Transaction {
		data += tx.Sender + tx.Receiver + strconv.Itoa(int(tx.Amount))
	}
	salt, _ := Salt()
	return Hash(data, salt)
}

// function to create the genesis block of thwe blockchain
func GenesisBlock() Block {
	var transaction []Transaction
	genesis := Block{
		ID:          0,
		Transaction: transaction,
		TimeStamp:   time.Now().String(),
		PrevHash:    "",
	}
	for {
		genesis.Hash = genesis.HashBlock()
		if IsValidHash(genesis.Hash) {
			break
		}
		genesis.Nonce++
	}
	return genesis
}

// function to verify the nonce of the block (PoW)
func IsValidHash(hash string) bool {
	return hash[:5] == "00000"
}

// function to create a new blockchain
func CreateBlockchain() Blockchain {
	genesis := GenesisBlock()
	return Blockchain{Blocks: []Block{genesis}}
}

// Function to verify the blocks before adding to the blockchain
func (bc *Blockchain) IsValidBlock(block Block) bool {
	if len(block.Hash) == 0 {
		return false
	}
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	if prevBlock.Hash == block.PrevHash {
		return true
	}
	return false
}

// function to Add the block to the blockchain
func (bc *Blockchain) AddBlock(transaction []Transaction) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := Block{
		ID:          prevBlock.ID + 1,
		Transaction: transaction,
		TimeStamp:   time.Now().String(),
		PrevHash:    prevBlock.Hash,
	}

	newBlock.MerkleRoot = CreateMerkelTree(transaction)
	for {
		newBlock.Hash = newBlock.HashBlock()
		if IsValidHash(newBlock.Hash) {
			break
		}
		newBlock.Nonce++
	}

	if !bc.IsValidBlock(newBlock) {
		log.Fatalln("The Block is an Invalid Block")
		return
	}
	bc.Blocks = append(bc.Blocks, newBlock)
}

// function to calculate the merkle root from a slice of transaction
func CreateMerkelTree(transaction []Transaction) string {
	var hashes []string

	if len(transaction) == 0 {
		return ""
	}
	salt, _ := Salt()
	for _, tx := range transaction {
		hashes = append(hashes, Hash(tx.Receiver+tx.Sender+strconv.Itoa(int(tx.Amount)), salt))
	}

	for len(hashes) > 1 {
		var newHash []string
		for i := 0; i < len(hashes); i += 2 {
			if i+1 < len(hashes) {
				newHash = append(newHash, Hash(hashes[i]+hashes[i+1], salt))
			} else {
				newHash = append(newHash, hashes[i])
			}
		}
		hashes = newHash
	}

	return hashes[0]
}

// function to print the output on the terminal of the CLI
func (bc *Blockchain) Display() {
	for _, block := range bc.Blocks {
		fmt.Printf("Block ID: %d\n", block.ID)
		fmt.Printf("  MerkleRoot: %s\n", block.MerkleRoot)
		fmt.Printf("  PrevHash: %s\n", block.PrevHash)
		fmt.Printf("  Hash: %s\n", block.Hash)
		fmt.Printf("  Nonce: %d\n", block.Nonce)
		for _, tx := range block.Transaction {
			fmt.Printf("Sender %s to Receiver %s amount %.f\n", tx.Receiver, tx.Sender, tx.Amount)
		}
	}
}

// function main to run the code on CLI mode
func main() {
	blockchain := CreateBlockchain()

	blockchain.AddBlock(
		[]Transaction{
			{
				Sender:   "paul",
				Receiver: "Smally",
				Amount:   150.0,
			},
		},
	)
	blockchain.Display()
}
