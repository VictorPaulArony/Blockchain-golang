# Exploring Basic Wallet Implementation in Golang

1. **Introduction**
This document aims to help those who are new to blockchain wallet creatinon process become familiar with the basic wallet implementation using Go language. If you use it in combination with the content of this document, you will be able to create your wallet from scratch and understand the underlying technologies better.
The wallet stores publick and private keys and interacts with various blockchain networks that allow the users to send or receive digital currencies or tokens within the network. The wallet doesn’t store the coins themselves but keeps track of balances of the addresses.

## Key Features of the Basic Wallet:
1. **Generate Keys**: The wallet generates the public and private keys using cryptographic algorithms i.e sha256.
2. **Transaction Management**: It allows users to create and manage transactions.
3. **Balance Tracking**: The wallet helps users to track their balance based on transactions they have done.

## Getting Started
### Key Components of the Wallet
1. **Wallet and Address Structures**
Create a sturct for the Wallet and the Address i.e the wallet is represented by a Wallet struct that contains a map of addresses. Each address is represented by the Address struct, which holds the private key, public key, and balance.

```go
type Wallet struct {
    Adresses map[string]*Address
}
type Address struct {
    Privatekey *ecdsa.PrivateKey
    PublicKey  *ecdsa.PublicKey
    Balance    float64
}
```
2. **Creating a New Wallet**

The CreateWallet function initializes a new wallet instance.
```go
func CreateWallet() *Wallet {
    return &Wallet{Adresses: make(map[string]*Address)}
}
```

3. **Key Generation**: Use the crypto/rand and ECDSA (Elliptic Curve Digital Signature Algorithm) packages to generate a secure random private key. The public key is derived from the private key using elliptic curve cryptography.
```go
func (w *Wallet) CreateAddress() string {
    privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        log.Fatalln(privateKey)
    }

    publicKey := privateKey.PublicKey
    publicKeyBytes := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
    add := sha256.Sum256(publicKeyBytes)
    address := hex.EncodeToString(add[:])

    w.Adresses[address] = &Address{
        Privatekey: privateKey,
        PublicKey:  &publicKey,
        Balance:    0,
    }
    return address
}
```
4. **Transferring Funds**

The Transfer method allows funds to be transferred between addresses. It first checks if both the sender and receiver addresses exist, then updates their balances according to the actions of either sending or receiving

```go
func (w *Wallet) Transfer(to, from string, amount float64) error {
    addFrom, exist := w.Adresses[from]
    if !exist {
        log.Fatalln(from)
    }

    addTo, exist := w.Adresses[to]
    if !exist {
        log.Fatalln(to)
    }
    
    if addFrom.Balance < amount {
        return fmt.Errorf("insufficient balance")
    }

    addFrom.Balance -= amount
    addTo.Balance += amount
    return nil
}
```

5. **Getting Balance**

The GetBalance method allows the user to check the balance of a specific address, by retrieving the balance from the wallet.

```go
func (w *Wallet) GetBalance(address string) (float64, error) {
    add, exist := w.Adresses[address]
    if !exist {
        return 0, fmt.Errorf("address not found")
    }
    return add.Balance, nil
}
```
6. **Putting It All Together**

The main function that demonstrate the wallet’s capabilities based on the functions and methods above:

```go
func main() {
    wallet := CreateWallet()

    addr1 := wallet.CreateAddress()
    addr2 := wallet.CreateAddress()

    fmt.Printf("THE FIRST ADDRESS: %s \n", addr1)
    fmt.Printf("THE SECOND ADDRESS: %s\n", addr2)

    wallet.Adresses[addr1].Balance = 136

    fmt.Printf("INITIAL BALANCE FOR FIRST ADDRESS: %f\n", wallet.Adresses[addr1].Balance)
    fmt.Printf("INITIAL BALANCE FOR SECOND ADDRESS: %f\n", wallet.Adresses[addr2].Balance)

    err := wallet.Transfer(addr1, addr2, 25)
    if err != nil {
        log.Fatalln(err)
    }
    
    fmt.Printf("CURRENT BALANCE FOR FIRST ADDRESS: %f\n", wallet.Adresses[addr1].Balance)
    fmt.Printf("CURRENT BALANCE FOR SECOND ADDRESS: %f\n", wallet.Adresses[addr2].Balance)
}
```

**Expected Output**

When you run the program, you should see the initial balances for both addresses and the updated balances after the transfer:

```bash
go run main.go
```

Output
```bash
THE FIRST ADDRESS: <address_1> 
THE SECOND ADDRESS: <address_2>
INITIAL BALANCE FOR FIRST ADDRESS: 136.000000
INITIAL BALANCE FOR SECOND ADDRESS: 0.000000
CURRENT BALANCE FOR FIRST ADDRESS: 111.000000
CURRENT BALANCE FOR SECOND ADDRESS: 25.000000
```

## Challenges and Lessons Learned

Throughout this project, I encountered several challenges, including understanding ECDSA key generation and managing balance updates accurately. This experience has reinforced the importance of security in blockchain applications and the necessity of robust error handling to ensure smooth operations.

## Conclusion

Building a basic blockchain wallet in Go has been an enlightening experience that has deepened my understanding of blockchain fundamentals. This project serves as a foundation for more complex applications and showcases the capabilities of the Go programming language in blockchain development. I encourage you to explore and expand upon this code, perhaps integrating more features like transaction history or additional security measures. Happy coding!
full code at https://github.com/VictorPaulArony/Blockchain-golang/tree/main/basic-wallet