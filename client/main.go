package main

import (
	"bufio"
	"chat-app/common"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Could not connect to server. Make sure server is running. Error: ", err)
	}
	defer client.Close()

	reader := bufio.NewReader(os.Stdin)
	var username string

	for {
		fmt.Print("Enter your unique username: ")
		input, _ := reader.ReadString('\n')
		username = strings.TrimSpace(input)
		if username == "" {
			fmt.Println("Username cannot be empty!")
			continue
		}

		args := &common.JoinArgs{Username: username}
		var reply common.JoinReply
		err = client.Call("ChatServer.Join", args, &reply)
		if err != nil {
			log.Fatal("RPC error:", err)
		}

		if reply.Success {
			fmt.Println("✅ Connected successfully! Type messages below (Type /quit to exit):")
			break
		} else {
			fmt.Printf("❌ Rejected: %s. Try another name.\n", reply.Message)
		}
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nDisconnecting...")
		var reply common.StandardReply
		_ = client.Call("ChatServer.Leave", &common.LeaveArgs{Username: username}, &reply)
		os.Exit(0)
	}()

	go func() {
		for {
			time.Sleep(1 * time.Second)
			args := &common.PollArgs{Username: username}
			var reply common.PollReply
			err := client.Call("ChatServer.GetUpdates", args, &reply)
			if err != nil {
				continue
			}
			for _, msg := range reply.Messages {
				fmt.Printf("\r%s\n> ", msg)
			}
		}
	}()

	fmt.Print("> ")
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "" {
			fmt.Print("> ")
			continue
		}

		if text == "/quit" {
			var reply common.StandardReply
			client.Call("ChatServer.Leave", &common.LeaveArgs{Username: username}, &reply)
			fmt.Println("Disconnected.")
			break
		}

		args := &common.MessageArgs{
			Username: username,
			Text:     text,
		}
		var reply common.StandardReply
		err = client.Call("ChatServer.SendMessage", args, &reply)
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
		fmt.Print("> ")
	}
}
