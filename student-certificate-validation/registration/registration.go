package registration

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
)

type Register struct {
	Name     string `json:"name"`
	RegNo    string `json:"regno"`
	Course   string `json:"course"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"pass"`
}
type Admin struct {
	AdminId    string `json:"adminid"`
	Name       string `json:"name"`
	Department string `json:"department"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Password   string `json:"pass"`
}
type CertificateRequest struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	RegNo     string `json:"regno"`
	Course    string `json:"course"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
}

type Certificate struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	RegNo     string `json:"regno"`
	Course    string `json:"course"`
	CreatedAt string `json:"created_at"`
	Hash      string `json:"hash"`
}

var (
	students []Register
	requests []CertificateRequest
	admins   []Admin
	// certificates     []Certificate
	registrationFile = "users.json"
	requestsFile     = "requests.json"
	adminFile        = "admins.json"
)

// SaveCertificates function saves the certificates to a storage (e.g., a file or database)
func SaveCertificates(certificates []Certificate) error {
	data, err := json.MarshalIndent(certificates, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("certificates.json", data, 0o644)
}

// LoadCertificates function loads the certificates from a storage (e.g., a file or database)
func LoadCertificates() ([]Certificate, error) {
	file, err := os.ReadFile("certificates.json")
	if err != nil {
		if os.IsNotExist(err) {
			return []Certificate{}, nil // Return an empty slice if the file doesn't exist
		}
		return nil, err
	}
	var certificates []Certificate
	if err := json.Unmarshal(file, &certificates); err != nil {
		return nil, err
	}
	return certificates, nil
}

// SaveRequests saves the certificate requests to a file
func SaveRequests(reqs []CertificateRequest) error {
	data, err := json.MarshalIndent(reqs, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(requestsFile, data, 0o644)
}

// LoadRequests loads the certificate requests from a file
func LoadRequests() error {
	file, err := os.ReadFile(requestsFile)
	if err != nil {
		if os.IsNotExist(err) {
			requests = []CertificateRequest{}
			return nil
		}
		return err
	}
	return json.Unmarshal(file, &requests)
}

// AddStudent saves the student data to a file
func AddStudent(students []Register) error {
	data, err := json.MarshalIndent(students, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(registrationFile, data, 0o644)
}

func LoadStudents() error {
	file, err := os.ReadFile(registrationFile)
	if err != nil {
		if os.IsNotExist(err) {
			students = []Register{}
			return nil
		}
		return err
	}
	return json.Unmarshal(file, &students)
}

// AddAdmin saves the admin data to a file
func AddAdmin(admins []Admin) error {
	data, err := json.MarshalIndent(admins, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(adminFile, data, 0o644)
}

// LoadAdmins loads the admin data from a file
func LoadAdmins() error {
	file, err := os.ReadFile(adminFile)
	if err != nil {
		if os.IsNotExist(err) {
			admins = []Admin{}
			return nil
		}
		return err
	}
	return json.Unmarshal(file, &admins)
}

// HashPassword hashes the password
func HashPassword(pass string) string {
	hash := sha256.New()
	hash.Write([]byte(pass))
	return hex.EncodeToString(hash.Sum(nil))
}
