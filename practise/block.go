package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

type Blockchain struct {
	blocks []Block 
}
type Block struct {
	ID        int    `json:"id"`
	Data      string `json:"data"`
	PrevHash  string `json:"prevHash"`
	Timestamp string `json:"timeStamp"`
	Hash      string `json:"hash"`
}

// function to create a new hash for the blocks
func (b *Block) CreateHash() string {
	res := strconv.Itoa(b.ID) + b.Data + b.PrevHash + b.Timestamp + b.Hash
	h := sha256.New()
	h.Write([]byte(res))
	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}

// function  to create a new block to add to the or for the blockchain
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := Block{
		ID:        prevBlock.ID + 1,
		Data:      data,
		PrevHash:  prevBlock.Hash,
		Timestamp: time.Now().String(),
	}
	newBlock.Hash = newBlock.CreateHash()
	bc.blocks = append(bc.blocks, newBlock)
}

// function to create the genesis block of the blockchain
func CreateGenesis() Block {
	genesis := Block{ID: 0, Data: "genesis Block", Timestamp: time.Now().String(), PrevHash: ""}
	genesis.Hash = genesis.CreateHash()
	return genesis
}

// function main that displays the blocks
func main() {
	genesis := CreateGenesis()

	blockchain := Blockchain{[]Block{genesis}}

	blockchain.AddBlock("First Block")
	blockchain.AddBlock("second Block")
	blockchain.AddBlock("Third Block")

	for _, blocks := range blockchain.blocks {
		fmt.Printf("%d\n", blocks.ID)
		fmt.Printf("%s\n", blocks.Data)
		fmt.Printf("%s\n", blocks.PrevHash)
		fmt.Printf("%s\n", blocks.Timestamp)
		fmt.Printf("%s\n", blocks.Hash)
	}
}
