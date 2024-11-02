package main

import "purple/internals/server"

func main() {
	purpleServer := server.NewServer(6379, "localhost")
	purpleServer.Listen()
}
