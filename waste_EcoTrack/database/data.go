package database

import (
	"crypto/sha512"
	"encoding/hex"
	"log"
	"os"
)

type Resident struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Location string `json:"location"`
	Password string `json:"password"`
}
type Staff struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Location string `json:"location"`
	Password string `json:"password"`
}

var FileName = "resident-registration.json"

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
