package client

import (
	"bytes"
	"log"
	"math"
	"purple/internals/my_resp/constants"
	types "purple/internals/my_resp/purple_data_types"
	arrayTypes "purple/internals/my_resp/purple_data_types/array"
	"strconv"
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
			key, err := purpleStringArray.GetElementAt(i + 1)
			if err != nil {
				return
			}
			value := client.get(key.Value)
			if value == nil {
				response.WriteString("null")
			} else {
				response.Write(client.resp.E.EncodeSimpleString(value.(string)))
			}
			i++
		} else if element.Value == "SET" {
			key, err := purpleStringArray.GetElementAt(i + 1)
			if err != nil {
				return
			}
			value, err := purpleStringArray.GetElementAt(i + 2)
			if err != nil {
				return
			}
			if purpleStringArray.GetLen() > i+3 {
				next, err := purpleStringArray.GetElementAt(i + 3)
				if err != nil {
					return
				}
				if next.Value == "px" {
					m, err := purpleStringArray.GetElementAt(i + 4)
					if err != nil {
						return
					}
					milliseconds, err := strconv.ParseInt(m.Value, 10, 64)
					if err != nil {
						log.Println("invalid milliseconds")
						return
					}
					client.set(key.Value, value.Value, milliseconds/1000)
				}
			} else {
				client.set(key.Value, value.Value, math.MaxInt64)
			}
			i += 2
		} else if element.Value == "CONFIG" {
			key, err := purpleStringArray.GetElementAt(i + 1)
			if err != nil {
				return
			}
			if key.Value == "GET" {
				key, err := purpleStringArray.GetElementAt(i + 2)
				if err != nil {
					return
				}
				switch key.Value {
				case "dir":
					encoded := client.resp.E.EncodeStringArray(
						[]string{"dir", client.rdbFile.GetDir()})
					log.Println(encoded)
					response.Write(encoded)
				case "dbfilename":
					response.Write(client.resp.E.EncodeStringArray(
						[]string{"dbfilename", client.rdbFile.GetDBFileName()}))
				default:
					log.Printf("%s not found", key.Value)
				}
			} else {
				log.Println("Not implemented yet")
				return
			}
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
