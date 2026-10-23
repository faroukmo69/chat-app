package main

import (
	"chat-app/common"
	"fmt"
	"log"
	"sync"
)

type ChatServer struct {
	mu        sync.Mutex
	users     map[string]bool
	mailboxes map[string][]string
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		users:     make(map[string]bool),
		mailboxes: make(map[string][]string),
	}
}

func (s *ChatServer) Join(args *common.JoinArgs, reply *common.JoinReply) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if args.Username == "" {
		reply.Success = false
		reply.Message = "Username cannot be empty"
		return nil
	}

	if s.users[args.Username] {
		reply.Success = false
		reply.Message = "Username already taken"
		return nil
	}

	s.users[args.Username] = true
	s.mailboxes[args.Username] = []string{}

	notif := fmt.Sprintf("📢 [System]: User '%s' has joined the chat.", args.Username)
	s.broadcast(args.Username, notif)

	reply.Success = true
	reply.Message = "Joined successfully"
	log.Printf("[+] User joined: %s\n", args.Username)
	return nil
}

func (s *ChatServer) SendMessage(args *common.MessageArgs, reply *common.StandardReply) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.users[args.Username] {
		reply.Success = false
		reply.Message = "You are not joined"
		return nil
	}

	msg := fmt.Sprintf("[%s]: %s", args.Username, args.Text)
	s.broadcast(args.Username, msg)

	reply.Success = true
	return nil
}

func (s *ChatServer) GetUpdates(args *common.PollArgs, reply *common.PollReply) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.users[args.Username] {
		reply.Messages = nil
		return nil
	}

	reply.Messages = s.mailboxes[args.Username]
	s.mailboxes[args.Username] = []string{}
	return nil
}

func (s *ChatServer) Leave(args *common.LeaveArgs, reply *common.StandardReply) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.users[args.Username] {
		delete(s.users, args.Username)
		delete(s.mailboxes, args.Username)

		notif := fmt.Sprintf("📢 [System]: User '%s' has left the chat.", args.Username)
		s.broadcast(args.Username, notif)
		log.Printf("[-] User left: %s\n", args.Username)
	}

	reply.Success = true
	return nil
}

func (s *ChatServer) broadcast(sender string, message string) {
	for user := range s.users {
		if user != sender {
			s.mailboxes[user] = append(s.mailboxes[user], message)
		}
	}
}
