package decoder

import (
	"errors"
	"log"
	"purple/internals/my_resp/constants"
	types "purple/internals/my_resp/purple_data_types"
)

func decodeBoolean(input []byte) (types.PurpleBoolean, error, uint64) {
	purpleBoolean := types.PurpleBoolean{}
	var counter uint64 = 1 //pehla byte
	if !isValid(input) {
		log.Println("Invalid cr lf bytes")
		return purpleBoolean, errors.New("invalid cr lf bytes"), counter
	}
	respBool := input[1]
	counter += 3 //<t/f>\r\n
	var response bool
	if respBool == constants.True {
		response = true
	} else if respBool == constants.False {
		response = false
	} else {
		panic("invalid boolean type")
	}
	purpleBoolean.Value = response
	return purpleBoolean, nil, counter
}
