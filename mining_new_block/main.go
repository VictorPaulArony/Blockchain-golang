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

type Block struct {
	ID        int
	Data      string
	TimeStamp string
	PreHash   string
	Hash      string
	Nonce     int
}

type Blockchain struct {
	Blocks []Block
	mu     sync.Mutex
}

// function to create the hashing function
func (b *Block) CreateHash() string {
	res := strconv.Itoa(b.ID) + b.Data + b.TimeStamp + b.PreHash + b.Hash
	hash := sha256.Sum256([]byte(res))
	return hex.EncodeToString(hash[:])
}

// function to add a block for the blockchain
func (bc *Blockchain) AddBlock(data string) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := Block{
		ID:        prevBlock.ID + 1,
		Data:      data,
		TimeStamp: time.Now().String(),
		PreHash:   prevBlock.Hash,
	}

	for {
		newBlock.Hash = newBlock.CreateHash()
		if IsValidHash(newBlock.Hash) {
			break
		}
		newBlock.Nonce++
	}

	if bc.IsValidBlock(newBlock) {
		bc.Blocks = append(bc.Blocks, newBlock)
	} else {
		log.Fatalln("Invalid Block")
	}
}

// function to verify the PoW
func IsValidHash(hash string) bool {
	return hash[:5] == "00000"
}

// function to create genesis block
func GenesisBlock() Block {
	genesis := Block{0, "Genesis Block", time.Now().String(), "", "", 0}
	genesis.Hash = genesis.CreateHash()
	return genesis
}

// function to create a new blockchain
func CreateBlockchain() Blockchain {
	genesis := GenesisBlock()
	return Blockchain{Blocks: []Block{genesis}}
}

// function to verify the blocks
func (bc *Blockchain) IsValidBlock(block Block) bool {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	if len(bc.Blocks) == 0 {
		return false
	}
	return prevBlock.Hash == block.PreHash
}

func main() {
	blockchain := CreateBlockchain()

	blockchain.AddBlock("second block")
	blockchain.AddBlock("third block")
	blockchain.AddBlock("forth block")
	blockchain.AddBlock("fifth block")
	blockchain.AddBlock("sixth block")
	blockchain.AddBlock("seventh block")

	for _, block := range blockchain.Blocks {
		fmt.Printf("Index: %d\n", block.ID)
		fmt.Printf("Timestamp: %s\n", block.TimeStamp)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("PrevHash: %s\n", block.PreHash)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Println()
	}
}
