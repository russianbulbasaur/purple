package client

import (
	"bytes"
	"log"
	"purple/internals/my_resp/constants"
	types "purple/internals/my_resp/purple_data_types"
	arrayTypes "purple/internals/my_resp/purple_data_types/array"
)

func (client *Client) evaluateArray(purpleArray arrayTypes.PurpleArray) {
	switch purpleArray.GetType() {
	case constants.PurpleBooleanArrayType:
		fallthrough
	case constants.PurpleStringArrayType:
		client.evaluateStringArray(purpleArray)
	}
}

func (client *Client) evaluateStringArray(purpleArray arrayTypes.PurpleArray) {
	var response *bytes.Buffer = bytes.NewBuffer(make([]byte, 0))
	arrayLength := purpleArray.GetLen()
	purpleStringArray := purpleArray.GetStruct().(*arrayTypes.PurpleStringArray)
	for i := 0; i < arrayLength; i++ {
		element, err := purpleStringArray.GetElementAt(i)
		log.Println(element.Value)
		if err != nil {
			return
		}
		if element.Value == "PING" {
			response.Write(client.resp.E.EncodeSimpleString("PONG"))
			client.writeChannel <- response.Bytes()
			return
		} else if element.Value == "ECHO" {
			parameter, err := purpleStringArray.GetElementAt(i + 1)
			if err != nil {
				return
			}
			response.Write(client.resp.E.EncodeSimpleString(parameter.Value))
			i++
		} else if element.Value == "GET" {

		} else if element.Value == "SET" {

		}
	}
	client.writeChannel <- response.Bytes()
}

func (client *Client) evaluateString(purpleString types.PurpleString) {
	var response *bytes.Buffer = bytes.NewBuffer(make([]byte, 0))
	if purpleString.Value == "PING" {
		response.Write(client.resp.E.EncodeSimpleString("PONG"))
		client.writeChannel <- response.Bytes()
		return
	} else {
		response.WriteString("Command not implmenented")
		client.writeChannel <- response.Bytes()
	}
}
