package main

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/sha256"
    "encoding/json"
    "encoding/hex"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "sync"
)

type User struct {
    Username   string
    Email      string
    Phone      string
    Address    string
    PrivateKey string
}

type Wallet struct {
    Users    map[string]*User
    mu       sync.Mutex
    Filename string
}

// CreateWallet initializes a new wallet
func CreateWallet(filename string) *Wallet {
    return &Wallet{Users: make(map[string]*User), Filename: filename}
}

// SaveToJSON saves users to a JSON file
func (w *Wallet) SaveToJSON() error {
    w.mu.Lock()
    defer w.mu.Unlock()

    data, err := json.Marshal(w.Users)
    if err != nil {
        return err
    }

    return ioutil.WriteFile(w.Filename, data, 0644)
}

// LoadFromJSON loads users from a JSON file
func (w *Wallet) LoadFromJSON() error {
    w.mu.Lock()
    defer w.mu.Unlock()

    // Check if the file exists
    if _, err := os.Stat(w.Filename); os.IsNotExist(err) {
        return nil // File doesn't exist, return nil to skip loading
    }

    data, err := ioutil.ReadFile(w.Filename)
    if err != nil {
        return err
    }

    return json.Unmarshal(data, &w.Users)
}

// GenerateWallet generates a new wallet for the user
func (w *Wallet) GenerateWallet(username, email, phone string) (string, error) {
    w.mu.Lock()
    defer w.mu.Unlock()

    // Check if user already exists
    for _, user := range w.Users {
        if user.Email == email {
            return "", fmt.Errorf("user with email %s already exists", email)
        }
    }

    privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        return "", err
    }

    publicKey := privateKey.PublicKey
    publicKeyBytes := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
    add := sha256.Sum256(publicKeyBytes)
    address := hex.EncodeToString(add[:])

    user := &User{
        Username:   username,
        Email:      email,
        Phone:      phone,
        Address:    address,
        PrivateKey: hex.EncodeToString(elliptic.Marshal(elliptic.P256(), privateKey.X, privateKey.Y)), // Store the private key
    }

    w.Users[email] = user
    err = w.SaveToJSON()
    if err != nil {
        return "", err
    }

    return address, nil
}

// RegistrationHandler handles user registration
func (w *Wallet) RegistrationHandler(wr http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(wr, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    var userDetails struct {
        Username string `json:"username"`
        Email    string `json:"email"`
        Phone    string `json:"phone"`
    }

    err := json.NewDecoder(r.Body).Decode(&userDetails)
    if err != nil {
        http.Error(wr, err.Error(), http.StatusBadRequest)
        return
    }

    address, err := w.GenerateWallet(userDetails.Username, userDetails.Email, userDetails.Phone)
    if err != nil {
        http.Error(wr, err.Error(), http.StatusBadRequest)
        return
    }

    fmt.Fprintf(wr, "Wallet created successfully! Address: %s", address)
}

// ServeHTML serves the HTML page for registration
func ServeHTML(wr http.ResponseWriter, r *http.Request) {
    html := `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Wallet Registration</title>
    </head>
    <body>
        <h1>Register for Wallet</h1>
        <form id="registerForm">
            <input type="text" id="username" placeholder="Enter your username" required />
            <input type="email" id="email" placeholder="Enter your email" required />
            <input type="text" id="phone" placeholder="Enter your phone number" required />
            <button type="submit">Register</button>
        </form>
        <h3>Message</h3>
        <p id="message"></p>

        <script>
            document.getElementById("registerForm").onsubmit = async (e) => {
                e.preventDefault();
                const username = document.getElementById("username").value;
                const email = document.getElementById("email").value;
                const phone = document.getElementById("phone").value;

                const response = await fetch("/register", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify({ username, email, phone })
                });

                const message = await response.text();
                document.getElementById("message").innerText = message;
            }
        </script>
    </body>
    </html>
    `
    wr.Header().Set("Content-Type", "text/html")
    wr.Write([]byte(html))
}

func main() {
    wallet := CreateWallet("/users.json")
    
    // Load existing users from JSON file
    err := wallet.LoadFromJSON()
    if err != nil {
        log.Println("Error loading users from JSON:", err)
    }

    // Serve the HTML page
    http.HandleFunc("/", ServeHTML)

    // API route for user registration
    http.HandleFunc("/register", wallet.RegistrationHandler)

    fmt.Println("Server started at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
