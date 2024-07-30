package database

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"os"
)

type Resident struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	UserId   string `json:"user_id"`
	Location string `json:"location"`
	Password string `json:"password"`
}

type Location struct {
	Building string `json:"building"`
	Region   string `json:"region"`
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
	Location  string `json:"location"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
}

var (
	FileName    = "resident-registration.json"
	StaffFile   = "staff-registration.json"
	RequestFile = "request.json"
	requests    []Request
	residents   []Resident
)

// SaveResident saves the resident data to the JSON file
func SaveResident(residents []Resident) error {
	data, err := json.MarshalIndent(residents, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(FileName, data, 0o644)
}

// LoadResident loads the resident data from the JSON file
func LoadResident() error {
	file, err := os.ReadFile(FileName)
	if err != nil {
		if os.IsNotExist(err) {
			residents = []Resident{}
		}
		return err
	}
	return json.Unmarshal(file, &residents)

}

// SaveStaff saves the staff data to the JSON file
func SaveStaff(staffs []Staff) error {
	data, err := json.MarshalIndent(staffs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(StaffFile, data, 0o644)
}

// LoadStaff loads the staff data from the JSON file
func LoadStaff() ([]Staff, error) {
	file, err := os.ReadFile(StaffFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Staff{}, nil
		}
		return nil, err
	}
	var staffs []Staff
	if err := json.Unmarshal(file, &staffs); err != nil {
		return nil, err
	}
	return staffs, nil
}

// CreateHash encrypts the user's password during registration
func CreateHash(password string) string {
	hash := sha512.New()
	hash.Write([]byte(password))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

// SaveRequest saves the requests to the JSON file
func SaveRequest() error {
	data, err := json.MarshalIndent(requests, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(RequestFile, data, 0o644)
}

// LoadRequest loads the requests from the JSON file
func LoadRequest() ([]Request, error) {
	file, err := os.ReadFile(RequestFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Request{}, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(file, &requests); err != nil {
		return nil, err
	}
	return requests, nil
}