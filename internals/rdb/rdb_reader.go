package rdb

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"os"
)

const auxiliaryOPCode byte = 0xfa
const databaseOPCode byte = 0xfe
const resizeDBOPCode byte = 0xfb
const secondsExpiryOPCode byte = 0xfd
const millisecondsExpiryOPCode byte = 0xfc

type BitCase int8

type Length int

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

func (reader *RDBReader) readHeader() string {
	buffer := make([]byte, 5)
	_, err := reader.bufReader.Read(buffer)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(buffer)
}

func (reader *RDBReader) IsValidRDB() bool {
	return reader.readHeader() == "REDIS"
}

func (reader *RDBReader) ReadEOFSection() {
	skipped, err := reader.osFile.Seek(0, io.SeekEnd)
	if err != nil {
		return
	}
	log.Printf("Skipped %d bytes", skipped)
}

func (reader *RDBReader) ReadAuxiliaryFields() {
	//var auxiliaryMap map[string]string = make(map[string]string)
	for {
		readByte, err := reader.bufReader.ReadByte()
		if err != nil {
			log.Fatalln(err)
			return
		}
		if readByte == auxiliaryOPCode {
			bitCase, length := reader.decodeLength()
			key := reader.readLengthPrefixedString(int(length))
			bitCase, length = reader.decodeLength()
			if bitCase == 3 {
				value := reader.readInteger(int(length))
				log.Printf("%s : %d", key, value)
			} else {
				value := reader.readLengthPrefixedString(int(length))
				log.Printf("%s : %s", key, value)
			}
		} else if readByte == databaseOPCode {
			err := reader.bufReader.UnreadByte()
			if err != nil {
				log.Fatalln(err)
			}
			break
		}
	}
	log.Println("Read auxiliary fields")
}

func (reader *RDBReader) ReadDatabase(result map[string]models.DataNode) {
	for {
		readByte, err := reader.bufReader.ReadByte()
		if err != nil {
			log.Fatalln("DB Read : ", err)
		}
		if readByte == databaseOPCode {
			hashTableSize := reader.getHashTableSize(false)
			expireHashTableSize := reader.getHashTableSize(true)
			log.Printf("Hash table size : %d", hashTableSize)
			log.Printf("Expire Hash table size : %d", expireHashTableSize)
			reader.readValues(hashTableSize, result)
			break
		}
	}
}

func (reader *RDBReader) getHashTableSize(expire bool) int {
	if expire {
		_, length := reader.decodeLength()
		return int(length)
	}
	for {
		readByte, err := reader.bufReader.ReadByte()
		if err != nil {
			log.Fatalln("DB Read : ", err)
		}
		if readByte == resizeDBOPCode {
			_, length := reader.decodeLength()
			return int(length)
		}
	}
	return 0
}

func (reader *RDBReader) readValues(hashTableSize int, result map[string]models.DataNode) {
	for i := 0; i < hashTableSize; i++ {
		readBytes, err := reader.bufReader.Peek(1)
		readByte := readBytes[0]
		if err != nil {
			log.Fatalln("value read : ", err)
		}
		var key string
		var value interface{}
		var expiry int64 = math.MaxInt64
		if readByte == secondsExpiryOPCode {
			key, value, expiry = reader.readSecondsExpiryValue()
		} else if readByte == millisecondsExpiryOPCode {
			key, value, expiry = reader.readMillisecondsExpiryValue()
		} else {
			key, value = reader.readNoExpiryValue()
		}
		result[key] = models.DataNode{
			Value:  value,
			Expiry: expiry,
		}
	}
}

func (reader *RDBReader) readMillisecondsExpiryValue() (string, string, int64) {
	readByte, err := reader.bufReader.ReadByte()
	if err != nil {
		log.Fatalln(err)
	}
	if readByte != millisecondsExpiryOPCode {
		log.Fatalln("should not be in ms expiry")
	}
	var timeBytes []byte = make([]byte, 8)
	_, err = reader.bufReader.Read(timeBytes)
	var time int64
	err = binary.Read(bytes.NewBuffer(timeBytes), binary.LittleEndian, &time)
	if err != nil {
		log.Fatalln(err)
	}
	key, value := reader.readNoExpiryValue()
	log.Printf("Key : %s Value : %s Expiry : %d", key, value, time)
	return key, value, time
}

func (reader *RDBReader) readSecondsExpiryValue() (string, string, int64) {
	readByte, err := reader.bufReader.ReadByte()
	if err != nil {
		log.Fatalln(err)
	}
	if readByte != secondsExpiryOPCode {
		log.Fatalln("should not be in seconds expiry")
	}
	var timeBytes []byte = make([]byte, 4)
	_, err = reader.bufReader.Read(timeBytes)
	var time int64
	err = binary.Read(bytes.NewBuffer(timeBytes), binary.LittleEndian, &time)
	if err != nil {
		log.Fatalln(err)
	}
	key, value := reader.readNoExpiryValue()
	log.Printf("Key : %s Value : %s Expiry : %d", key, value, time)
	return key, value, time
}

func (reader *RDBReader) readNoExpiryValue() (string, string) {
	valueTypeByte, err := reader.bufReader.ReadByte()
	if err != nil {
		log.Fatalln(err)
	}
	valueType := int(valueTypeByte)
	bitCase, length := reader.decodeLength()
	if bitCase == 3 {
		log.Fatalln("Unexpected bitcase. String encoded value cannot have case 3")
	}
	var keyBuffer []byte = make([]byte, int(length))
	_, err = reader.bufReader.Read(keyBuffer)
	if err != nil {
		log.Fatalln(err)
	}
	key := string(keyBuffer)
	log.Printf("Key : %s", key)
	switch valueType {
	case 0:
		//string value
		bitCase, length = reader.decodeLength()
		if bitCase == 3 {
			log.Fatalln("Unexpected bitcase. String encoded value cannot have case 3")
		}
		var valueBuffer []byte = make([]byte, int(length))
		_, err = reader.bufReader.Read(valueBuffer)
		if err != nil {
			log.Fatalln(err)
		}
		value := string(valueBuffer)
		log.Printf("value : %s", value)
		return key, value
	default:
		log.Println("value type read not implemented yet")
	}
	return "", ""
}

func (reader *RDBReader) ReadKeys() []string {
	var result map[string]models.DataNode = make(map[string]models.DataNode)
	reader.ReadDatabase(result)
	var keys []string
	for key := range result {
		keys = append(keys, key)
	}
	return keys
}

func (reader *RDBReader) decodeLength() (BitCase, Length) {
	firstByte, err := reader.bufReader.ReadByte()
	if err != nil {
		log.Fatalln(err)
	}
	//checking the first two msbs
	firstByteShifted := firstByte >> 6
	switch int8(firstByteShifted) {
	case 0:
		//The next 6 bits represent the length
		stringLength := (firstByte << 2) >> 2
		return 0, Length(stringLength)
	case 1:
		//Read one additional byte. The combined 14 bits represent the length
		additionalByte, err := reader.bufReader.ReadByte()
		remainingBits := (firstByte << 2) >> 2
		if err != nil {
			log.Fatalln(err)
		}
		var result uint16 = uint16(remainingBits)
		result = result << 8
		result = result | uint16(additionalByte)
		return 1, Length(result)
	case 2:
		//Discard the remaining 6 bits. The next 4 bytes from the stream represent the length
		var nextFourBytes []byte = make([]byte, 4)
		_, err := reader.bufReader.Read(nextFourBytes)
		if err != nil {
			log.Fatalln(err)
		}
		var length int32
		err = binary.Read(bytes.NewBuffer(nextFourBytes), binary.BigEndian, &length)
		if err != nil {
			log.Fatalln(err)
		}
		return 2, Length(length)
	case 3:
		//The next object is encoded in a special format. The remaining 6 bits indicate the format
		remainingBits := (firstByte << 2) >> 2
		return 3, Length(remainingBits)
	}
	return 0, -1
}

func (reader *RDBReader) readLengthPrefixedString(stringLength int) string {
	var buffer []byte = make([]byte, stringLength)
	_, err := reader.bufReader.Read(buffer)
	if err != nil {
		log.Fatalln(err)
	}
	return string(buffer)
}

func (reader *RDBReader) readInteger(length int) int {
	var buffer []byte
	if length == 0 {
		//8 bit int
		buffer = make([]byte, 1)
		_, err := reader.bufReader.Read(buffer)
		if err != nil {
			log.Fatalln(err)
		}
		var result int8
		err = binary.Read(bytes.NewBuffer(buffer), binary.LittleEndian, &result)
		if err != nil {
			log.Fatalln(err)
		}
		return int(result)
	} else if length == 1 {
		//8 bit int
		buffer = make([]byte, 2)
		_, err := reader.bufReader.Read(buffer)
		if err != nil {
			log.Fatalln(err)
		}
		var result int16
		err = binary.Read(bytes.NewBuffer(buffer), binary.LittleEndian, &result)
		if err != nil {
			log.Fatalln(err)
		}
		return int(result)
	} else if length == 2 {
		//8 bit int
		buffer = make([]byte, 4)
		_, err := reader.bufReader.Read(buffer)
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
