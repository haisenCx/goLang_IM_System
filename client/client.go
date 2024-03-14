package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int // current operation
}

func NewClient(serverIp string, serverPort int) *Client {
	// Connect to the server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("Error connecting:", err)
		return nil
	}
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		conn:       conn,
		flag:       999,
	}

	return client
}

var ServerIp string
var ServerPort int

// use flag to init client
func init() {
	flag.StringVar(&ServerIp, "ip", "127.0.0.1", "Set server IP address(default is 127.0.0.1)")
	flag.IntVar(&ServerPort, "port", 8888, "Set server port(default is 8888)")
}

// add a menu , to select the operation
// 1. public chat
// 2. private chat
// 3. rename
// 4. quit
func (client *Client) menu() bool {
	var flag int
	fmt.Println("1. Public chat")
	fmt.Println("2. Private chat")
	fmt.Println("3. Rename")
	fmt.Println("0. Quit")
	fmt.Println("Please select(1-4):")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("select error:")
		return false
	}
}

// create a goroutine to deal with the server response , and print it to the console
func (client *Client) DealResponse() {
	//get the server response
	io := client.conn
	buf := make([]byte, 1024)
	for {
		n, err := io.Read(buf)
		if n == 0 {
			fmt.Println("Server disconnected")
			return
		}
		if err != nil {
			fmt.Println("Server read error:", err)
			return
		}
		fmt.Println(string(buf[:n]))
	}
}

// Run the client
func (client *Client) Run() {
	for client.flag != 0 {
		for !client.menu() {

		}
		switch client.flag {
		case 1:
			fmt.Println("Public chat")
			client.PublicChat()
			break
		case 2:
			client.PrivateChat()
			break
		case 3:
			// Rename
			//call the server to rename
			//-user rename|newName
			client.updateName()
			break
		case 0:
			fmt.Println("Quit")
			break
		default:
			fmt.Println("Select error")
		}

	}

}

// private chat
func (client *Client) PrivateChat() {
	//show the online users
	client.ShowOnlineUsers()
	//select the user to chat
	fmt.Println("Please enter the username:")
	var remoteName string
	fmt.Scanln(&remoteName)
	fmt.Println("Please enter the message, type 'exit' to quit:")
	var chatMsg string
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			client.sendMsg("to|" + remoteName + "|" + chatMsg + "\n")
		}
		chatMsg = ""
		fmt.Scanln(&chatMsg)
	}
}

// show the online users
func (client *Client) ShowOnlineUsers() {
	client.sendMsg("-user ls\n")
}

// update the user name
func (client *Client) updateName() bool {
	fmt.Println("Please enter your new name:")
	fmt.Scanln(&client.Name)
	msg := "-user rename|" + client.Name + "\n"
	//send the msg to the server
	return client.sendMsg(msg)
}

// public chat
func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println("Please enter the message, type 'exit' to quit:")
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			client.sendMsg(chatMsg)
		}
		chatMsg = ""
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) sendMsg(msg string) bool {
	_, err := client.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("conn write err:", err)
		return true
	}
	return false
}

func main() {
	flag.Parse()

	// Create a new client
	client := NewClient(ServerIp, ServerPort)

	go client.DealResponse()

	if client == nil {
		fmt.Println(">>>>>>>>>>>>>> Failed to connect to server")
		return
	}

	fmt.Println(">>>>>>>>>>>>>> Connected to the server")

	client.Run()
}
