package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

const (
	MasterRole = "master"
	SlaveRole  = "slave"
	LOCAL      = "localhost"
)

type NodeConfig struct {
	Host         string
	Port         string
	Role         string
	MasterConfig *NodeConfig
}

func (v *NodeConfig) String() string {
	var master string
	if v.MasterConfig != nil {
		master = fmt.Sprintf("%s:%s", v.MasterConfig.Host, v.MasterConfig.Port)
		return fmt.Sprintf("Role: %s at %s:%s from master %s", v.Role, v.Host, v.Port, master)
	}

	return fmt.Sprintf("Role: %s at %s:%s", v.Role, v.Host, v.Port)
}

var nodeConfig *NodeConfig

func main() {
	nodeConfig = NewNodeConfig()
	fmt.Printf("NodeConfig %s\n", nodeConfig.String())

	if nodeConfig.Role == MasterRole {
		fmt.Println("Start new master server")
		startMasterServer()
	} else {
		fmt.Println("Start new slave server")
		StartReplica()
	}
}

func NewNodeConfig() *NodeConfig {
	host, port, masterHost, masterPort := parseArgs()

	role := MasterRole
	var masterConfig *NodeConfig
	if len(masterHost) > 0 && len(masterPort) > 0 {
		role = SlaveRole
		masterConfig = &NodeConfig{
			Host: masterHost,
			Port: masterPort,
		}
	}

	return &NodeConfig{
		Host:         host,
		Port:         port,
		Role:         role,
		MasterConfig: masterConfig,
	}
}

func parseArgs() (string, string, string, string) {
	port := flag.String("port", "6379", "Redis port [Default: 6379]")
	masterHost := flag.String("replicaof", "", "master_host")
	flag.Parse()

	host := LOCAL
	tailArgs := flag.Args()

	var masterPort string
	if len(tailArgs) > 0 {
		masterPort = tailArgs[0]
	}

	return host, *port, *masterHost, masterPort
}

func startMasterServer() {
	addr := fmt.Sprintf("%s:%s", nodeConfig.Host, nodeConfig.Port)
	fmt.Printf("startMasterServer: on %s\n", addr)

	tcp, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Failed to bind to port %s\n", nodeConfig.Port)
		os.Exit(1)
	}

	defer tcp.Close()

	for {
		conn, err := tcp.Accept()
		// conn.SetDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			return
		}

		go HandleClient(conn)
	}
}
