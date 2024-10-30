package encoder

import (
	"bytes"
	"purple/internals/my_resp/constants"
	"strconv"
	"strings"
)

func (encoder *MyRespEncoder) EncodeInteger(input int) []byte {
	response := bytes.NewBuffer(make([]byte, 0))
	response.WriteByte(constants.IntegerPrefix)
	inputString := strings.Trim(strconv.Itoa(input), "-+")
	if input < 0 {
		response.WriteByte('-')
	} else {
		response.WriteByte('+')
	}
	response.WriteString(inputString)
	response.Write([]byte{constants.CR, constants.LF})
	return response.Bytes()
}
