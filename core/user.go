package core

import (
	"fmt"
	"io"
	"net"
	"runtime"
)

type User struct {
	Name string
	Addr string
	C    chan string
	coon net.Conn
	quit chan struct{}
}

func NewUser(conn net.Conn) *User {
	addr := conn.RemoteAddr().String()
	u := &User{addr, addr, make(chan string), conn, make(chan struct{})}
	go u.LMWriteUser()
	return u
}

// Online /*用户下线
func (this *User) Online(server *Server) {
	server.mapLock.Lock()
	server.OnlineMap[this.Name] = this
	server.mapLock.Unlock()
	server.BroadCast(this, "已上线")
	go this.LMRead(server)
}

// Offline /*用户上线
func (this *User) Offline(server *Server) {
	server.mapLock.Lock()
	delete(server.OnlineMap, this.Name)
	server.mapLock.Unlock()
	server.BroadCast(this, "已下线")
}

// LMWriteUser /*向用户发送消息
func (this *User) LMWriteUser() {
	for {
		select {
		case massage := <-this.C:
			this.coon.Write([]byte(massage + "\n"))
		case <-this.quit: // 接收退出信号
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
		n, err := this.coon.Read(buf)
		if n == 0 {
			this.Offline(server)
			close(this.quit)
			runtime.Goexit()
		}
		if err != nil && err != io.EOF {
			fmt.Println("coon Read err ", err)
		}
		server.BroadCast(this, string(buf[:n-1]))
	}
}
