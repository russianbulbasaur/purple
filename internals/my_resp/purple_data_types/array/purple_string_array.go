package array

import (
	types "purple/internals/my_resp/purple_data_types"
)

type PurpleStringArray struct {
	elements []types.PurpleString
	len      uint64
}

func (p *PurpleStringArray) AddElement(i interface{}) {
	p.elements = append(p.elements, i.(types.PurpleString))
}

func (p *PurpleStringArray) GetLen() int {
	return len(p.elements)
}
