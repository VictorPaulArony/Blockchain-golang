package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	// "io/os"
	"log"
	"os"
	"time"
)

const (
	UTXOFile       = "simple_utxos.json"
	BlockchainFile = "simple_blockchain.json"
	CoinbaseReward = 10.0
)

// Transaction represents a single transaction
type Transaction struct {
	ID      string
	Inputs  []TXInput
	Outputs []TXOutput
}

// TXInput represents a transaction input
type TXInput struct {
	TxID     string
	OutIndex int
}

// TXOutput represents a transaction output
type TXOutput struct {
	Value   float64
	Address string
}

// Block represents a single block in the blockchain
type Block struct {
	Index        int
	Timestamp    string
	Transactions []Transaction
	PrevHash     string
	Hash         string
}

// Blockchain represents the entire blockchain
type Blockchain struct {
	Blocks  []Block
	UTXOSet UTXOSet
}

// UTXOSet represents all unspent transaction outputs
type UTXOSet struct {
	UTXOs map[string][]TXOutput
}

// NewBlockchain creates a new blockchain with a genesis block
func NewBlockchain(minerAddress string) *Blockchain {
	coinbaseTx := createCoinbaseTx(minerAddress)
	genesisBlock := createBlock([]Transaction{coinbaseTx}, "")

	utxoSet := UTXOSet{UTXOs: make(map[string][]TXOutput)}
	utxoSet.UTXOs[coinbaseTx.ID] = coinbaseTx.Outputs

	bc := Blockchain{
		Blocks:  []Block{genesisBlock},
		UTXOSet: utxoSet,
	}

	bc.saveBlockchain()
	bc.saveUTXOSet()

	return &bc
}

// createCoinbaseTx creates a new coinbase transaction
func createCoinbaseTx(address string) Transaction {
	output := TXOutput{Value: CoinbaseReward, Address: address}
	tx := Transaction{
		ID:      generateTxID(),
		Inputs:  nil,
		Outputs: []TXOutput{output},
	}
	return tx
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(transactions []Transaction) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := createBlock(transactions, prevBlock.Hash)

	bc.Blocks = append(bc.Blocks, newBlock)
	bc.updateUTXOSet(transactions)
	bc.saveBlockchain()
	bc.saveUTXOSet()
}

// createBlock creates a new block
func createBlock(transactions []Transaction, prevHash string) Block {
	block := Block{
		Index:        len(transactions) + 1,
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		PrevHash:     prevHash,
	}
	block.Hash = calculateBlockHash(block)
	return block
}

// calculateBlockHash calculates the hash of a block
func calculateBlockHash(block Block) string {
	data := fmt.Sprintf("%d%s%x%s", block.Index, block.Timestamp, block.Transactions, block.PrevHash)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// updateUTXOSet updates the UTXO set with new transactions
func (bc *Blockchain) updateUTXOSet(transactions []Transaction) {
	for _, tx := range transactions {
		// Add new outputs to UTXO set
		bc.UTXOSet.UTXOs[tx.ID] = tx.Outputs
	}
}

// generateTxID generates a new transaction ID
func generateTxID() string {
	id := make([]byte, 32)
	rand.Read(id)
	return hex.EncodeToString(id)
}

// saveUTXOSet saves the UTXO set to a JSON file
func (bc *Blockchain) saveUTXOSet() {
	data, err := json.Marshal(bc.UTXOSet)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(UTXOFile, data, 0o644)
	if err != nil {
		log.Fatal(err)
	}
}

// saveBlockchain saves the blockchain to a JSON file
func (bc *Blockchain) saveBlockchain() {
	data, err := json.Marshal(bc.Blocks)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(BlockchainFile, data, 0o644)
	if err != nil {
		log.Fatal(err)
	}
}

// loadBlockchain loads the blockchain from a JSON file
func loadBlockchain() ([]Block, error) {
	data, err := os.ReadFile(BlockchainFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var blocks []Block
	err = json.Unmarshal(data, &blocks)
	return blocks, err
}

// loadUTXOSet loads the UTXO set from a JSON file
func loadUTXOSet() (UTXOSet, error) {
	data, err := os.ReadFile(UTXOFile)
	if err != nil {
		if os.IsNotExist(err) {
			return UTXOSet{UTXOs: make(map[string][]TXOutput)}, nil
		}
		return UTXOSet{}, err
	}
	var utxoSet UTXOSet
	err = json.Unmarshal(data, &utxoSet)
	return utxoSet, err
}

func main() {
	// Initialize a new blockchain and add a coinbase transaction
	address := "your-miner-address"
	bc := NewBlockchain(address)

	// Add a block to the blockchain
	coinbaseTx := createCoinbaseTx(address)
	bc.AddBlock([]Transaction{coinbaseTx})

	fmt.Println("Blockchain initialized and first block added.")
}
