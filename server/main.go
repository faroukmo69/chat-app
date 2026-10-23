package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

func main() {
	server := NewChatServer()
	_ = rpc.Register(server)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
	defer listener.Close()

	fmt.Println("🚀 Chat Server is running on port 1234...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(conn)
	}
}
