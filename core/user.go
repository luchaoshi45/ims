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

const readTimeout = time.Duration(100 * time.Second) // 编译时直接替换
const WhoOnline string = "WhoOnline"
const Rename string = "Rename"
const PivateChat string = "PivateChat"

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

		buf_str := string(buf[:n])
		if buf_str[0:len(WhoOnline)] == WhoOnline {
			this.ReplyWhoOnline(server)
		} else if buf_str[0:len(Rename)] == Rename && len(buf_str) > len(Rename)+1 {
			this.ReplyRename(buf_str, server)
		} else if buf_str[0:len(PivateChat)] == PivateChat && len(buf_str) > len(PivateChat)+1 {
			this.ReplyPivateChat(buf_str, server)
		} else {
			server.BroadCast(this, buf_str[:n-1])
		}
	}
}

// ReplyWhoOnline /*查询在线用户
func (this *User) ReplyWhoOnline(server *Server) {
	for _, user := range server.OnlineMap {
		msg := user.Name
		this.C <- msg
	}
}

// ReplyRename /*修改用户名
func (this *User) ReplyRename(buf_str string, server *Server) {
	buf_str_split := strings.Split(buf_str, "|")
	new_name := buf_str_split[1]
	new_name = strings.Replace(new_name, "\n", "", -1)
	_, ok := server.OnlineMap[new_name]
	if ok {
		this.C <- "User name is used"
	} else {
		server.mapLock.Lock()
		delete(server.OnlineMap, this.Name)
		this.Name = new_name
		server.OnlineMap[new_name] = this
		server.mapLock.Unlock()
	}
}

// ReplyPivateChat /*私聊
// PivateChat|lcs|hello
func (this *User) ReplyPivateChat(buf_str string, server *Server) {
	buf_str_split := strings.Split(buf_str, "|")
	if len(buf_str_split) != 3 {
		this.C <- "Pivate chat grammatical errors"
	}
	name := buf_str_split[1]
	msg := "[" + name + "]: " + buf_str_split[2]
	server.OnlineMap[name].C <- msg
}
