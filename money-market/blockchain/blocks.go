package blockchains

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	helpers "money-market/utils"
)

type Block struct {
	Index        int                   `json:"index"`
	Timestamp    string                `json:"timestamp"`
	Transactions []helpers.Transaction `json:"transactions"`
	PrevHash     string                `json:"prev_hash"`
	Hash         string                `json:"hash"`
}

type Blockchain struct {
	Blocks  []Block               `json:"blocks"`
	Mempool []helpers.Transaction `json:"mempool"`
}

const BlockchainFile = "blockchain.json"

var blockchain Blockchain

// function to initialize the blockchain
func InitializeBlockchain() {
	blockchain = LoadBlockchain()

	// Create a genesis block if the blockchain is empty
	if len(blockchain.Blocks) == 0 {
		genesisBlock := Block{
			Index:        0,
			Timestamp:    time.Now().Format(time.RFC3339),
			Transactions: []helpers.Transaction{},
			PrevHash:     "0",
			Hash:         helpers.GenerateHash("Genesis Block"),
		}
		blockchain.Blocks = append(blockchain.Blocks, genesisBlock)
		SaveBlockchain()
	}
}

// function to load the blockchain fron the db
func LoadBlockchain() Blockchain {
	data, err := os.ReadFile(BlockchainFile)
	if err != nil {
		if os.IsNotExist(err) {
			return Blockchain{} // Return an empty blockchain
		}
		log.Fatalf("Failed to load blockchain: %v", err)
	}
	var bc Blockchain
	json.Unmarshal(data, &bc)
	return bc
}

func SaveBlockchain() {
	data, err := json.MarshalIndent(blockchain, "", "  ")
	if err != nil {
		log.Fatalf("Failed to save blockchain: %v", err)
	}
	os.WriteFile(BlockchainFile, data, 0o644)
}

func MineBlock() {
	if len(blockchain.Mempool) == 0 {
		log.Println("No transactions to mine.")
		return
	}

	prevBlock := blockchain.Blocks[len(blockchain.Blocks)-1]
	newBlock := Block{
		Index:        len(blockchain.Blocks),
		Timestamp:    time.Now().Format(time.RFC3339),
		Transactions: blockchain.Mempool,
		PrevHash:     prevBlock.Hash,
	}

	// Generate the hash for the new block
	blockData := fmt.Sprintf("%d%s%v%s", newBlock.Index, newBlock.Timestamp, newBlock.Transactions, newBlock.PrevHash)
	newBlock.Hash = helpers.GenerateHash(blockData)

	blockchain.Blocks = append(blockchain.Blocks, newBlock)
	blockchain.Mempool = []helpers.Transaction{} // Clear the mempool

	SaveBlockchain()
	log.Printf("Block %d mined successfully.", newBlock.Index)
}

func AddTransactionToMempool(transaction helpers.Transaction) {
	blockchain.Mempool = append(blockchain.Mempool, transaction)
	SaveBlockchain()
	log.Printf("Transaction %s added to mempool.", transaction.ID)
}
