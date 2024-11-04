package commands

import (
	"bytes"
	resp "purple/internals/my_resp"
)

func Echo(resp *resp.MyRespObject, value string) []byte {
	var response *bytes.Buffer = bytes.NewBuffer(make([]byte, 0))
	response.Write(resp.E.EncodeSimpleString(value))
	return response.Bytes()
}
