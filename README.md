# Distributed Chat Room using Go RPC

A concurrent chat application built in Go using client-server architecture with Go's `net/rpc` package and Polling communication mechanism.

## Features
- Central Server: Manages connected users, distributes messages, and handles join/leave notifications.
- Unique Usernames: Validates and prevents duplicate login names.
- Polling Mechanism: Clients periodically fetch pending messages from the server using RPC.
- Concurrency Safety: Uses `sync.Mutex` and `goroutines` to safely handle multiple concurrent clients.
- Graceful Shutdown: Handles disconnections cleanly via `/quit` or `Ctrl+C`.

## Project Structure
```text
chat-app/
├── go.mod
├── common/
│   └── types.go
├── server/
│   ├── chat.go
│   └── main.go
└── client/
    └── main.go