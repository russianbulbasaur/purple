package server

import (
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"purple/internals/client"
	"purple/internals/master"
	"purple/internals/my_resp"
	"purple/internals/rdb"
	"purple/models"
	"strconv"
	"time"
)

type Server struct {
	port          int
	address       string
	serverContext context.Context
	resp          *my_resp.MyRespObject
	void          map[string]models.DataNode
	rdbFile       *rdb.RDBFile
	rdbReader     *rdb.RDBReader
	config        map[string]interface{}
	masterNode    *master.Master
}

func NewServer(port int, address string, config map[string]interface{}) Server {
	resp := my_resp.Init()
	var masterNode *master.Master
	if config["role"] == "slave" {
		masterNode = connectToMaster(config["master"].(string),
			config["master_port"].(string), resp, port)
	}
	config["master_replid"] = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
	config["master_repl_offset"] = 0
	rdbFile := rdb.NewRDBFile(config["dbfilename"].(string), config["dir"].(string))
	rdbReader := rdb.NewRDBReader(rdbFile)
	void := make(map[string]models.DataNode)
	if rdbReader.IsValidRDB() {
		rdbReader.ReadDatabase(void)
	}
	return Server{
		port,
		address,
		context.Background(),
		resp,
		void,
		rdbFile,
		rdbReader,
		config,
		masterNode,
	}
}

func connectToMaster(masterAddress string, masterPort string,
	resp *my_resp.MyRespObject, serverPort int) *master.Master {
	port, err := strconv.Atoi(masterPort)
	if err != nil {
		log.Fatalln(err)
	}
	return master.NewMaster(masterAddress, port, serverPort, resp)
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
		newClient := client.NewClient(connection, server.set, server.get,
			server.rdbFile, server.getAll, server.config)
		go newClient.Handle()
	}
}

func (server *Server) getAll() map[string]models.DataNode {
	return server.void
}

func (server *Server) set(key string, value interface{}, expiry int64) {
	if expiry != math.MaxInt64 {
		log.Printf("setting expiry to %d ms", expiry)
		expiry = time.Now().UnixMilli() + expiry
	}
	data := models.DataNode{
		Value: value, Expiry: expiry,
	}
	server.void[key] = data
}

func (server *Server) get(key string) interface{} {
	data, ok := server.void[key]
	if !ok {
		return nil
	}
	if data.Expiry == math.MaxInt64 {
		return data.Value
	}
	log.Printf("expiry : %d Current : %d", data.Expiry, time.Now().UnixMilli())
	if data.Expiry >= time.Now().UnixMilli() {
		return data.Value
	}
	return nil
}
