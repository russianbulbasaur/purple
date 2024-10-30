package encoder

import (
	"bytes"
	"purple/internals/my_resp/constants"
)

func (encoder *MyRespEncoder) EncodeSimpleError(error string) []byte {
	response := bytes.NewBuffer(make([]byte, 0))
	response.WriteByte(constants.ErrorPrefix)
	response.WriteString(error)
	response.Write([]byte{constants.CR, constants.LF})
	return response.Bytes()
}
