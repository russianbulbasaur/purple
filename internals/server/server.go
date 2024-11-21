package server

import (
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"purple/internals/client"
	"purple/internals/rdb"
	"strconv"
	"time"
)

type Server struct {
	port          int
	address       string
	serverContext context.Context
	void          map[string]dataNode
	rdbFile       *rdb.RDBFile
}

func NewServer(port int, address string, config map[string]string) Server {
	rdbFile := rdb.NewRDBFile(config["dbfilename"], config["dir"])
	return Server{
		port,
		address,
		context.Background(),
		make(map[string]dataNode),
		rdbFile,
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
		newClient := client.NewClient(connection, server.set, server.get, server.rdbFile)
		go newClient.Handle()
	}
}

func (server *Server) set(key string, value interface{}, expiry int64) {
	if expiry != math.MaxInt64 {
		expiry = time.Now().Unix()*1000 + expiry
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
	if data.expiry > time.Now().Unix()*1000 {
		return data.value
	}
	return nil
}
