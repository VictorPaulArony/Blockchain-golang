# Mempool 
This Golang code demonstrates a basic blockchain system with a mempool for pending transactions, proof of work (PoW) for mining blocks, and functions to add new blocks to the chain. Let's break down the code and explain how it works in detail:

## Key Components:

1. Transaction Structure:
    - A transaction in this system includes the `Sender`, `Receiver`, `Amount`, and `TimeStamp`.
    - Example: A transaction where "Smally" sends 100 tokens to "Paul".
    
 ```go
    type Transaction struct {
    Sender    string
    Receiver  string
    Amount    int
    TimeStamp string
}
```
2. Block Structure:

   - Each block contains an `ID`, a list of `Transaction` entries, the `PrevHash` (link to the previous block), a `TimeStamp` when the block was created, a `Hash` (unique identifier of the block), and a `Nonce` (used for mining and proof of work).
   - The block's hash is calculated by hashing its content, including its ID, previous hash, transactions, and nonce.

```go
type Block struct {
    ID          int
    Transaction []Transaction
    PrevHash    string
    TimeStamp   string
    Hash        string
    Nonce       int
}
```
3. Mempool:

    - The mempool stores transactions that haven't yet been included in a block. This is essentially a buffer of transactions waiting to be mined.
    - Transactions can be added to the mempool using `AddTransaction`.
    - The `sync.Mutex` ensures that only one thread can access the mempool at a time, preventing race conditions in concurrent environments.

```go
type Mempool struct {
    Transaction []Transaction
    mu          sync.Mutex
}
```
4. Blockchain Structure:
    - The `Blockchain` is an array of blocks. It contains functions to add new blocks and verify the integrity of the blockchain.
    - `mu` is a mutex for locking the blockchain during the mining process to ensure data consistency.

```go
type Blockchain struct {
    Blocks []Block
    mu     sync.Mutex
}
```
## Key Functions:

1. **CreateHash():**
    - The `CreateHash()` method for a `Block` generates a unique SHA-256 hash for the block based on its `ID`, previous hash, transactions, timestamp, and nonce.
    - The function is crucial for validating the integrity of the block.

```go
func (b *Block) CreateHash() string {
    res := strconv.Itoa(b.ID) + b.PrevHash + b.TimeStamp + strconv.Itoa(b.Nonce)
    for _, tx := range b.Transaction {
        res += tx.Receiver + tx.Sender + tx.TimeStamp + fmt.Sprint(tx.Amount)
    }
    hash := sha256.Sum256([]byte(res))
    return hex.EncodeToString(hash[:])
}
```
2. GenesisBlock():

    - The Genesis Block is the first block in the blockchain and has no previous block to link to, so its `PrevHash` is an empty string.
    - The genesis block is created with an empty list of transactions.

```go
func GenesisBlock() Block {
    var transaction []Transaction
    genesis := Block{0, transaction, "", time.Now().String(), "", 0}
    genesis.Hash = genesis.CreateHash()
    return genesis
}
```

3. **AddTransaction():**

    - This function adds a transaction to the mempool by locking it with a mutex to ensure safe concurrent access. It creates a new Transaction and appends it to the mempool.

 ```go
    func (mp *Mempool) AddTransaction(from, to string, amount int) {
    mp.mu.Lock()
    defer mp.mu.Unlock()

    newTransaction := Transaction{
        Sender:    from,
        Receiver:  to,
        Amount:    amount,
        TimeStamp: time.Now().String(),
    }
    mp.Transaction = append(mp.Transaction, newTransaction)
}
```

4. **AddBlock():**

    -This function retrieves transactions from the mempool, constructs a new block, and performs proof of work (PoW) by adjusting the block’s Nonce until the hash meets the difficulty criteria (starts with 00000).
    - If the block is valid, it is added to the blockchain, and the mempool is cleared.

```go
func (bc *Blockchain) AddBlock() {
    bc.mu.Lock()
    defer bc.mu.Unlock()

    mempool.mu.Lock()
    defer mempool.mu.Unlock()

    transaction := mempool.Transaction
    if len(transaction) == 0 {
        log.Println("No Transaction to Mine")
        return
    }

    prevBlock := bc.Blocks[len(bc.Blocks)-1]

    newBlock := Block{
        ID:          prevBlock.ID + 1,
        Transaction: transaction,
        PrevHash:    prevBlock.Hash,
        TimeStamp:   time.Now().String(),
    }

    for {
        newBlock.Hash = newBlock.CreateHash()
        if IsValidHash(newBlock.Hash) {
            break
        }
        newBlock.Nonce++
    }

    if !bc.IsValidBlock(newBlock) {
        log.Fatalln("Invalid Block: Previous Hash Does not Match")
    }

    bc.Blocks = append(bc.Blocks, newBlock)
    mempool.Transaction = []Transaction{}
}
```

5. **IsValidHash():**

    A simple function to verify if a block's hash is valid. In this case, a hash is valid if it starts with 00000, a simplified form of Proof of Work (PoW).
    
```go
func IsValidHash(hash string) bool {
    return hash[:5] == "00000"
}
```

6. **NewBlockchain():**

    - This function initializes the blockchain with the Genesis Block. The genesis block is the starting point, and future blocks link to it.

```go
func NewBlockchain() *Blockchain {
    genesis := GenesisBlock()
    return &Blockchain{Blocks: []Block{genesis}}
}
```

7. **IsValidBlock():**

    This function checks if the newly created block is valid by verifying that its PrevHash matches the hash of the previous block.

```go
func (bc *Blockchain) IsValidBlock(newBlock Block) bool {
    if len(bc.Blocks) == 0 {
        return false
    }
    prevBlock := bc.Blocks[len(bc.Blocks)-1]
    return newBlock.PrevHash == prevBlock.Hash
}
```

8. displayBlockchain():

    - This utility function prints out the contents of the blockchain, showing the ID, timestamp, hashes, and transactions of each block.

```go
func displayBlockchain(bc *Blockchain) {
    for _, block := range bc.Blocks {
        fmt.Printf("Block ID: %d\n", block.ID)
        fmt.Printf("Timestamp: %s\n", block.TimeStamp)
        fmt.Printf("Previous Hash: %s\n", block.PrevHash)
        fmt.Printf("Hash: %s\n", block.Hash)
        fmt.Println("Transactions:")
        for _, tx := range block.Transaction {
            fmt.Printf("\t%s -> %s: %d\n", tx.Sender, tx.Receiver, tx.Amount)
        }
        fmt.Println()
    }
}
```

## How it Works:

    1. Mempool:
        The mempool stores all pending transactions. Transactions are added to the mempool via AddTransaction.

   2. Mining Process:
        To mine a block, transactions are retrieved from the mempool, and a new block is created. The system performs Proof of Work (PoW) by adjusting the Nonce value until the block hash starts with 00000. Once a valid hash is found, the block is added to the blockchain, and the mempool is cleared.

    3. Proof of Work (PoW):
        The Nonce is incremented until the block's hash meets the difficulty target (00000). This is a basic PoW mechanism that ensures some computational effort is required to mine a block.

    4. Chain Validation:
        After a block is mined, it is checked against the previous block in the chain. If the PrevHash of the new block matches the hash of the last block, the block is considered valid and added to the chain.

## Example Workflow:

   - A transaction is added to the mempool: "Alice" sends 100 units to "Bob".
    - The blockchain mines a new block containing this transaction. During the mining process, the system performs proof of work to generate a valid block hash.
    - The new block is added to the blockchain, and the mempool is cleared.

This basic blockchain system simulates the core concepts of mempool, block mining, PoW, and transaction management in a blockchain.

## full Code
```go
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

type Transaction struct {
	Sender    string
	Receiver  string
	Amount    int
	TimeStamp string
}

type Block struct {
	ID          int
	Transaction []Transaction
	PrevHash    string
	TimeStamp   string
	Hash        string
	Nonce       int
}

type Mempool struct {
	Transaction []Transaction
	mu          sync.Mutex
}

type Blockchain struct {
	Blocks []Block
	mu     sync.Mutex
}

var mempool = CreateMempool()

// Method to create a hashing function for Blocks
func (b *Block) CreateHash() string {
	res := strconv.Itoa(b.ID) + b.PrevHash + b.TimeStamp + strconv.Itoa(b.Nonce)
	for _, tx := range b.Transaction {
		res += tx.Receiver + tx.Sender + tx.TimeStamp + fmt.Sprint(tx.Amount)
	}
	hash := sha256.Sum256([]byte(res))
	return hex.EncodeToString(hash[:])
}

// function to create Genesis Block of the blockchain
func GenesisBlock() Block {
	var transaction []Transaction
	genesis := Block{0, transaction, "", time.Now().String(), "", 0}
	genesis.Hash = genesis.CreateHash()
	return genesis
}

// function to create a new mempool
func CreateMempool() Mempool {
	return Mempool{Transaction: []Transaction{}}
}

// function to create a new transaction to the mempool
func (mp *Mempool) AddTransaction(from, to string, amount int) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	newTransaction := Transaction{
		Sender:    from,
		Receiver:  to,
		Amount:    amount,
		TimeStamp: time.Now().String(),
	}
	mp.Transaction = append(mp.Transaction, newTransaction)
}

// function to add New block to blockchain
func (bc *Blockchain) AddBlock() {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// lock mempool to retrive transaction
	mempool.mu.Lock()
	defer mempool.mu.Unlock()

	// retrieve trasactions from the mempool
	transaction := mempool.Transaction
	if len(transaction) == 0 {
		log.Println("No Transaction to Mine")
		return
	}

	prevBlock := bc.Blocks[len(bc.Blocks)-1]

	newBlock := Block{
		ID:          prevBlock.ID + 1,
		Transaction: transaction,
		PrevHash:    prevBlock.Hash,
		TimeStamp:   time.Now().String(),
	}

	// simple PoW
	for {
		newBlock.Hash = newBlock.CreateHash()
		if IsValidHash(newBlock.Hash) {
			break
		}
		newBlock.Nonce++

	}
	if !bc.IsValidBlock(newBlock) {
		log.Fatalln("Invalid Block: Previous Hash Does not Match")
	}

	bc.Blocks = append(bc.Blocks, newBlock)

	// clear the mempool
	mempool.Transaction = []Transaction{}
}

// function to verify the hash for nonce(PoW)
func IsValidHash(hash string) bool {
	return hash[:5] == "00000"
}

// function to create new blockchain
func NewBlockchain() *Blockchain {
	genesis := GenesisBlock()
	return &Blockchain{Blocks: []Block{genesis}}
}

// function to validate blocks to be added in the blockchain
func (bc *Blockchain) IsValidBlock(newBlock Block) bool {
	if len(bc.Blocks) == 0 {
		return false
	}
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	return newBlock.PrevHash == prevBlock.Hash
}

// Function to display the entire blockchain
func displayBlockchain(bc *Blockchain) {
	for _, block := range bc.Blocks {
		fmt.Printf("Block ID: %d\n", block.ID)
		fmt.Printf("Timestamp: %s\n", block.TimeStamp)
		fmt.Printf("Previous Hash: %s\n", block.PrevHash)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Println("Transactions:")
		for _, tx := range block.Transaction {
			fmt.Printf("\t%s -> %s: %d\n", tx.Sender, tx.Receiver, tx.Amount)
		}
		fmt.Println()
	}
}

// the main function using the CLI
func main() {
	blockchain := NewBlockchain()

	// Adding transaction to the mempool
	mempool.AddTransaction("smally", "pauls", 100)

	blockchain.AddBlock()
	displayBlockchain(blockchain)
}

```