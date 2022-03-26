package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*Add)(nil)
var _ Emittable = (*Add)(nil)

type Add struct {
	Statement
	source1     flamego.Register
	source2     flamego.Register
	destination flamego.Register
}

func NewAdd(s1, s2, d flamego.Register, c string) *Add {
	return &Add{
		Statement: Statement{
			comment: c,
		},
		source1:     s1,
		source2:     s2,
		destination: d,
	}
}

func (a *Add) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Add) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Add) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Add) Instruction() flamego.Instruction {
	return isa.NewAdd(a.source1, a.source2, a.destination)
}
