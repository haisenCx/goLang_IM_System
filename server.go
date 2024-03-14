package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//onlineUsersMap
	OnlineMap map[string]*User
	mapLcok   sync.RWMutex

	//msg broadcast chaneel
	Message chan string
}

// create server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// ListenMessager goruntine
func (server *Server) ListenMessager() {
	for {
		msg := <-server.Message

		//msg broadcast to all online user
		server.mapLcok.Lock()
		for _, cli := range server.OnlineMap {
			cli.C <- msg
		}
		server.mapLcok.Unlock()
	}
}

// broadcast msg
func (server *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	server.Message <- sendMsg
}

// do handler
// Handler handles the incoming connection and manages the user's online status.
// It reads messages from the user and performs actions based on the received messages.
// If the user is inactive for more than 10 seconds, they will be forcefully disconnected.
func (server *Server) Handler(conn net.Conn) {
	//fmt.Println("Conn created");
	user := NewUser(conn, server)

	//user Online
	user.Online()
	var isLive = make(chan bool)

	//listen msgs from users--
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("conn read err:", err)
				return
			}
			//get msg from users,(except "/n")
			msg := string(buf[:n-1])

			user.DoMessage(msg)

			isLive <- true
		}
	}()
	for {
		select {
		case <-isLive:

		case <-time.After(3000 * time.Second):
			user.SendMsg("you not live")
			user.Offline()
			conn.Close()
			return

		}
	}

}

func (server *Server) Start() {
	fmt.Println("Start")

	//socket Listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("net.listen err:", err)
		return
	}
	//close listen socket
	defer listener.Close()

	//start a ListenMessager goruntine
	go server.ListenMessager()
	fmt.Println("ListenMessager Started")

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		//do handler (a user onlin)
		go server.Handler(conn)
	}

}
