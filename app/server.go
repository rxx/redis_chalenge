package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Listen port 6379")
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()

	for {
		fmt.Println("Wait for the client...")
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		handleClient(conn)

	}
}

func handleClient(conn net.Conn) {
	for {
		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err != nil {
			fmt.Println("Error on reading from socket", err.Error())
			os.Exit(1)
		}

		fmt.Printf("Received: %s with %d bytes\n", data[:n], n)
		conn.Write(formatString("PONG"))
	}
}

func formatString(response string) []byte {
	result := fmt.Sprintf("+%s\r\n", response)

	fmt.Println("Response to the client with", result)

	return []byte(result)
}
