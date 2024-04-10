package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	port := flag.String("port", "6379", "Redis port [Default: 6379]")
	flag.Parse()

	tcp, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", *port))
	if err != nil {
		fmt.Printf("Failed to bind to port %v\n", *port)
		os.Exit(1)
	}

	defer tcp.Close()

	for {
		conn, err := tcp.Accept()
		conn.SetDeadline(time.Now().Add(20 * time.Second))

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			return
		}

		go HandleClient(conn)

	}
}
