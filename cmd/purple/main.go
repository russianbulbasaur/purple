package main

import "purple/internals/server"

func main() {
	purpleServer := server.NewServer(8000, "localhost")
	purpleServer.Listen()
}
