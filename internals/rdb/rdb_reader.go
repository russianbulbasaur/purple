package rdb

import (
	"bufio"
	"log"
)

type RDBReader struct {
	bufReader *bufio.Reader
}

func NewRDBReader(file *RDBFile) *RDBReader {
	bufferedReader := bufio.NewReader(file.file)
	return &RDBReader{
		bufferedReader,
	}
}

func (reader *RDBReader) ReadHeader() {
	buffer := make([]byte, 9)
	_, err := reader.bufReader.Read(buffer)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Read : %s", string(buffer))
}
