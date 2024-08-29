package main

import (
	"time"

	"github.com/luchaoshi45/ims/run"
)

func main() {
	go run.RunServer()
	time.Sleep(10 * time.Millisecond)
	run.RunClient()
}
