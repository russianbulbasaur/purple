package my_resp

import (
	"purple/internals/my_resp/decoder"
	"purple/internals/my_resp/encoder"
)

type MyRespObject struct {
	E *encoder.MyRespEncoder
	D *decoder.MyRespDecoder
}

func (ob *MyRespObject) newEncoder() {
	ob.E = &encoder.MyRespEncoder{}
}

func (ob *MyRespObject) newDecoder() {
	ob.D = &decoder.MyRespDecoder{}
}
