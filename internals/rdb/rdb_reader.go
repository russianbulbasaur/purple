package rdb

import (
	"bufio"
	"io"
	"log"
	"os"
)

type RDBReader struct {
	osFile    *os.File
	bufReader *bufio.Reader
}

func NewRDBReader(file *RDBFile) *RDBReader {
	bufferedReader := bufio.NewReader(file.file)
	return &RDBReader{
		file.file,
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

func (reader *RDBReader) ReadEOFSection() {
	skipped, err := reader.osFile.Seek(0, io.SeekEnd)
	if err != nil {
		return
	}
	log.Printf("Skipped %d bytes", skipped)
}

func (reader *RDBReader) ReadAuxiliaryFields() {

}
