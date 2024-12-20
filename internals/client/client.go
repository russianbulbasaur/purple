package client

import (
	"context"
	"log"
	"net"
	resp "purple/internals/my_resp"
	types "purple/internals/my_resp/purple_data_types"
	arrayTypes "purple/internals/my_resp/purple_data_types/array"
	"purple/internals/rdb"
	"runtime"
)

type Client struct {
	conn          net.Conn
	resp          *resp.MyRespObject
	writeChannel  chan []byte
	clientContext context.Context
	cancel        context.CancelFunc
	set           func(string, interface{}, int64)
	get           func(string) interface{}
	getAll        func() map[string]models.DataNode
	rdbFile       *rdb.RDBFile
	serverConfig  map[string]interface{}
}

func NewClient(connection net.Conn, set func(string, interface{}, int64), get func(string) interface{},
	rdbFile *rdb.RDBFile, getAll func() map[string]models.DataNode,
	config map[string]interface{}) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		connection,
		resp.Init(),
		make(chan []byte),
		ctx,
		cancel,
		set,
		get,
		getAll,
		rdbFile,
		config,
	}
}

func (client *Client) Handle() {
	log.Printf("Goroutines %d", runtime.NumGoroutine())
	go client.forkReader()
	go client.forkWriter()
}

func (client *Client) forkWriter() {
	for {
		select {
		case <-client.clientContext.Done():
			log.Println("Killing writer")
			log.Printf("Goroutines %d", runtime.NumGoroutine())
			return
		case input := <-client.writeChannel:
			client.conn.Write(input)
		}
	}
}

func (client *Client) forkReader() {
	var buffer []byte = make([]byte, 4096)
	for {
		count, err := client.conn.Read(buffer)
		if err != nil {
			log.Println("Killing reader", err)
			client.kill()
			break
		}
		message := buffer[0:count]
		purpleType, err, _ := client.resp.D.Decode(message)
		if err != nil {
			log.Println(err)
		}
		switch purpleType.(type) {
		case arrayTypes.PurpleArray:
			client.evaluateArray(purpleType.(arrayTypes.PurpleArray))
		case types.PurpleString:
			client.evaluateString(purpleType.(types.PurpleString))
		}
	}
}

func (client *Client) kill() {
	client.cancel()
}
