# GoLang IM System

This project is a simple Instant Messaging (IM) system implemented in Go. It consists of a server and a client, which can communicate with each other over the network.

## Features

- Public chat: All users can send messages to a public chat room, and all messages in the public chat room can be seen by all users.
- Private chat: Users can send messages to a specific user privately.
- Online user list: Users can see a list of all online users.
- User rename: Users can change their username.

## Usage

### Server

To start the server, navigate to the server directory and build the server with the following command:

```bash
go build -o server server.go
```

Then, run the server with:

```bash
./server
```

### Client

To start a client, navigate to the client directory and build the client with the following command:

```bash
go build -o client client.go
```

Then, run the client with:

```bash
./client
```

Once the client is running, follow the prompts in the terminal to chat publicly, chat privately, see the online user list, or change your username.

## Note

This is a simple IM system for learning purposes. It does not implement many features that a full-fledged IM system would have, such as authentication, encryption, file transfer, voice chat, etc. Use it as a starting point for learning network programming with Go.
