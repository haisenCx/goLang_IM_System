package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// create new User
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}

	go user.ListenMessage()

	return user

}

func (user *User) Online() {
	//add user into the OnlineMap
	user.server.mapLcok.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLcok.Unlock()

	//broadcast msg to online user
	user.server.BroadCast(user, "online")

}

func (user *User) Offline() {
	//delete user from the OnlineMap
	user.server.mapLcok.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLcok.Unlock()

	//broadcast msg to online user
	user.server.BroadCast(user, "offline")
}

func (user *User) SendMsg(msg string) {
	user.conn.Write([]byte(msg))
}

func (thisUser *User) DoMessage(msg string) {
	if msg == "-user ls" {
		thisUser.server.mapLcok.Lock()
		for _, tmpUser := range thisUser.server.OnlineMap {
			onlineMsg := "[" + tmpUser.Addr + "]" + tmpUser.Name + ":online"
			thisUser.SendMsg(onlineMsg)

		}
		thisUser.server.mapLcok.Unlock()
	} else if strings.HasPrefix(msg, "-user rename|") {
		newName := strings.TrimPrefix(msg, "-user rename|")
		newName = strings.TrimSpace(newName)

		// Check if the new username already exists
		thisUser.server.mapLcok.Lock()
		_, exists := thisUser.server.OnlineMap[newName]
		thisUser.server.mapLcok.Unlock()

		if exists {
			thisUser.SendMsg("Username already exists")
		} else {
			// Update the username
			thisUser.server.mapLcok.Lock()
			delete(thisUser.server.OnlineMap, thisUser.Name)
			thisUser.Name = newName
			thisUser.server.OnlineMap[thisUser.Name] = thisUser
			thisUser.server.mapLcok.Unlock()

			thisUser.SendMsg("Username updated successfully")
		}

	} else if strings.HasPrefix(msg, "to|") {
		//msg format: to|username|msg
		//1.get the	username

		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			thisUser.SendMsg("msg format error, should be to|username|msg")
			return
		}
		//2.get userObj from OnlineMap
		remoteUser, ok := thisUser.server.OnlineMap[remoteName]
		if !ok {
			thisUser.SendMsg("username not exists")
			return
		}

		//3.get the msg
		content := strings.Split(msg, "|")[2]
		if content == "" {
			thisUser.SendMsg("msg content is empty")
			return
		}
		remoteUser.SendMsg(thisUser.Name + " to you:" + content)

	} else {
		thisUser.server.BroadCast(thisUser, msg)
	}
}

// Listen Current User C channel, send msg directly when a msg accepted
func (user *User) ListenMessage() {
	for {
		msg := <-user.C

		user.conn.Write([]byte(msg + "\n"))
	}
}
