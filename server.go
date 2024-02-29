package main

import "ims/core"

func main() {
	Ip := "127.0.0.1"
	Port := 8888
	server := core.NewServer(Ip, Port)
	server.Start()
}
