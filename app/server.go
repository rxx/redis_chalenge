package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

type NodeConfig struct {
	Port string
	Role string
}

var nodeConfig *NodeConfig

func main() {
	port := flag.String("port", "6379", "Redis port [Default: 6379]")
	replicaof := flag.String("replicaof", "", "")
	flag.Parse()

	nodeConfig = &NodeConfig{Port: *port}
	role := "master"
	if len(*replicaof) > 0 {
		role = "slave"
	}

	nodeConfig.Role = role

	fmt.Printf("NodeConfig %v", nodeConfig)

	tcp, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", nodeConfig.Port))
	if err != nil {
		fmt.Printf("Failed to bind to port %s\n", nodeConfig.Port)
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
