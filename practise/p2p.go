package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"sync"
	"time"
)

type Transaction struct {
	ID        int
	Sender    string
	Receiver  string
	TimeStamp string
	Amount    float64
	Signature string
}

type Mempool struct {
	Transaction []Transaction
	mu          sync.Mutex
}

type Block struct {
	ID          int
	Transaction []Transaction
	TimeStamp   string
	MerkleRoot  string
	PrevHash    string
	Hash        string
	Nonce       int
}

type Blockchain struct {
	Blocks []Block
	mu     sync.Mutex
}

type Address struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Balance    float64
}

type Wallet struct {
	Address map[string]*Address
}

// function to generate random value for the salt
func Salt() (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

// function to generate the hashing functions
func Hash(data, salt string) string {
	hash := sha256.Sum256([]byte(data + salt))
	return hex.EncodeToString(hash[:])
}

// function to create the wallate
func CreateWallet() Wallet {
	return Wallet{make(map[string]*Address)}
}

// function to create the private and public key of the address
func (w *Wallet) CreateAddress() string {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Println(err)
	}

	publicKey := privateKey.PublicKey
	publicKeyBytes := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
	AdrHash := sha256.Sum256(publicKeyBytes)
	address := hex.EncodeToString(AdrHash[:])

	w.Address[address] = &Address{
		PrivateKey: privateKey,
		PublicKey:  &publicKey,
		Balance:    0,
	}
	return address
}

// function to create a mempool for the transaction
func CreateMempool() Mempool {
	return Mempool{Transaction: []Transaction{}}
}

// function create a transaction
func (w *Wallet) CreateTransaction(from, to string, amount float64) error {
	mp := &Mempool{}
	mp.mu.Lock()
	defer mp.mu.Unlock()

	sender, exist := w.Address[from]
	if !exist {
		log.Fatalln(from)
	}

	receiver, exist := w.Address[to]
	if !exist {
		log.Fatalln(to)
	}

	if w.Address[from].Balance < amount {
		log.Println("Insufficient Funds In the Wallet")
		return fmt.Errorf("insufficient funds")
	}

	// Check if there are previous transactions
	var prevTxID int
	if len(mp.Transaction) > 0 {
		prevTx := mp.Transaction[len(mp.Transaction)-1]
		prevTxID = prevTx.ID
	} else {
		prevTxID = 0
	}

	transaction := &Transaction{
		ID:        prevTxID + 1,
		Sender:    from,
		Receiver:  to,
		TimeStamp: time.Now().String(),
		Amount:    amount,
	}

	err := w.SignTransaction(transaction)
	if err != nil {
		log.Println(err)
	}

	if !w.IsValidTransaction(transaction) {
		log.Println("Invalid Transaction")
	}

	sender.Balance -= amount
	receiver.Balance += amount

	mp.Transaction = append(mp.Transaction, *transaction)
	return nil
}

// Fuunction to sign the transaction
func (w *Wallet) SignTransaction(tx *Transaction) error {
	addr, exist := w.Address[tx.Sender]
	if !exist {
		log.Fatalln(tx.Sender)
	}

	salt, err := Salt()
	if err != nil {
		log.Println(err)
	}

	hash := Hash((strconv.Itoa(tx.ID) + tx.Sender + tx.Receiver + tx.TimeStamp + strconv.Itoa(int(tx.Amount))), salt)
	r, s, err := ecdsa.Sign(rand.Reader, addr.PrivateKey, []byte(hash))
	if err != nil {
		log.Println(err)
	}

	signature := append(r.Bytes(), s.Bytes()...)

	tx.Signature = hex.EncodeToString(signature)

	return nil
}

// function to verify the transaction
func (w *Wallet) IsValidTransaction(tx *Transaction) bool {
	addr, exist := w.Address[tx.Sender]
	if !exist {
		log.Println(tx.Sender)
	}

	salt, err := Salt()
	if err != nil {
		log.Println(err)
	}
	hash := Hash(strconv.Itoa(tx.ID)+tx.Receiver+tx.Sender+tx.TimeStamp+strconv.Itoa(int(tx.Amount)), salt)
	signature, err := hex.DecodeString(tx.Signature)
	if err != nil {
		log.Println(err)
	}

	r := big.Int{}
	s := big.Int{}

	signLen := len(signature)

	r.SetBytes(signature[:(signLen / 2)])
	s.SetBytes(signature[(signLen / 2):])
	return ecdsa.Verify(&addr.PrivateKey.PublicKey, []byte(hash[:]), &r, &s)
}

// fnction to form the merkel tree root for a transaction
func MerkleRoots(transaction []Transaction) string {
	if len(transaction) == 0 {
		return ""
	}

	var hashes []string

	for _, tx := range transaction {
		salt, err := Salt()
		if err != nil {
			log.Println(err)
		}
		hash := Hash(strconv.Itoa(tx.ID)+tx.Receiver+tx.Sender+tx.TimeStamp+strconv.Itoa(int(tx.Amount)), salt)
		hashes = append(hashes, hash)
	}

	for len(hashes) > 1 {
		var newHashes []string
		for i := 0; i < len(hashes); i += 2 {
			if i+1 < len(hashes) {
				newHashes = append(newHashes, Hash(hashes[i]+hashes[i+1], "")) // Fixed hashing logic for combining pairs
			} else {
				newHashes = append(newHashes, hashes[i])
			}
		}
		hashes = newHashes
	}
	return hashes[0]
}

// functon to create a hashing function fo the block
func (b *Block) CreateHash() string {
	salt, err := Salt()
	if err != nil {
		log.Println(err)
	}

	hash := strconv.Itoa(b.ID) + b.MerkleRoot + b.TimeStamp + b.PrevHash + b.Hash + strconv.Itoa(b.Nonce)
	// var tx *Transaction
	// hash += strconv.Itoa(tx.ID) + tx.Receiver + tx.Sender + tx.Signature + tx.TimeStamp + strconv.Itoa(int(tx.Amount))

	return Hash(hash, salt)
}

// function to create a genesis block of the blockchain
func CreateGenesis() Block {
	var transaction []Transaction
	genesis := Block{
		ID:          0,
		Transaction: transaction,
		TimeStamp:   time.Now().String(),
		PrevHash:    "",
	}

	for {
		genesis.Hash = genesis.CreateHash()
		if !IsValidHash(genesis.Hash) {
			break
		}
		genesis.Nonce++
	}

	return genesis
}

// function to create the concensus ie. PoW
func IsValidHash(hash string) bool {
	return hash[:5] == "00000"
}

// function to craeta a new blockchain
func CreateBlockchain() Blockchain {
	genesis := CreateGenesis()
	return Blockchain{Blocks: []Block{genesis}}
}

// function to create and add blocks to the blockchain
func (bc *Blockchain) AddBlock() {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	var mp *Mempool
	mp = &Mempool{}

	mp.mu.Lock()
	defer mp.mu.Unlock()

	// retrieve transaction from mempool
	transaction := mp.Transaction
	if len(transaction) == 0 {
		log.Println("N Transaction to Mine")
		return
	}

	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := Block{
		ID:          prevBlock.ID + 1,
		Transaction: transaction,
		TimeStamp:   time.Now().String(),
		PrevHash:    prevBlock.Hash,
	}

	newBlock.MerkleRoot = MerkleRoots(transaction)

	for {
		newBlock.Hash = newBlock.CreateHash()
		if !IsValidHash(newBlock.Hash) {
			break
		}
		newBlock.Nonce++
	}
	if !bc.IsValidBlock(newBlock) {
		log.Println("Invalid Block")
		return
	}
	bc.Blocks = append(bc.Blocks, newBlock)

	// clear mempool for next transaction block
	mp.Transaction = []Transaction{}
}

// function to validate the blocks before adding to the blockchain
func (bc *Blockchain) IsValidBlock(block Block) bool {
	if len(block.Hash) == 0 {
		return false
	}

	prevBlock := bc.Blocks[len(bc.Blocks)-1]

	if prevBlock.Hash == block.PrevHash {
		return true
	}
	return false
}

// main function to run the project
func main() {
	wallet := CreateWallet()

	blockchain := CreateBlockchain()

	address1 := wallet.CreateAddress()
	address2 := wallet.CreateAddress()

	fmt.Printf("Address 1: %s\n", address1)
	fmt.Printf("Address 2: %s\n", address2)

	wallet.Address[address1].Balance = 100

	fmt.Printf("Initial Balance of Address 1: %.2f\n", wallet.Address[address1].Balance)
	fmt.Printf("Initial Balance of Address 2: %.2f\n", wallet.Address[address2].Balance)

	err := wallet.CreateTransaction(address1, address2, 50)
	if err != nil {
		log.Fatal(err)
	}

	blockchain.AddBlock()

	fmt.Printf("Balance of Address 1 after transfer: %.2f\n", wallet.Address[address1].Balance)
	fmt.Printf("Balance of Address 2 after transfer: %.2f\n", wallet.Address[address2].Balance)

	// Print the blockchain
	for _, block := range blockchain.Blocks {
		fmt.Printf("Block ID: %d\n", block.ID)
		fmt.Printf("Timestamp: %s\n", block.TimeStamp)
		fmt.Printf("Previous Hash: %s\n", block.PrevHash)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Printf("Transactions:\n")
		for _, tx := range block.Transaction {
			fmt.Printf("  %s -> %s: %f\n", tx.Sender, tx.Receiver, tx.Amount)
		}
		fmt.Println()
	}
}
