package decoder

import (
	"errors"
	"log"
	"purple/internals/my_resp/constants"
)

func decodeSimpleString(input []byte) (constants.PurpleString, error) {
	if !isValid(input) {
		log.Println("simple string : invalid lf and cr bytes")
		return "", errors.New("invalid lf and cr bytes")
	}

	return constants.PurpleString(string(input[1 : len(input)-2])), nil
}

//func decodeBulkString(input []byte) (constants.PurpleString, error) {
//	if !isValid(input) {
//		log.Println("invalid lf and cr bytes")
//		return "", errors.New("invalid lf and cr bytes")
//	}
//	lengthAndPayload := bytes.Split(input[1:], []byte("\r\n"))
//	if len(lengthAndPayload) < 2 {
//		log.Println("Invalid bulk string")
//		return "", errors.New("invalid bulk string")
//	}
//	bulkLength, err := strconv.Atoi(string(lengthAndPayload[0]))
//	if err != nil {
//		log.Println(err)
//		return "", err
//	}
//	payload := string(lengthAndPayload[1])
//	if len(payload) != bulkLength {
//		log.Println("Payload length and specified length do not match")
//		return "", errors.New("payload length and specified length do not match")
//	}
//	return constants.PurpleString(payload), nil
//}
