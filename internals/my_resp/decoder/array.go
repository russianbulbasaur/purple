package decoder

import (
	"bytes"
	"errors"
	"log"
	"purple/internals/my_resp/constants"
	types "purple/internals/my_resp/purple_data_types"
	arrayTypes "purple/internals/my_resp/purple_data_types/array"
	"strconv"
)

func decodeArray(input []byte) (arrayTypes.PurpleArray, error, uint64) {
	var counter uint64 = 1 //pehla byte
	if !isValid(input) {
		log.Println("CR and LF bytes not in place")
		return nil, errors.New("CR and LF bytes not in place"), counter
	}
	if input[0] != constants.ArrayPrefix {
		log.Println("Invalid array prefix")
		return nil, errors.New("invalid array prefix"), counter
	}
	numberOfElementsBuffer := bytes.NewBuffer(make([]byte, 0))
	var pointer uint64 = 1
	for _, b := range input[pointer:] {
		pointer++
		if b == constants.CR {
			pointer++ // lf flag accounting
			break
		}
		numberOfElementsBuffer.WriteByte(b)
	}
	numberOfElements, err := strconv.Atoi(string(numberOfElementsBuffer.Bytes()))
	if err != nil {
		log.Println("Invalid number of elements")
		return nil, errors.New("invalid number of elements"), counter
	}
	log.Printf("Number of elements %d", numberOfElements)
	decoder := MyRespDecoder{}
	var result arrayTypes.PurpleArray = nil
	for pointer < uint64(len(input)) {
		purpleDataType, err, decodeCount := decoder.Decode(input[pointer:])
		if err != nil {
			log.Println(err)
			break
		}
		pointer += decodeCount
		switch purpleDataType.(type) {
		case types.PurpleString:
			if result == nil {
				result = &arrayTypes.PurpleStringArray{}
			}
			result.AddElement(purpleDataType)
		case types.PurpleBoolean:
			if result == nil {
				result = &arrayTypes.PurpleBooleanArray{}
			}
			result.AddElement(purpleDataType)
		}
	}
	//need to push commit
	return result, nil, counter
}
