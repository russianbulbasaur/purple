package decoder

import (
	"bytes"
	"errors"
	"log"
	types "purple/internals/my_resp/purple_data_types"
	"strconv"
)

func decodeSimpleString(input []byte) (types.PurpleString, error, uint64) {
	var counter uint64 = 1 //first byte ignore
	purpleString := types.PurpleString{}
	if !isValid(input) {
		log.Println("simple string : invalid lf and cr bytes")
		return purpleString, errors.New("invalid lf and cr bytes"), counter
	}
	purpleString.Value = string(input[1 : len(input)-2])
	purpleString.Len = uint64(len(purpleString.Value))
	counter += purpleString.Len
	counter += 2 //end ke CR and LF bytes
	return purpleString, nil, counter
}

func decodeBulkString(input []byte) (types.PurpleString, error, uint64) {
	var counter uint64 = 1 //ignore the first byte
	log.Println(string(input))
	purpleString := types.PurpleString{}
	if !isValid(input) {
		log.Println("invalid lf and cr bytes")
		return purpleString, errors.New("invalid lf and cr bytes"), counter
	}
	lengthAndPayload := bytes.Split(input[1:], []byte("\r\n"))
	if len(lengthAndPayload) < 2 {
		log.Println("Invalid bulk string")
		return purpleString, errors.New("invalid bulk string"), counter
	}
	lengthString := string(lengthAndPayload[0])
	counter += uint64(len(lengthString)) + 2
	bulkLength, err := strconv.Atoi(lengthString)
	if err != nil {
		log.Println(err)
		return purpleString, err, counter
	}
	payload := string(lengthAndPayload[1])
	if len(payload) != bulkLength {
		log.Println("Payload length and specified length do not match")
		return purpleString, errors.New("payload length and specified length do not match"), counter
	}
	counter += uint64(len(payload)) + 2
	purpleString.Value = payload
	purpleString.Len = uint64(bulkLength)
	return purpleString, nil, counter
}
