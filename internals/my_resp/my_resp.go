package my_resp

func Init() *MyRespObject {
	myRespObject := &MyRespObject{}
	myRespObject.newEncoder()
	myRespObject.newDecoder()
	return myRespObject
}
