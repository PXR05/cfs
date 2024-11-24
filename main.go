package main

import (
	"cfs/server"
)

func main() {
	server := server.Server{}
	server.Init()
	defer server.Close()

	server.Run()
}
