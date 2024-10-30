package decoder

import (
	"purple/internals/my_resp/constants"
)

type MyRespDecoder struct {
}

func (decoder *MyRespDecoder) Decode(input []byte) (interface{}, error, uint64) {
	switch input[0] {
	case constants.SimpleStringPrefix:
		return decodeSimpleString(input)
	case constants.BulkStringPrefix:
		return decodeBulkString(input)
	case constants.IntegerPrefix:
		return decodeInteger(input)
	//case constants.ErrorPrefix:
	//	return decodeError(input)
	case constants.BooleanPrefix:
		return decodeBoolean(input)
	case constants.ArrayPrefix:
		return decodeArray(input)
	default:
		panic("Invalid input : " + string(input))
	}
}

func isValid(input []byte) bool {
	if len(input) < 3 {
		return false
	}
	lfIndex := len(input) - 1
	crIndex := lfIndex - 1
	if input[lfIndex] != constants.LF || input[crIndex] != constants.CR {
		return false
	}
	return true
}
