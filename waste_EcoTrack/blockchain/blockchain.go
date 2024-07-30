package blockchain

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"os"
	"strconv"
	"sync"
	"time"
)

type Collection struct {
	ID        int    `json:"id"`
	Data      string `json:"data"`
	TimeStamp string `json:"timestamp"`
	PrevHash  string `json:"prevhash"`
	Hash      string `json:"hash"`
}
type Reward struct {
	ID        int    `json:"id"`
	Token     string `json:"token"`
	TimeStamp string `json:"timestamp"`
	PrevHash  string `json:"prevhash"`
	Hash      string `json:"hash"`
}
type Blockchain struct {
	sync.Mutex
	collections []Collection
}

var fileName = "blocks.json"

func CreateHash(col Collection) string {
	res := strconv.Itoa(col.ID) + col.Data + col.TimeStamp + col.PrevHash + col.Hash
	hash := sha512.New()
	hash.Write([]byte(res))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

func GenerateGenesis() Collection {
	genesis := Collection{0, "Genesis Colloction", time.Now().String(), "", " "}
	genesis.Hash = CreateHash(genesis)
	return genesis
}
func (bc *Blockchain) AddBlock(data string) string {
	bc.Lock()
	defer bc.Unlock()

	prevBlock := bc.collections[len(bc.collections)-1]
	newCollection := Collection{
		ID:        prevBlock.ID + 1,
		Data:      data,
		TimeStamp: time.Now().String(),
		PrevHash:  prevBlock.Hash,
	}
	newCollection.Hash = CreateHash(newCollection)
	bc.collections = append(bc.collections, newCollection)
	return newCollection.Hash

}

//function to save blockchain to the json file
func (bc *Blockchain) SaveBlock() error {
	data, err := json.MarshalIndent(bc, "", " ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(fileName, data, 0o644); err != nil {
		return err
	}
	return nil
}

//function to load the blockchain from the json file
func (bc *Blockchain) LoadBlock() error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(nil) {
			bc.collections = []Collection{GenerateGenesis()}
			return bc.SaveBlock()
		}
		return err
	}
	err = json.Unmarshal(file, &bc)
	if err != nil {
		return err
	}
	return nil
}
