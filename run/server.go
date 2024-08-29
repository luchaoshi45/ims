package run

import "github.com/luchaoshi45/ims/core"

func RunServer() {
	Ip := "127.0.0.1"
	Port := 8888
	server := core.NewServer(Ip, Port)
	server.Start()
}
