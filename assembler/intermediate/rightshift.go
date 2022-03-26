package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*RightShift)(nil)
var _ Emittable = (*RightShift)(nil)

type RightShift struct {
	Statement
	source1     flamego.Register
	source2     flamego.Register
	destination flamego.Register
}

func NewRightShift(s1, s2, d flamego.Register, c string) *RightShift {
	return &RightShift{
		Statement: Statement{
			comment: c,
		},
		source1:     s1,
		source2:     s2,
		destination: d,
	}
}

func (a *RightShift) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *RightShift) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *RightShift) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *RightShift) Instruction() flamego.Instruction {
	return isa.NewRightShift(a.source1, a.source2, a.destination)
}
