package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

const (
  MasterRole = "master"
  SlaveRole = "slave"
)

type NodeConfig struct {
  Host string
	Port string
	Role string
  MasterConfig *NodeConfig
}

var nodeConfig *NodeConfig


func main() {
	port, masterHost, masterPort := parseArgs()

	nodeConfig = &NodeConfig{Port: *port}
	role := "master"
	if len(masterConfig.host) > 0 {
		role = "slave"
		fmt.Printf("Master %v:%v\n", masterConfig.host, masterConfig.port)
	}

	nodeConfig.Role = role

	fmt.Printf("NodeConfig %v\n", nodeConfig)

	startServer()
}

func parseArgs() (string, string, string) {
  port := flag.String("port", "6379", "Redis port [Default: 6379]")
	masterHost := flag.String("replicaof", "", "master_host")
	tailArgs := flag.Args()
	flag.Parse()

  var masterPort string
	if len(tailArgs) > 0 {
    masterPort = tailArgs[0]
	}
  
  return *port, *masterHost, masterPort
}

func startServer() {
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
