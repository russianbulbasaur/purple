package commands

import (
	"bytes"
	resp "purple/internals/my_resp"
)

func Ping(resp *resp.MyRespObject) []byte {
	var response *bytes.Buffer = bytes.NewBuffer(make([]byte, 0))
	response.Write(resp.E.EncodeSimpleString("PONG"))
	return response.Bytes()
}
