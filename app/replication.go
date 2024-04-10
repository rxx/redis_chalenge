package main

import (
	"fmt"
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
