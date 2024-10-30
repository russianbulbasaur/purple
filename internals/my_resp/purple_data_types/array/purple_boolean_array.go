package array

import (
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
