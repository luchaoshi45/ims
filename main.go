package main

import "ims/core"

func main() {
	Ip := "127.0.0.1"
	Port := 8888
	s := core.NewServer(Ip, Port)
	s.Start()
}
