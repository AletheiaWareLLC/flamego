package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*Subtract)(nil)
var _ Emittable = (*Subtract)(nil)

type Subtract struct {
	Statement
	source1     flamego.Register
	source2     flamego.Register
	destination flamego.Register
}

func NewSubtract(s1, s2, d flamego.Register, c string) *Subtract {
	return &Subtract{
		Statement: Statement{
			comment: c,
		},
		source1:     s1,
		source2:     s2,
		destination: d,
	}
}

func (a *Subtract) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Subtract) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Subtract) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Subtract) Instruction() flamego.Instruction {
	return isa.NewSubtract(a.source1, a.source2, a.destination)
}
