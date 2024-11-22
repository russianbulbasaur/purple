package master

import (
	"context"
	"fmt"
	"log"
	"net"
	"purple/internals/my_resp"
	"strconv"
)

type Master struct {
	conn          net.Conn
	port          int
	host          string
	resp          *my_resp.MyRespObject
	masterContext context.Context
	cancel        context.CancelFunc
	writePipe     chan []byte
	ReadPipe      chan []byte
}

func NewMaster(host string, port int, serverPort int, resp *my_resp.MyRespObject) *Master {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalln(err, "Cannot connect to master")
	}
	ctx, cancel := context.WithCancel(context.Background())
	masterNode := &Master{
		conn,
		port,
		host,
		resp,
		ctx,
		cancel,
		make(chan []byte),
		make(chan []byte),
	}
	go masterNode.forkReader()
	go masterNode.forkWriter()
	masterNode.handshake(serverPort)
	return masterNode
}
func (master *Master) Write(message []byte) {
	master.writePipe <- message
}

func (master *Master) handshake(serverPort int) {
	master.Write(master.resp.E.EncodeBulkStringArray([]string{"PING"}))
	log.Println("Master says : ", string(<-master.ReadPipe))
	master.Write(master.resp.E.EncodeBulkStringArray(
		[]string{"REPLCONF", "listening-port", strconv.Itoa(serverPort)}))
	log.Println("Master says : ", string(<-master.ReadPipe))
	master.Write(master.resp.E.EncodeBulkStringArray(
		[]string{"REPLCONF", "capa", "psync2"}))
	log.Println("Master says : ", string(<-master.ReadPipe))
	master.Write(master.resp.E.EncodeBulkStringArray(
		[]string{"PSYNC", "?", "-1"}))
	log.Println("Master says : ", string(<-master.ReadPipe))
}

func (master *Master) forkReader() {
	var buffer []byte = make([]byte, 1024)
	for {
		n, err := master.conn.Read(buffer)
		if err != nil {
			log.Println("Master reader error : ", err)
			master.cancel()
			return
		}
		message := buffer[0:n]
		master.ReadPipe <- message
	}
}

func (master *Master) forkWriter() {
	for {
		select {
		case <-master.masterContext.Done():
			log.Println("Master writer breaking")
			return
		case message := <-master.writePipe:
			n, err := master.conn.Write(message)
			if n != len(message) {
				log.Println("Message too long")
			}
			if err != nil {
				log.Println(err)
			}
		}
	}
}
