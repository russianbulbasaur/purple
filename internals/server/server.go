package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"purple/internals/client"
	"strconv"
)

type Server struct {
	port          int
	address       string
	serverContext context.Context
}

func NewServer(port int, address string) Server {
	return Server{
		port,
		address,
		context.Background(),
	}
}

func (server *Server) Listen() {
	listener, err := net.Listen("tcp4",
		fmt.Sprintf("%s:%s", server.address, strconv.Itoa(server.port)))
	if err != nil {
		panic(err)
	}
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Println("Error accept client", err)
			continue
		}
		newClient := client.NewClient(connection)
		go newClient.Handle()
	}
}
