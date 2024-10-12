# Simple Blockchain Implementation in Go

## Overview

This project demonstrates a simplified version of a blockchain implemented in Go. It showcases the fundamental concepts of blockchain technology, including blocks, hashing, and proof of work (PoW). The implementation allows users to create a chain of blocks, each containing data, a timestamp, a reference to the previous block, and a cryptographic hash.

## Table of Contents

- [What is a Blockchain?](#what-is-a-blockchain)
- [Key Components](#key-components)
- [How the Code Works](#how-the-code-works)
  - [Block Structure](#block-structure)
  - [Blockchain Structure](#blockchain-structure)
  - [Creating the Genesis Block](#creating-the-genesis-block)
  - [Adding Blocks](#adding-blocks)
  - [Hashing and Validating Blocks](#hashing-and-validating-blocks)
- [Running the Code](#running-the-code)
- [Conclusion](#conclusion)

## What is a Blockchain?

A blockchain is a decentralized digital ledger that records transactions across many computers in such a way that the registered transactions cannot be altered retroactively. Each block contains a list of transactions and is linked to the previous block, forming a chain. This structure makes blockchains resistant to modification and fraud.

## Key Components

1. **Block**: The fundamental unit of a blockchain that contains data, a timestamp, a hash of the previous block, and its own hash.
2. **Blockchain**: A collection of blocks that are linked together in a specific order.
3. **Hashing**: The process of generating a fixed-size string of characters (a hash) from data, which serves as a unique identifier for that data.
4. **Proof of Work (PoW)**: A consensus mechanism used to validate new blocks by requiring computational effort to generate a valid hash.

## How the Code Works

### Block Structure

The `Block` struct represents an individual block in the blockchain with the following fields:

- `ID`: The unique identifier for the block.
- `Data`: The information stored in the block.
- `TimeStamp`: The time at which the block was created.
- `PreHash`: The hash of the previous block in the chain.
- `Hash`: The hash of the current block.
- `Nonce`: A variable used in the PoW algorithm to find a valid hash.

### Blockchain Structure

The `Blockchain` struct holds a slice of `Block` objects and a mutex for safe concurrent access:

```go
type Blockchain struct {
    Blocks []Block
    mu     sync.Mutex
}
```

### Creating the Genesis Block

The first block in the blockchain is called the "genesis block." It is created with a specific structure and serves as the foundation for the subsequent blocks:

```go 
func GenesisBlock() Block {
    genesis := Block{0, "Genesis Block", time.Now().String(), "", "", 0}
    genesis.Hash = genesis.CreateHash()
    return genesis
}
```
### Adding Blocks
New blocks can be added to the blockchain using the AddBlock method. This method performs the following steps:

1. Locks the blockchain for safe concurrent access.
2. Creates a new block based on the previous block's hash.
3. Uses a nonce to find a valid hash that meets the PoW requirement (in this case, the hash must start with five zeros).
4. Validates the new block and appends it to the blockchain.

### Hashing and Validating Blocks

The CreateHash method generates a hash for the block using SHA-256, and the IsValidHash function checks if the hash meets the criteria for PoW.

```go 
func (b *Block) CreateHash() string {
    res := strconv.Itoa(b.ID) + b.Data + b.TimeStamp + b.PreHash + b.Hash
    hash := sha256.Sum256([]byte(res))
    return hex.EncodeToString(hash[:])
}
```
### Running the Code
To run the blockchain implementation, follow these steps:

1. Install Go: Ensure you have Go installed on your machine. If not, you can download it from golang.org.
2. Create a new Go file: Copy the provided code into a new file named main.go.
Run the application:
2. Open a terminal and navigate to the directory containing your main.go file. Run the following command:

```bash
go run main.go
```
### Conclusion
This simple blockchain implementation in Go demonstrates the core principles of blockchain technology, including block creation, hashing, and proof of work. While this example is simplified, it provides a foundational understanding of how blockchains work and can be expanded with additional features such as transaction handling, network protocols, and consensus mechanisms.