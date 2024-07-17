# **Simple Blockchain Implementation in Go**
This documentation provides an overview and detailed explanation of a simple blockchain implementation in Go. The code is divided into sections for clarity and ease of understanding.

## Table of Contents
- Introduction
- Block Structure
- Blockchain Structure
- Hash Function
- Creating the Genesis Block
- Adding New Blocks
- Main Function
- Complete Code

## Introduction

A blockchain is a series of blocks, each of which carries data and is connected to the ones before it by a cryptographic hash. The basic ideas of blockchain, such as block creation, hashing, and block linking, are illustrated by this implementation.

## Block Structure
Each block in the blockchain has the following attributes:

- `ID:` A unique identifier for the block.
- `TimeStamp:` The time when the block was created.
- `Data:` The data contained in the block.
- `PrevHash:` The hash of the previous block.
- `Hash:` The hash of the current block.

```go
type Block struct {
    ID        int
    TimeStamp string
    Data      string
    PrevHash  string
    Hash      string
}
```
## Blockchain Structure
The blockchain itself is a slice of blocks. It starts with a genesis block, and new blocks are added to it.

```go
type Blockchain struct {
    blocks []Block
}
```

## Hash Function
The CreateHash function generates a cryptographic hash for a block using its ID, timestamp, data, and previous hash. The MD5 algorithm is used for hashing in this example.

```go
func CreateHash(b *Block) string {
    res := strconv.Itoa(b.ID) + b.TimeStamp + b.Data + b.PrevHash
    h := md5.New()
    h.Write([]byte(res))
    hashed := h.Sum(nil)
    return hex.EncodeToString(hashed)
}
```

## Creating the Genesis Block
The genesis block is the first block in the blockchain and is created with predefined data. It does not have a previous hash.

```go
func CreateGenesis() Block {
    genesis := Block{ID: 0, TimeStamp: time.Now().String(), Data: "Genesis Block"}
    genesis.Hash = CreateHash(&genesis)
    return genesis
}
```

## Adding New Blocks
New blocks are added to the blockchain with the AddBlock method. Each new block contains data, the current timestamp, and the hash of the previous block.

```go
func (bc *Blockchain) AddBlock(data string) {
    prevBlock := bc.blocks[len(bc.blocks)-1]
    newBlock := Block{
        ID:        prevBlock.ID + 1,
        TimeStamp: time.Now().String(),
        Data:      data,
        PrevHash:  prevBlock.Hash,
    }
    newBlock.Hash = CreateHash(&newBlock)
    bc.blocks = append(bc.blocks, newBlock)
}
```
## Main Function
The main function initializes the blockchain with a genesis block and adds a few more blocks. It then prints the details of each block in the blockchain.

```go
func main() {
    genesisBlock := CreateGenesis()
    blockchain := Blockchain{[]Block{genesisBlock}}

    blockchain.AddBlock("second block")
    blockchain.AddBlock("third block")
    blockchain.AddBlock("fourth block")

    for _, block := range blockchain.blocks {
        fmt.Printf("ID: %d\n", block.ID)
        fmt.Printf("TimeStamp: %s\n", block.TimeStamp)
        fmt.Printf("Data: %s\n", block.Data)
        fmt.Printf("Previous Hash: %s\n", block.PrevHash)
        fmt.Printf("Hash: %s\n", block.Hash)
    }
}
```
## Complete Code
Here is the complete code for the simple blockchain implementation:

```go 
package main

import (
    "crypto/md5"
    "encoding/hex"
    "fmt"
    "strconv"
    "time"
)

type Block struct {
    ID        int
    TimeStamp string
    Data      string
    PrevHash  string
    Hash      string
}

type Blockchain struct {
    blocks []Block
}

func CreateHash(b *Block) string {
    res := strconv.Itoa(b.ID) + b.TimeStamp + b.Data + b.PrevHash
    h := md5.New()
    h.Write([]byte(res))
    hashed := h.Sum(nil)
    return hex.EncodeToString(hashed)
}

func CreateGenesis() Block {
    genesis := Block{ID: 0, TimeStamp: time.Now().String(), Data: "Genesis Block"}
    genesis.Hash = CreateHash(&genesis)
    return genesis
}

func (bc *Blockchain) AddBlock(data string) {
    prevBlock := bc.blocks[len(bc.blocks)-1]
    newBlock := Block{
        ID:        prevBlock.ID + 1,
        TimeStamp: time.Now().String(),
        Data:      data,
        PrevHash:  prevBlock.Hash,
    }
    newBlock.Hash = CreateHash(&newBlock)
    bc.blocks = append(bc.blocks, newBlock)
}

func main() {
    genesisBlock := CreateGenesis()
    blockchain := Blockchain{[]Block{genesisBlock}}

    blockchain.AddBlock("second block")
    blockchain.AddBlock("third block")
    blockchain.AddBlock("fourth block")

    for _, block := range blockchain.blocks {
        fmt.Printf("ID: %d\n", block.ID)
        fmt.Printf("TimeStamp: %s\n", block.TimeStamp)
        fmt.Printf("Data: %s\n", block.Data)
        fmt.Printf("Previous Hash: %s\n", block.PrevHash)
        fmt.Printf("Hash: %s\n", block.Hash)
    }
}
```


## output
```bash
ID: 0   TimeStamp: 2024-07-17 16:48:14.299038516 +0300 EAT m=+0.000026802       Data: Genesis Block  Previous Hash:   Hash: 3849ac0e830963c34fe18b18ea3d945a 
ID: 1   TimeStamp: 2024-07-17 16:48:14.299127941 +0300 EAT m=+0.000116227       Data: second Block   Previous Hash: 3849ac0e830963c34fe18b18ea3d945a  Hash: fedb1d4405b0330e0e2049e69ccce0de 
ID: 2   TimeStamp: 2024-07-17 16:48:14.299135028 +0300 EAT m=+0.000123316       Data: third Block    Previous Hash: fedb1d4405b0330e0e2049e69ccce0de  Hash: dbcfa81389c07da02b9dbaa429d9eed1 
ID: 3   TimeStamp: 2024-07-17 16:48:14.299139949 +0300 EAT m=+0.000128236       Data: forth Block    Previous Hash: dbcfa81389c07da02b9dbaa429d9eed1  Hash: 62e57c727a2b48d23c28b1d6db08f85b 
```

This documentation provides a detailed explanation of each part of the blockchain code, ensuring a clear understanding of how the implementation works.