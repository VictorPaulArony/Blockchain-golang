package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"strconv"
	"time"
)

type Document struct {
	ID        int
	Name      string
	TimeStamp string
	PrevHash  string
	Hash      string
}
type Blockchain struct {
	Document []Document
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

func (bc *Blockchain) AddBlock(name string, hash string, reply *string) error {
	prevDocument := bc.Document[len(bc.Document)-1]
	neWDocument := Document{
		ID:        prevDocument.ID + 1,
		Name:      name,
		TimeStamp: time.Now().String(),
		PrevHash:  prevDocument.PrevHash,
		Hash:      hash,
	}
	bc.Document = append(bc.Document, neWDocument)
	*reply = neWDocument.Hash
	return nil
}

func (bc *Blockchain) DocumentHisttory(args int, reply *[]Document) error {
	*reply = bc.Document
	return nil
}

func StartServer(port string, blockchain *Blockchain) {
	rpc.Register(blockchain)
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalln("ERROR STARTING THE SERVER: ", err)
		return
	}
	defer l.Close()
	fmt.Println("Server started at port: ", port)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalln("CONNECTION ERROR: ", err)
			continue
		}
		go rpc.ServerConn(conn)
	}
}

func AddDocument(client *rpc.Client, name string, hash string) {
	var reply string
	err := client.Call("Blockchain.AddBlock", 0, &reply)
	if err != nil {
		log.Fatalln("ERROR ADDING DOCUMMENT: ", err)
		return
	}
	fmt.Println("ADDED DOCUMENT WITH HASH", reply)
}

func GetDocument(client *rpc.Client) []Document {
}
