package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"time"
)

type Block struct {
	ID        int
	TimeStamp string
	Data      string
	PrevHash  string
	Hash      string
	Signature string
	PubKey    *ecdsa.PublicKey
}

type Blockchain struct {
	blocks []Block
}

func CreateHash(b *Block) string {
	res := strconv.Itoa(b.ID) + b.TimeStamp + b.Data + b.PrevHash
	h := sha256.New()
	h.Write([]byte(res))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func CreateGenesis(privKey *ecdsa.PrivateKey) Block {
	genesis := Block{ID: 0, TimeStamp: time.Now().String(), Data: "Genesis Block"}
	genesis.Hash = CreateHash(&genesis)

	r, s, err := ecdsa.Sign(rand.Reader, privKey, []byte(genesis.Hash))
	if err != nil {
		log.Fatalf("Failed to sign genesis block: %v", err)
	}

	genesis.Signature = fmt.Sprintf("%s%s", r.String(), s.String())
	genesis.PubKey = &privKey.PublicKey

	return genesis
}

func (bc *Blockchain) AddBlock(data string, privKey *ecdsa.PrivateKey) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := Block{
		ID:        prevBlock.ID + 1,
		TimeStamp: time.Now().String(),
		Data:      data,
		PrevHash:  prevBlock.Hash,
	}
	newBlock.Hash = CreateHash(&newBlock)

	r, s, err := ecdsa.Sign(rand.Reader, privKey, []byte(newBlock.Hash))
	if err != nil {
		log.Fatalf("Failed to sign new block: %v", err)
	}

	newBlock.Signature = fmt.Sprintf("%s%s", r.String(), s.String())
	newBlock.PubKey = &privKey.PublicKey

	bc.blocks = append(bc.blocks, newBlock)
}

func main() {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	genesisBlock := CreateGenesis(privKey)
	blockchain := Blockchain{[]Block{genesisBlock}}

	blockchain.AddBlock("second block", privKey)
	blockchain.AddBlock("third block", privKey)
	blockchain.AddBlock("fourth block", privKey)

	for _, block := range blockchain.blocks {
		fmt.Printf("ID: %d\n", block.ID)
		fmt.Printf("TimeStamp: %s\n", block.TimeStamp)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Previous Hash: %s\n", block.PrevHash)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Printf("Signature: %s\n", block.Signature)
		fmt.Printf("PublicKey: %s\n\n", elliptic.Marshal(elliptic.P256(), block.PubKey.X, block.PubKey.Y))
	}
}
