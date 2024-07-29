package database

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"log"
	"os"
)

type Resident struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	UserId   string `json:"user_id"`
	Location `json:"location"`
	Password string `json:"password"`
}
type Location struct {
	Building string
	Region   string
}
type Staff struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Location string `json:"location"`
	Password string `json:"password"`
}
type Request struct {
	ID        int    `json:"id"`
	UserId    string `json:"user_id"`
	Nature    string `json:"nature"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
}

var FileName = "resident-registration.json"
var requests []Request

//function to save the resident data to the json file
func SaveResident() {
	data, err := os.ReadFile(FileName)
	if err != nil {
		log.Fatal("ERROR OPENING THE FILE: ", err)
	}
	err = os.WriteFile(FileName, data, 0o644)
	if err != nil {
		log.Fatal("FILE DOES NOT EXIT")
	}
}

//function to save the staff data to the json
func SaveStaff() {
	data, err := os.ReadFile(FileName)
	if err != nil {
		log.Fatal("ERROR OPENING THE FILE: ", err)
	}
	err = os.WriteFile(FileName, data, 0o644)
	if err != nil {
		log.Fatal("FILE DOES NOT EXIT")
	}
}

//function to encrypt the users passwords during registration
func CreateHash(password string) string {
	hash := sha512.New()
	hash.Write([]byte(password))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}
func SaveRequest() error {
	data, err := json.MarshalIndent(requests, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile("request.json", data, 0o644)
}
func LoadRequest() ([]Request, error) {
	file, err := os.ReadFile("request.json")
	if err != nil {
		if os.IsNotExist(nil) {
			return []Request{}, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(file, &requests); err != nil {
		return nil, err
	}
	return requests, nil
}
