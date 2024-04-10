package main

import (
	"net"
	"testing"
)

func sendRequest(conn net.Conn, msg string) (response string, err error) {
	conn.Write([]byte(msg))
	data, err := readData(conn)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func TestHandleClient(t *testing.T) {
	tests := []struct {
		request  string
		response string
	}{
		{
			request:  "+PING\r\n",
			response: "+PONG\r\n",
		},
		{
			request:  "+echo\r\nhey",
			response: "+hey\r\n",
		},
		{
			request:  "FOO",
			response: "Commands Array Expected",
		},
	}

	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	go HandleClient(server)

	for _, tt := range tests {
		response, err := sendRequest(client, tt.request)
		if err != nil {
			t.Errorf("Error on reading client data due to " + err.Error())
		}

		if response != tt.response {
			t.Errorf("Expected for %q but received %q", tt.response, response)
		}
	}
}
