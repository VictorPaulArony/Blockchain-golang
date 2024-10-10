package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"
)

// Claim represents an insurance claim
type Claim struct {
	ID          string `json:"id"`
	PolicyID    string `json:"policy_id"`
	ClaimAmount int    `json:"claim_amount"`
	Status      string `json:"status"` 
	Timestamp   string `json:"timestamp"`
}

// Block represents a single block in the blockchain
type Block struct {
	ID        int
	Timestamp string
	Claims    []Claim
	PrevHash  string
	Hash      string
}

// Blockchain represents the entire blockchain
type Blockchain struct {
	Blocks []Block
	mu     sync.Mutex
}

// CalculateHash computes the hash for the block
func (b *Block) CalculateHash() string {
	data := strconv.Itoa(b.ID) + b.Timestamp + b.PrevHash
	for _, claim := range b.Claims {
		claimJSON, _ := json.Marshal(claim)
		data += string(claimJSON)
	}
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// AddBlock creates a new block and adds it to the blockchain
func (bc *Blockchain) AddBlock(claims []Claim) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	var prevHash string
	if len(bc.Blocks) > 0 {
		prevHash = bc.Blocks[len(bc.Blocks)-1].Hash
	}

	newBlock := Block{
		ID:        len(bc.Blocks) + 1,
		Timestamp: time.Now().String(),
		Claims:    claims,
		PrevHash:  prevHash,
	}

	newBlock.Hash = newBlock.CalculateHash()
	bc.Blocks = append(bc.Blocks, newBlock)
}

// RegisterClaim registers a new claim
func (bc *Blockchain) RegisterClaim(policyID string, amount int) Claim {
	claim := Claim{
		ID:          strconv.Itoa(len(bc.Blocks) + 1),
		PolicyID:    policyID,
		ClaimAmount: amount,
		Status:      "Pending",
		Timestamp:   time.Now().String(),
	}
	return claim
}

// ApproveClaim approves a claim if conditions are met
func (bc *Blockchain) ApproveClaim(claimID string) error {
	for i := range bc.Blocks {
		for j := range bc.Blocks[i].Claims {
			if bc.Blocks[i].Claims[j].ID == claimID {
				if bc.Blocks[i].Claims[j].Status != "Pending" {
					return fmt.Errorf("claim is already processed: %s", bc.Blocks[i].Claims[j].Status)
				}
				bc.Blocks[i].Claims[j].Status = "Approved"
				return nil
			}
		}
	}
	return fmt.Errorf("claim not found: %s", claimID)
}

// DenyClaim denies a claim
func (bc *Blockchain) DenyClaim(claimID string) error {
	for i := range bc.Blocks {
		for j := range bc.Blocks[i].Claims {
			if bc.Blocks[i].Claims[j].ID == claimID {
				if bc.Blocks[i].Claims[j].Status != "Pending" {
					return fmt.Errorf("claim is already processed: %s", bc.Blocks[i].Claims[j].Status)
				}
				bc.Blocks[i].Claims[j].Status = "Denied"
				return nil
			}
		}
	}
	return fmt.Errorf("claim not found: %s", claimID)
}

// PrintBlockchain prints the entire blockchain
func (bc *Blockchain) PrintBlockchain() {
	for _, block := range bc.Blocks {
		fmt.Printf("Block ID: %d\n", block.ID)
		fmt.Printf("Timestamp: %s\n", block.Timestamp)
		for _, claim := range block.Claims {
			fmt.Printf("  Claim ID: %s, Policy ID: %s, Amount: %d, Status: %s\n", claim.ID, claim.PolicyID, claim.ClaimAmount, claim.Status)
		}
		fmt.Printf("  Previous Hash: %s\n", block.PrevHash)
		fmt.Printf("  Hash: %s\n\n", block.Hash)
	}
}

func main() {
	blockchain := &Blockchain{}

	// Register claims
	claim1 := blockchain.RegisterClaim("policy_123", 1000)
	claim2 := blockchain.RegisterClaim("policy_456", 2000)

	// Add claims to the blockchain
	blockchain.AddBlock([]Claim{claim1})
	blockchain.AddBlock([]Claim{claim2})

	// Approve a claim
	if err := blockchain.ApproveClaim(claim1.ID); err != nil {
		fmt.Println(err)
	}

	// Deny a claim
	if err := blockchain.DenyClaim(claim2.ID); err != nil {
		fmt.Println(err)
	}

	// Print the blockchain
	blockchain.PrintBlockchain()
}