package array

import (
	"purple/internals/my_resp/constants"
	types "purple/internals/my_resp/purple_data_types"
)

type PurpleBooleanArray struct {
	elements []types.PurpleBoolean
	len      uint64
}

func (p *PurpleBooleanArray) AddElement(i interface{}) {
	p.elements = append(p.elements, i.(types.PurpleBoolean))
}

func (p *PurpleBooleanArray) GetLen() int {
	return len(p.elements)
}

func (p *PurpleBooleanArray) GetType() int {
	return constants.PurpleBooleanArrayType
}

func (p *PurpleBooleanArray) GetStruct() interface{} {
	return p
}
