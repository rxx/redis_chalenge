package main

import (
	"fmt"
	"net"
	"strings"
)

// role:master
// connected_slaves:0
// master_replid:8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb
// master_repl_offset:0
// second_repl_offset:-1
// repl_backlog_active:0
// repl_backlog_size:1048576
// repl_backlog_first_byte_offset:0
// repl_backlog_histlen:

func ReplicationInfoToString() string {
	replicationInfo := map[string]string{
		"role": nodeConfig.Role,
		// 40 sized random string
		"master_replid":      "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		"master_repl_offset": "0",
	}

	var result strings.Builder

	for key, value := range replicationInfo {
		rValue := &StringValue{value: fmt.Sprintf("%s:%s", key, value)}
		result.WriteString(rValue.String())
	}

	return result.String()
}

func StartReplica() {
	fmt.Printf("StartReplica: MasterConfig %s\n", nodeConfig.MasterConfig)
	if nodeConfig.MasterConfig == nil {
		return
	}

	SendMessageToMaster("PING")
}

// Try 5 times to send message
func SendMessageToMaster(msg string) {
	fmt.Printf("Send message to master: %s\n", msg)

	for repeat := 5; repeat > 0; repeat-- {
		if err := writeMessage("PING"); err != nil {
			fmt.Printf("Failed to connect to master node due to %v", err)
			continue
		}

		return
	}
}

func writeMessage(msg string) error {
	masterAddr := fmt.Sprintf("%s:%s", nodeConfig.MasterConfig.Host, nodeConfig.MasterConfig.Port)

	fmt.Printf("masterAddr: %s\n", masterAddr)

	tcp, err := net.Dial("tcp", masterAddr)
	if err != nil {
		return fmt.Errorf("WriteMessage: %w", err)
	}

	defer tcp.Close()
	// tcp.SetDeadline(time.Now().Add(5 * time.Second))

	_, err = tcp.Write(buildCommand(msg))
	if err != nil {
		return fmt.Errorf("WriteMessage: %w", err)
	}

	return nil
}

func buildCommand(cmd string) []byte {
	rCmd := &StringValue{value: cmd}
	rValue := &ArrayValue{}
	rValue.values = append(rValue.values, rCmd)

	return rValue.Bytes()
}
