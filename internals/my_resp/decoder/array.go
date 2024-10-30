package decoder

import (
	"bytes"
	"log"
	"purple/internals/my_resp/constants"
	"strconv"
)

func decodeArray(input []byte) (interface{}, string) {
	if !isValid(input) {
		log.Println("CR and LF bytes not in place")
		return nil, ""
	}
	if input[0] != constants.ArrayPrefix {
		log.Println("Invalid array prefix")
		return nil, ""
	}
	numberOfElementsBuffer := bytes.NewBuffer(make([]byte, 0))
	counter := 1
	for _, b := range input[counter:] {
		counter++
		if b == constants.CR {
			counter++ // lf flag accounting
			break
		}
		numberOfElementsBuffer.WriteByte(b)
	}
	numberOfElements, err := strconv.Atoi(string(numberOfElementsBuffer.Bytes()))
	if err != nil {
		log.Println("Invalid number of elements")
		return nil, ""
	}
	log.Printf("Number of elements %d", numberOfElements)
	decoder := MyRespDecoder{}
	log.Println(string(input[counter:]))
	decoder.Decode(input[counter:])
	return nil, ""
}
