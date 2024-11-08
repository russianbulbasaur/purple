package server

import (
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"purple/internals/client"
	"strconv"
	"time"
)

type Server struct {
	port          int
	address       string
	serverContext context.Context
	void          map[string]dataNode
}

func NewServer(port int, address string) Server {
	return Server{
		port,
		address,
		context.Background(),
		make(map[string]dataNode),
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
		newClient := client.NewClient(connection, server.set, server.get)
		go newClient.Handle()
	}
}

func (server *Server) set(key string, value interface{}, expiry int64) {
	if expiry != math.MaxInt64 {
		expiry = time.Now().Unix() + expiry
	}
	data := dataNode{
		value, expiry,
	}
	server.void[key] = data
}

func (server *Server) get(key string) interface{} {
	data, ok := server.void[key]
	if !ok {
		return nil
	}
	if data.expiry == math.MaxInt64 {
		return data.value
	}
	if data.expiry >= time.Now().Unix() {
		return data.value
	}
	return nil
}
