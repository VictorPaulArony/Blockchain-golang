package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// Block structure
type Block struct {
	Index     int    // Position of the data record in the blockchain
	Timestamp string // Time of the block creation
	Data      string // Data stored in the block
	PrevHash  string // Hash of the previous block
	Hash      string // Hash of the current block
}

// Blockchain is a slice of Blocks
var (
	Blockchain []Block
	mutex      = &sync.Mutex{} // Mutex for safe concurrent access
)

// Network peers
var peers = []string{"localhost:9002", "localhost:9003"}

// Compute the SHA-256 hash for a block
func calculateHash(block Block) string {
	record := fmt.Sprintf("%d%s%s%s", block.Index, block.Timestamp, block.Data, block.PrevHash)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Generate a new block based on the previous block
func generateBlock(oldBlock Block, data string) Block {
	newBlock := Block{
		Index:     oldBlock.Index + 1,
		Timestamp: time.Now().String(),
		Data:      data,
		PrevHash:  oldBlock.Hash,
		Hash:      "",
	}

	newBlock.Hash = calculateHash(newBlock)
	return newBlock
}

// Validate if the block is valid by checking the hashes
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

// Replace the chain if a longer valid chain is found
func replaceChain(newBlocks []Block) {
	mutex.Lock()
	defer mutex.Unlock()

	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

// Handle incoming connections from peers
func handleConnection(conn net.Conn) {
	defer conn.Close()

	io.WriteString(conn, "Enter data for the block:\n")

	// Read the input from the connection
	scanData := make([]byte, 1024)
	n, err := conn.Read(scanData)
	if err != nil {
		log.Println("Error reading from connection:", err)
		return
	}

	data := string(scanData[:n])

	// Create a new block with the input data
	mutex.Lock()
	newBlock := generateBlock(Blockchain[len(Blockchain)-1], data)
	mutex.Unlock()

	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		mutex.Lock()
		Blockchain = append(Blockchain, newBlock)
		mutex.Unlock()
	}

	// Broadcast the new block to peers
	for _, peer := range peers {
		sendBlockToPeer(peer, newBlock)
	}

	io.WriteString(conn, "\nNew block mined: \n")
	bytes, err := json.MarshalIndent(newBlock, "", "  ")
	if err != nil {
		log.Println("Error marshaling block:", err)
		return
	}
	io.WriteString(conn, string(bytes)+"\n")
}

// Function to send block to a peer
func sendBlockToPeer(peerAddress string, block Block) {
	conn, err := net.Dial("tcp", peerAddress)
	if err != nil {
		log.Println("Error connecting to peer:", err)
		return
	}
	defer conn.Close()

	bytes, err := json.Marshal(block)
	if err != nil {
		log.Println("Error marshaling block:", err)
		return
	}

	conn.Write(bytes)
}

// Listen for incoming connections on a specific port
func startServer(port string) {
	server, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	log.Println("Listening on port:", port)

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

// Start the node, connect to peers, and listen for blocks
func startNode(port string) {
	go startServer(port)

	// Connect to peers to sync blockchain
	for _, peer := range peers {
		conn, err := net.Dial("tcp", peer)
		if err != nil {
			log.Println("Unable to connect to peer", peer)
			continue
		}
		defer conn.Close()

		// Request the peer's blockchain
		fmt.Fprintln(conn, "GET_BLOCKCHAIN")
		handleConnection(conn)
	}
}

func main() {
	// Create the Genesis block
	genesisBlock := Block{0, time.Now().String(), "Genesis Block", "", ""}
	genesisBlock.Hash = calculateHash(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	// Start the node on port 9001
	startNode("9002")

	// Keep the main function running
	select {}
}