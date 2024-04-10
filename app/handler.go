package main

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"strings"
)

func HandleClient(conn net.Conn) {
	defer conn.Close()

	var response RValue

	for {
		data, err := readData(conn)
		if err != nil {
			handleError(conn, err)
			return
		}

		response, err = handleRequest(data)
		if err != nil {
			handleError(conn, err)
			return
		}

		conn.Write(response.Bytes())
	}
}

func handleError(conn net.Conn, err error) {
	fmt.Fprintln(os.Stderr, err.Error())

	response := &ErrorValue{value: err.Error()}

	conn.Write(response.Bytes())
}

func handleRequest(data []byte) (RValue, error) {
	request, err := NewRValue(data[0])
	if err != nil {
		return nil, err
	}

	_, err = request.Parse(data)
	if err != nil {
		return nil, err
	}

	var response RValue

	// TODO create interface Command.Execute() and create different commsnd types to handle own logic
	switch reflect.TypeOf(request) == reflect.TypeOf(&ArrayValue{}) {
	case true:
		response = executeCommand(request)
	default:
		response = &ErrorValue{value: "Commands Array Expected"}
	}

	return response, nil
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
		return &ErrorValue{value: fmt.Sprintf("Invalid command %q", cmd)}
	}
}

func readData(conn net.Conn) ([]byte, error) {
	data := make([]byte, 1024)
	size, err := conn.Read(data)
	if err != nil {
		return nil, fmt.Errorf("Error on reading from socket due to %w", err)
	}
	return data[:size], nil
}
