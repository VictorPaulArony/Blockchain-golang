package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Index     int
	TimeStamp string
	Data      string
	PreHash   string
	Hash      string
	Nonce     int
}

type Blockchain struct {
	block []Block
	diff  int
}

//var block = Block{}

func Calc(block Block) string {
	res := strconv.Itoa(block.Index) + block.TimeStamp + block.Data + block.PreHash + strconv.Itoa(block.Nonce)
	h := sha256.New()
	h.Write([]byte(res))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func (b *Block) MineBlock(diff int) {
	target := strings.Repeat("0", diff)
	for !strings.HasPrefix(b.Hash, target) {
		b.Nonce++
		b.Hash = Calc(*b)
	}
}

func CreateGenesis() Block {
	genesisBlock := Block{0, time.Now().String(), "First block", "", "", 0}
	genesisBlock.Hash = Calc(genesisBlock)
	genesisBlock.MineBlock(3)
	return genesisBlock
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.block[len(bc.block)-1]
	block := Block{
	Index: prevBlock.Index + 1,
	TimeStamp: time.Now().String(),
	Data: data,
	PreHash: prevBlock.Hash,
	}

	block.MineBlock(bc.diff)
	bc.block = append(bc.block, block)
}

func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.block); i++ {
		prevBlock := bc.block[i-1]
		currentBlock := bc.block[i]
		if currentBlock.PreHash != prevBlock.Hash {
			return false
		}
		if currentBlock.Hash != Calc(currentBlock) {
			return false
		}

	}
	return true
}

func main() {
	currentGenesis := CreateGenesis()
	blockchain := Blockchain{[]Block{currentGenesis}, 3}

	blockchain.AddBlock("second block")
	blockchain.AddBlock("third block")
	blockchain.AddBlock("forth block")
	blockchain.AddBlock("fifth block")
	blockchain.AddBlock("sixth block")
	blockchain.AddBlock("seventh block")

	for _, block := range blockchain.block {
		fmt.Printf("Index: %d\n", block.Index)
		fmt.Printf("Timestamp: %s\n", block.TimeStamp)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("PrevHash: %s\n", block.PreHash)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Println()
	}

	isValid := blockchain.IsValid()
	fmt.Printf("Is blockchain valid? %v\n", isValid)

	

}
