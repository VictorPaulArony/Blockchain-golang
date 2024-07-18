package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Document struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	TimeStamp string `json:"timestamp"`
	PrevHash  string `json:"prevhash"`
	Hash      string `json:"hash"`
}

type Blockchain struct {
	sync.Mutex
	Documents []Document
}

func CalcHash(doc Document) string {
	res := strconv.Itoa(doc.ID) + doc.Name + doc.TimeStamp + doc.PrevHash
	h := sha256.New()
	h.Write([]byte(res))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func CreateGenesis() Document {
	genesis := Document{0, "Genesis Document", time.Now().String(), "", ""}
	genesis.Hash = CalcHash(genesis)
	return genesis
}

func (bc *Blockchain) AddBlock(name string) string {
	bc.Lock()
	defer bc.Unlock()

	prevDocument := bc.Documents[len(bc.Documents)-1]
	newDocument := Document{
		ID:        prevDocument.ID + 1,
		Name:      name,
		TimeStamp: time.Now().String(),
		PrevHash:  prevDocument.Hash,
	}
	newDocument.Hash = CalcHash(newDocument)
	bc.Documents = append(bc.Documents, newDocument)
	return newDocument.Hash
}

func (bc *Blockchain) DocumentHistory() []Document {
	bc.Lock()
	defer bc.Unlock()
	return bc.Documents
}

var bc = &Blockchain{
	Documents: []Document{CreateGenesis()},
}

func addDocumentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		hash := bc.AddBlock(name)
		w.Write([]byte(hash))
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func getDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bc.DocumentHistory())
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func main() {
	fs := http.FileServer(http.Dir("styles"))
	http.Handle("/styles/", http.StripPrefix("/styles/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/add-document", addDocumentHandler)
	http.HandleFunc("/get-documents", getDocumentsHandler)

	fmt.Println("Server started at port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
