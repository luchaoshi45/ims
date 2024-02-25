package core

import (
	"fmt"
	"log"
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

// Start /*服务器运行循环
func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		log.Println("net.Listen err: ", err)
		return
	}
	defer listener.Close()

	go this.LMWriteAllUser()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept err: ", err)
			continue
		}
		user := NewUser(conn)
		user.Online(this)
	}

}

// BroadCast /*广播函数，写入Message
func (this *Server) BroadCast(user *User, msg string) {
	msg = "[" + user.Addr + "]" + " " + user.Name + " " + msg
	this.Message <- msg
}

// LMWriteAllUser /*向所以用户发送消息
func (this *Server) LMWriteAllUser() {
	for {
		msg := <-this.Message
		this.mapLock.Lock()
		for _, user := range this.OnlineMap {
			user.C <- msg
		}
		this.mapLock.Unlock()
	}
}
