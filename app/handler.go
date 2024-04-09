package main

import (
	"fmt"
	"net"
	"strings"
)

func HandleClient(conn net.Conn) {
	for {
		data, err := readData(conn)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Received", data)

		response := handleRequest(data)
		conn.Write(response.Bytes())
	}
}

func handleRequest(data []byte) RValue {
	rValue := NewRValue(data[0])
	_, err := rValue.Parse()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var response RValue

	// TODO create interface Command.Execute() and create different commsnd types to handle own logic
	switch rValue.Value.(type) {
	case ArrayValue{}:
		response = executeCommand(strings.ToLower(rValue.Value[0].Value))
	default:
		response = make(ErrorValue{Value: "ERR: Command Expected"})
	}

	fmt.Println("We are about to send", response.Value)
	return response
}

func executeCommand(cmd string) RValue {
	var result RValue

	switch cmd {
	case "ping":
		make(SimpleStringValue{Value: "PONG"})
	case "echo":
		message := rValue[1].Value
		make(StringValue{Value: message})
	}

	return result
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
