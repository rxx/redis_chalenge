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

type replicaFlags []string

func (f *replicaFlags) String() string {
	if len(*f) == 2 {
		return (*f)[0] + ":" + (*f)[1]
	}

	return ""
}

func (f *replicaFlags) Set(value string) error {
	fmt.Printf("Set %v\n", value)
	*f = append(*f, value)
	return nil
}

var masterConfig struct {
	host string
	port string
}

func main() {
	port := flag.String("port", "6379", "Redis port [Default: 6379]")
	masterHost := flag.String("replicaof", "", "master_host")
	masterPort := flag.String("", "", "port")

	flag.Parse()

	fmt.Printf("OS args %v\n", os.Args)
	fmt.Printf("Flag args %v\n", flag.Args())

	// OS args [/tmp/tmp.cFFFnP --port 6380 --replicaof localhost 6379]
	// [your_program] Flag args [6379]

	masterConfig.host = *masterHost
	masterConfig.port = *masterPort

	nodeConfig = &NodeConfig{Port: *port}
	role := "master"
	if len(masterConfig.host) > 0 {
		role = "slave"
		fmt.Printf("Master %v:%v\n", masterConfig.host, masterConfig.port)
	}

	nodeConfig.Role = role

	fmt.Printf("NodeConfig %v\n", nodeConfig)

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
