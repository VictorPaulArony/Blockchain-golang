package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Initialize the Go module with the specified module name
	err := goModInit("module-name")
	if err != nil {
		log.Fatal(err)
	}

	// Set the endpoint for the Ethereum client
	endpoint := "https://mainnet.infura.io/v3/3599b276e8d4460eb1a195a19691085f"

	// Connect to the Ethereum client
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	// Get the latest block number
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Print the latest block number
	fmt.Println("Latest Block Number:", blockNumber)
}

// goModInit initializes the Go module with the specified module name
func goModInit(moduleName string) error {
	// Command: go mod init <module-name>
	cmd := exec.Command("go", "mod", "init", moduleName)

	// Run the command and check for any errors
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
