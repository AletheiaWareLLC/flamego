package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*LeftShift)(nil)
var _ Emittable = (*LeftShift)(nil)

type LeftShift struct {
	Statement
	source1     flamego.Register
	source2     flamego.Register
	destination flamego.Register
}

func NewLeftShift(s1, s2, d flamego.Register, c string) *LeftShift {
	return &LeftShift{
		Statement: Statement{
			comment: c,
		},
		source1:     s1,
		source2:     s2,
		destination: d,
	}
}

func (a *LeftShift) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *LeftShift) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *LeftShift) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *LeftShift) Instruction() flamego.Instruction {
	return isa.NewLeftShift(a.source1, a.source2, a.destination)
}
