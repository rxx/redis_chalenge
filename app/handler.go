package main

import (
	"fmt"
	"net"
)

func HandleClient(conn net.Conn) {
	for {
		data, err := readData(conn)
		if err != nil {
			return
		}

		fmt.Println("Received", data)

		// parseData(data)

		conn.Write([]byte(formatSimpleString("PONG")))
	}
}

func readData(conn net.Conn) ([]byte, error) {
	data := make([]byte, 1024)
	size, err := conn.Read(data)
	if err != nil {
		fmt.Println("Error on reading from socket", err.Error())
		return nil, err
	}
	return data[:size], nil
}
