package rdb

import (
	"bufio"
	"log"
	"os"
	"path"
)

type RDBReader struct {
	rdbFile    *os.File
	readBuffer *bufio.Reader
}

func NewRDBReader(filepath string, filename string) *RDBReader {
	fullPath := path.Join(filepath, filename)
	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		log.Fatalf("Cannot open file %s", fullPath)
		return nil
	}
	bufferedReader := bufio.NewReader(file)
	return &RDBReader{
		file,
		bufferedReader,
	}
}

func (reader *RDBReader) ReadHeader() {

}
