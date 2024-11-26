package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	server, _ := net.Listen("tcp", ":1234")

	defer server.Close()
	fmt.Println("server started at port 1234")
	for {
		conn, _ := server.Accept()

		fmt.Println("Client Connected at: ", conn.RemoteAddr())

		go ConnectionHandler(conn)
	}
}

// function for the client to dial from
func ConnectionHandler(conn net.Conn) {
	defer conn.Close()

	//handle outging sms
	go func() {
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			sms := input.Text() + "\n"
			conn.Write([]byte( sms))
		}
	}()
		
	//handle iincomming sms
	for {
		
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected")
			return
		}
		fmt.Println(response)
		fmt.Print("Server: ")
	}
}
