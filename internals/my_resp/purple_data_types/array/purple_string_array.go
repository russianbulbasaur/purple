package array

import (
	"errors"
	"log"
	"purple/internals/my_resp/constants"
	types "purple/internals/my_resp/purple_data_types"
)

type PurpleStringArray struct {
	elements []types.PurpleString
	len      int
}

func (p *PurpleStringArray) GetType() int {
	return constants.PurpleStringArrayType
}

func (p *PurpleStringArray) GetStruct() interface{} {
	return p
}

func (p *PurpleStringArray) AddElement(i interface{}) {
	p.elements = append(p.elements, i.(types.PurpleString))
	p.len = len(p.elements)
}

func (p *PurpleStringArray) GetLen() int {
	return p.len
}

func (p *PurpleStringArray) GetElements() []types.PurpleString {
	return p.elements
}

func (p *PurpleStringArray) GetElementAt(index int) (types.PurpleString, error) {
	if index >= p.len {
		log.Println("Index out of range")
		return types.PurpleString{}, errors.New("Index out of range")
	}
	return p.elements[index], nil
}
