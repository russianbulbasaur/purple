package rdb

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"os"
	"strconv"
)

const auxiliaryOPCode byte = 0xfa
const dbOPCode byte = 0xfe

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
	var auxiliaryMap map[string]string = make(map[string]string)
	for {
		readByte, err := reader.bufReader.ReadByte()
		if err != nil {
			log.Fatalln(err)
			return
		}
		if readByte == auxiliaryOPCode {
			key := reader.length_decode_string()
			log.Printf("Key : %s", key)
			value := reader.length_decode_string()
			auxiliaryMap[key] = value
			log.Printf("%s : %s", key, value)
		} else if readByte == dbOPCode {
			err := reader.bufReader.UnreadByte()
			if err != nil {
				log.Fatalln(err)
			}
			break
		}
	}
	log.Println("Read auxiliary fields")
}

func (reader *RDBReader) length_decode_string() string {
	firstByte, err := reader.bufReader.ReadByte()
	if err != nil {
		log.Fatalln(err)
		return ""
	}
	//checking the first two msbs
	firstByteShifted := firstByte >> 6
	switch int8(firstByteShifted) {
	case 0:
		//The next 6 bits represent the length
		stringLength := (firstByte << 2) >> 2
		var stringBytes []byte = make([]byte, stringLength)
		_, err := reader.bufReader.Read(stringBytes)
		if err != nil {
			log.Fatalln(err)
			return ""
		}
		value := string(stringBytes)
		return value
	case 1:
		//Read one additional byte. The combined 14 bits represent the length
	case 2:
		//Discard the remaining 6 bits. The next 4 bytes from the stream represent the length
	case 3:
		//The next object is encoded in a special format. The remaining 6 bits indicate the format
		remainingBits := (firstByte << 2) >> 2
		result := reader.readInteger(remainingBits)
		return strconv.Itoa(result)
	}
	return ""
}

func (reader *RDBReader) readInteger(remainingBits byte) int {
	flag := int8(remainingBits)
	var buffer []byte
	if flag == 0 {
		//8 bit int
		buffer = make([]byte, 1)
		_, err := reader.bufReader.Read(buffer)
		log.Println(buffer)
		if err != nil {
			log.Fatalln(err)
		}
		var result int8
		err = binary.Read(bytes.NewBuffer(buffer), binary.LittleEndian, &result)
		if err != nil {
			log.Fatalln(err)
		}
		return int(result)
	} else if flag == 1 {
		//8 bit int
		buffer = make([]byte, 2)
		_, err := reader.bufReader.Read(buffer)
		log.Println(buffer)
		if err != nil {
			log.Fatalln(err)
		}
		var result int16
		err = binary.Read(bytes.NewBuffer(buffer), binary.LittleEndian, &result)
		if err != nil {
			log.Fatalln(err)
		}
		return int(result)
	} else if flag == 2 {
		//8 bit int
		buffer = make([]byte, 4)
		_, err := reader.bufReader.Read(buffer)
		log.Println(buffer)
		if err != nil {
			log.Fatalln(err)
		}
		var result int32
		err = binary.Read(bytes.NewBuffer(buffer), binary.LittleEndian, &result)
		if err != nil {
			log.Fatalln(err)
		}
		return int(result)
	}
	return -1
}
