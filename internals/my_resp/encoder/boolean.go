package encoder

import (
	"bytes"
	"purple/internals/my_resp/constants"
)

func (encoder *MyRespEncoder) EncodeBoolean(input bool) []byte {
	response := bytes.NewBuffer(make([]byte, 0))
	response.WriteByte(constants.BooleanPrefix)
	if input {
		response.WriteByte(constants.True)
	} else {
		response.WriteByte(constants.False)
	}
	response.Write([]byte{constants.CR, constants.LF})
	return response.Bytes()
}
