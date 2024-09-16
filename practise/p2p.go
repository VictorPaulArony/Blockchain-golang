package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

type Block struct {
	ID        int
	Data      string
	TimeStamp string
	PrevHash  string
	Hash      string
	Nonce     int
}

type Blockchain struct {
	blocks []Block
}

// function to create the hash function for the block
func (b *Block) CreatHash() string {
	res := strconv.Itoa(b.ID) + b.Data + b.TimeStamp + b.PrevHash + b.Hash + strconv.Itoa(b.Nonce)
	hash := sha256.New()
	hash.Write([]byte(res))
	hashed := hash.Sum(nil)

	return hex.EncodeToString(hashed)
}

// function to create a method to add a new block to the blockchain
func (bc *Blockchain) NewBlock(data string, diff int) {
	prevHash := bc.blocks[len(bc.blocks)-1]
	newBlock := Block{
		ID:        prevHash.ID + 1,
		Data:      data,
		TimeStamp: time.Now().String(),
		PrevHash:  prevHash.Hash,
		Nonce:     0,
	}

	for {
		newBlock.Hash = newBlock.CreatHash()
		if IsValidHash(newBlock.Hash, diff) {
			break
		}
		newBlock.Nonce++

	}

	bc.blocks = append(bc.blocks, newBlock)
}

// function to create the genesis block of the blockchain
func (bc *Blockchain) CreateGenesis() Block {
	genesis := Block{ID: 0, Data: "Genesis Block", TimeStamp: time.Now().String(), PrevHash: "", Nonce: 0}
	genesis.Hash = genesis.CreatHash()
	return genesis
}

// function to create blockchain and add the genesis block to the blockchain
func NewBlockchain() Blockchain {
	bc := Blockchain{}
	genesis := bc.CreateGenesis()
	bc.blocks = append(bc.blocks, genesis)
	return bc
}

// function to verify the hash if it meets the difficulity
func IsValidHash(hash string, diff int) bool {
	prefix := ""
	for i := 0; i < diff; i++ {
		prefix += "0"
	}
	return hash[:diff] == prefix
}

func main() {
	difficulty := 7 // Number of leading zeros required in the hash
	blockchain := NewBlockchain()

    prevBlockHash := blockchain.blocks[len(blockchain.blocks)-1].Hash
	

	fmt.Println("Mining...")
	blockchain.NewBlock("Transaction Data 1", difficulty)
	fmt.Printf("New Block Mined: %s\n",prevBlockHash )

	blockchain.NewBlock("Transaction Data 2", difficulty)
	fmt.Printf("New Block Mined: %s\n", prevBlockHash)

	for _, blk := range blockchain.blocks {
		fmt.Printf("Timestamp: %s\nData: %s\nPrevHash: %s\nHash: %s\nNonce: %d\n\n",
			blk.TimeStamp, blk.Data, blk.PrevHash, blk.Hash, blk.Nonce)
	}
}
