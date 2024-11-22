package my_resp

import (
	"purple/internals/my_resp/decoder"
	"purple/internals/my_resp/encoder"
)

type MyRespObject struct {
	E *encoder.MyRespEncoder
	D *decoder.MyRespDecoder
}

func Init() *MyRespObject {
	myRespObject := &MyRespObject{}
	myRespObject.E = &encoder.MyRespEncoder{}
	myRespObject.D = &decoder.MyRespDecoder{}
	return myRespObject
}
