package main

import (
	"fmt"
	"net"
	"reflect"
	"strings"
)

func HandleClient(conn net.Conn) {
	for {
		data, err := readData(conn)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		response := handleRequest(data)
		conn.Write(response.Bytes())
	}
}

func handleRequest(data []byte) RValue {
	request := NewRValue(data[0])
	_, err := request.Parse(data)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	var response RValue

	// TODO create interface Command.Execute() and create different commsnd types to handle own logic
	switch reflect.TypeOf(request).Name() {
	case "ArrayValue":
		response = executeCommand(request)
	default:
		response = &ErrorValue{value: "ERR: Commands Array Expected"}
	}

	return response
}

func executeCommand(req RValue) RValue {
	commands := req.Value().([]string)
	cmd := strings.ToLower(commands[0])

	switch cmd {
	case "ping":
		return &SimpleStringValue{value: "PONG"}
	case "echo":
		return &StringValue{value: commands[1]}
	default:
		return &ErrorValue{value: fmt.Sprintf("Invalid command %s", cmd)}
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
