package client

import (
	"log"
	"net"
	resp "purple/internals/my_resp"
)

type Client struct {
	conn net.Conn
	resp *resp.MyRespObject
}

func NewClient(connection net.Conn) *Client {
	return &Client{
		connection,
		resp.Init(),
	}
}

func (client *Client) Handle() {
	var buffer []byte = make([]byte, 4096)
	for {
		count, err := client.conn.Read(buffer)
		if err != nil {
			log.Println(err)
			break
		}
		message := buffer[0 : count-1]
		idk, purpleType := client.resp.D.Decode(message)
		log.Printf("%#v %#v", idk, purpleType)
	}
}
