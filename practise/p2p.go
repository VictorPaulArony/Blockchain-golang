package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// Block structure
type Block struct {
	Index     int
	Timestamp string
	Data      string
	PrevHash  string
	Hash      string
	Nonce     int
}

// Blockchain structure
type Blockchain struct {
	blocks []*Block
	mu     sync.Mutex
}

// Create a new block
func NewBlock(index int, data string, prevHash string, nonce int) *Block {
	block := &Block{
		Index:     index,
		Timestamp: time.Now().String(),
		Data:      data,
		PrevHash:  prevHash,
		Nonce:     nonce,
	}
	block.Hash = block.calculateHash()
	return block
}

// Calculate the block's hash
func (b *Block) calculateHash() string {
	record := fmt.Sprintf("%d%s%s%s%d", b.Index, b.Timestamp, b.Data, b.PrevHash, b.Nonce)
	hash := sha256.New()
	hash.Write([]byte(record))
	return hex.EncodeToString(hash.Sum(nil))
}

// Create a new blockchain
func NewBlockchain() *Blockchain {
	return &Blockchain{blocks: []*Block{NewBlock(0, "Genesis Block", "", 0)}}
}

// Add a block to the blockchain
func (bc *Blockchain) AddBlock(data string) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	lastBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(lastBlock.Index+1, data, lastBlock.Hash, 0)

	// Simple proof of work
	for !isValidHash(newBlock.Hash) {
		newBlock.Nonce++
		newBlock.Hash = newBlock.calculateHash()
	}

	bc.blocks = append(bc.blocks, newBlock)
}

// Check if the hash is valid (e.g., starts with 0000)
func isValidHash(hash string) bool {
	return hash[:5] == "00000"
}

func main() {
	blockchain := NewBlockchain()

	// Add blocks
	blockchain.AddBlock("First block after Genesis")
	blockchain.AddBlock("Second block after Genesis")

	// Print the blockchain
	for _, block := range blockchain.blocks {
		fmt.Printf("Index: %d\n", block.Index)
		fmt.Printf("Timestamp: %s\n", block.Timestamp)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Previous Hash: %s\n", block.PrevHash)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Printf("Nonce: %d\n\n", block.Nonce)
	}
}
