package encoder

import (
	"bytes"
	"purple/internals/my_resp/constants"
	"strconv"
)

func (encoder *MyRespEncoder) EncodeBulkString(input string) []byte {
	response := bytes.NewBuffer(make([]byte, 0))
	response.WriteByte(constants.BulkStringPrefix)
	stringLen := len(input)
	if stringLen == 0 {
		//null string
		response.Write([]byte{'-', '1'})
		response.Write([]byte{constants.CR, constants.LF})
		return response.Bytes()
	}
	response.WriteString(strconv.Itoa(stringLen))
	response.WriteString(input)
	response.Write([]byte{constants.CR, constants.LF})
	return response.Bytes()
}

func (encoder *MyRespEncoder) EncodeSimpleString(input string) []byte {
	response := bytes.NewBuffer(make([]byte, 0))
	response.WriteByte(constants.SimpleStringPrefix)
	response.WriteString(input)
	response.Write([]byte{constants.CR, constants.LF})
	return response.Bytes()
}
