package core

import (
	"context"
	"io"
	"log"
	"net"
	"runtime"
	"strings"
	"time"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	ctx    context.Context
	cancel context.CancelFunc
}

const readTimeout = time.Duration(10 * time.Second) // 编译时直接替换
const WhoOnline string = "WhoOnline"
const Rename string = "Rename"

func NewUser(conn net.Conn) *User {
	addr := conn.RemoteAddr().String()
	ctx, cancel := context.WithCancel(context.Background())
	u := &User{
		Name:   addr,
		Addr:   addr,
		C:      make(chan string),
		conn:   conn,
		ctx:    ctx,
		cancel: cancel,
	}
	return u
}

// Online /*用户下线
func (this *User) Online(server *Server) {
	server.mapLock.Lock()
	server.OnlineMap[this.Name] = this
	server.mapLock.Unlock()
	server.BroadCast(this, "Online")
	go this.LMWriteUser()
	go this.LMRead(server)
}

// Offline /*用户上线
func (this *User) Offline(server *Server) {
	server.mapLock.Lock()
	delete(server.OnlineMap, this.Name)
	server.mapLock.Unlock()
	server.BroadCast(this, "Offline")
}

// LMWriteUser /*向用户发送消息
func (this *User) LMWriteUser() {
	for {
		select {
		case massage := <-this.C:
			this.conn.Write([]byte(massage + "\n"))
		case <-this.ctx.Done(): // 接收退出信号
			runtime.Goexit()
		}
	}
}

// LMRead /*读取用户消息
/*
	server: 服务器
	关闭用户连接
	广播用户消息
*/
func (this *User) LMRead(server *Server) {
	buf := make([]byte, 4096)
	for {
		err := this.conn.SetReadDeadline(time.Now().Add(readTimeout)) // timeout
		if err != nil {
			log.Println("setReadDeadline failed:", err)
		}
		n, err := this.conn.Read(buf)

		if err != nil {
			//if netErr, ok := err.(net.Error); ok && netErr.Timeout()
			//Timeout 强制下线
			if err.(net.Error).Timeout() {
				this.C <- "You hava been forced offline (Timeout)"
				this.Offline(server)
				this.cancel()
				runtime.Goexit()
			}
			if err != io.EOF {
				log.Println("coon Read err ", err)
			}
		}

		//主动下线
		if n == 0 {
			this.Offline(server)
			this.cancel()
			runtime.Goexit()
		}

		buf_str := string(buf)
		if buf_str[0:len(WhoOnline)] == WhoOnline {
			this.ReplyWhoOnline(server)
		} else if buf_str[0:len(Rename)] == Rename && len(buf_str) > len(Rename)+1 {
			name := strings.Replace(buf_str[len(Rename)+1:], "\n", "", -1)
			this.ReplyRename(name, server)
		} else {
			server.BroadCast(this, buf_str[:n-1])
		}
	}
}

func (this *User) ReplyWhoOnline(server *Server) {
	for _, user := range server.OnlineMap {
		msg := user.Name
		this.C <- msg
	}
}

func (this *User) ReplyRename(name string, server *Server) {
	_, ok := server.OnlineMap[name]
	if ok {
		this.C <- "User name is used"
	} else {
		server.mapLock.Lock()
		delete(server.OnlineMap, this.Name)
		this.Name = name
		server.OnlineMap[name] = this
		server.mapLock.Unlock()
	}
}
