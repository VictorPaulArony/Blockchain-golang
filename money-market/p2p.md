# Updated Code with P2P Networking:

Here’s a simplified version of how to modify the current code for P2P networking:
## Add a Peer Struct
``` go
type Peer struct {
	Address string `json:"address"`
}
```

## Add Networking Logic
``` go
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

var peers []Peer
var peerLock sync.Mutex

// AddPeer adds a new peer to the network
func AddPeer(address string) {
	peerLock.Lock()
	defer peerLock.Unlock()
	for _, peer := range peers {
		if peer.Address == address {
			return // Peer already exists
		}
	}
	peers = append(peers, Peer{Address: address})
}

// BroadcastTransaction sends a transaction to all connected peers
func BroadcastTransaction(tx Transaction) {
	for _, peer := range peers {
		go func(peer Peer) {
			conn, err := net.Dial("tcp", peer.Address)
			if err != nil {
				log.Printf("Failed to connect to peer %s: %v", peer.Address, err)
				return
			}
			defer conn.Close()

			txData, _ := json.Marshal(tx)
			conn.Write(txData)
		}(peer)
	}
}

// StartServer starts the P2P server
func StartServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	log.Printf("Server started on port %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

// Handle incoming connections
func handleConnection(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		var tx Transaction
		err := json.Unmarshal(scanner.Bytes(), &tx)
		if err != nil {
			log.Printf("Failed to parse transaction: %v", err)
			continue
		}
		log.Printf("Received transaction: %+v", tx)

		// Validate and add transaction
		ValidateAndAddTransaction(tx)
	}
}

// ValidateAndAddTransaction validates and adds a transaction
func ValidateAndAddTransaction(tx Transaction) {
	// Load users and validate sender and receiver
	users := LoadUsers()
	var senderUser, receiverUser *User
	for i := range users {
		if users[i].Wallet == tx.Sender {
			senderUser = &users[i]
		} else if users[i].Wallet == tx.Receiver {
			receiverUser = &users[i]
		}
	}

	if senderUser == nil || receiverUser == nil || senderUser.Balance < tx.Amount {
		log.Printf("Invalid transaction: %+v", tx)
		return
	}

	// Process transaction
	senderUser.Balance -= tx.Amount
	receiverUser.Balance += tx.Amount
	SaveUsers(users)

	// Save transaction
	transactions := LoadTransactions()
	transactions = append(transactions, tx)
	SaveTransactions(transactions)

	log.Printf("Transaction added: %+v", tx)
}
```

## Update TransactionHandler to Broadcast Transactions

Modify the TransactionHandler to broadcast transactions to peers:

``` go
func TransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	sender := r.FormValue("sender_wallet")
	receiver := r.FormValue("receiver_wallet")
	amountStr := r.FormValue("amount")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	users := LoadUsers()
	var senderUser, receiverUser *User
	for i := range users {
		if users[i].Wallet == sender {
			senderUser = &users[i]
		} else if users[i].Wallet == receiver {
			receiverUser = &users[i]
		}
	}

	if senderUser == nil || receiverUser == nil {
		http.Error(w, "Invalid wallet address", http.StatusBadRequest)
		return
	}
	if senderUser.Balance < amount {
		http.Error(w, "Insufficient balance", http.StatusBadRequest)
		return
	}

	// Update balances
	senderUser.Balance -= amount
	receiverUser.Balance += amount
	SaveUsers(users)

	// Record transaction
	transaction := Transaction{
		ID:        uuid.New().String(),
		Sender:    sender,
		Receiver:  receiver,
		Amount:    amount,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	transactions := LoadTransactions()
	transactions = append(transactions, transaction)
	SaveTransactions(transactions)

	// Broadcast the transaction to peers
	BroadcastTransaction(transaction)

	// Redirect back to the dashboard
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
```
## To Run the P2P Network:

    Start Nodes:
        Run multiple instances of your application, each listening on a different port.
        For example, go run main.go --port=8081, go run main.go --port=8082, etc.

    Connect Nodes:
        Use an API or configuration file to connect nodes.
        For example, one node connects to another by adding the peer's address.

    Test Transactions:
        Make a transaction from one node. The transaction will be broadcast to all connected peers and validated.

    Synchronization:
        Implement a method to synchronize the blockchain state across peers when a new node joins the network.

This setup provides the foundation for a P2P transaction network. You can extend it with features like consensus mechanisms and advanced synchronization algorithms for scalability and fault tolerance.