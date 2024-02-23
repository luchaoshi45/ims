package core

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(Ip string, Port int) *Server {
	s := &Server{Ip, Port}
	return s
}

func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}
	defer listener.Close()

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
	fmt.Println("连接成功")
}
