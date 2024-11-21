package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	client, _ := net.Dial("tcp", "192.168.89.77:1234")

	defer client.Close()
	fmt.Println("Connected to the server")

	for {
		fmt.Print("client: ")

		input := bufio.NewReader(os.Stdin)
		sms, _ := input.ReadString('\n')

		client.Write([]byte(sms))
		response, _ := bufio.NewReader(client).ReadString('\n')

		fmt.Println("Server: ", response)
	}
}
