package core

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	Message   chan string
}

func NewServer(Ip string, Port int) *Server {
	s := &Server{
		Ip:        Ip,
		Port:      Port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return s
}

func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}
	defer listener.Close()

	go this.ListenMessage()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err: ", err)
			continue
		}
		go this.Handler(conn)
	}

}

func (this *Server) Handler(conn net.Conn) {
	//fmt.Println("连接成功")
	user := NewUser(conn)
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()
	this.BroadCast(user, "已上线")
}

func (this *Server) BroadCast(user *User, msg string) {
	msg = "[" + user.Addr + "]" + " " + user.Name + " " + msg
	this.Message <- msg
}

func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		this.mapLock.Lock()
		for _, user := range this.OnlineMap {
			user.C <- msg
		}
		this.mapLock.Unlock()
	}
}
