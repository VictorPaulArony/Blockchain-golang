package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Certificate struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	TimeStamp string `json:"timestamp"`
	PrevHash  string `json:"prevhash"`
	Hash      string `json:"hash"`
}

type Blockchain struct {
	Certificates []Certificate `json:"certificates"`
}

var fileName = "blocks.json"

// GenerateHash creates a SHA-256 hash of the certificate data
func GenerateHash(b *Certificate) string {
	res := strconv.Itoa(b.ID) + b.PrevHash + b.TimeStamp + b.Name + b.Hash
	hash := sha256.New()
	hash.Write([]byte(res))
	return hex.EncodeToString(hash.Sum(nil))
}

// CreateGenesis creates the first block
func CreateGenesis() Certificate {
	genesis := Certificate{0, "Genesis Certificate", time.Now().String(), "", ""}
	genesis.Hash = GenerateHash(&genesis)
	return genesis
}

// AddBlock adds the block to the blockchain
func (bc *Blockchain) AddBlock(cert string) {
	prevCert := bc.Certificates[len(bc.Certificates)-1]
	newCert := Certificate{
		ID:        prevCert.ID + 1,
		Name:      cert,
		PrevHash:  prevCert.Hash,
		TimeStamp: time.Now().String(),
	}
	newCert.Hash = GenerateHash(&newCert)
	bc.Certificates = append(bc.Certificates, newCert)
}

// SaveBlocks saves blockchain data to the JSON file
func (bc *Blockchain) SaveBlocks() error {
	data, err := json.MarshalIndent(bc, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling blockchain data: %w", err)
	}
	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing blockchain data to file: %w", err)
	}
	return nil
}

// LoadBlockchain loads the blockchain from the JSON file
func (bc *Blockchain) LoadBlockchain() error {
	file, err := os.ReadFile(fileName)
	if os.IsNotExist(err) {
		fmt.Println("Blockchain file not found, creating a new one with genesis block.")
		bc.Certificates = []Certificate{CreateGenesis()}
		return bc.SaveBlocks()
	}
	if err != nil {
		return fmt.Errorf("error reading blockchain file: %w", err)
	}
	err = json.Unmarshal(file, bc)
	if err != nil {
		return fmt.Errorf("error unmarshalling blockchain data: %w", err)
	}
	return nil
}
