package core

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	coon net.Conn
}

func NewUser(conn net.Conn) *User {
	addr := conn.RemoteAddr().String()
	u := &User{addr, addr, make(chan string), conn}
	go u.ListenMessage()
	return u
}

func (this *User) ListenMessage() {
	for {
		massage := <-this.C
		this.coon.Write([]byte(massage + "\n"))
	}
}
