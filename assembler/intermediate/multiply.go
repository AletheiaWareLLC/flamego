package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*Multiply)(nil)
var _ Emittable = (*Multiply)(nil)

type Multiply struct {
	Statement
	source1     flamego.Register
	source2     flamego.Register
	destination flamego.Register
}

func NewMultiply(s1, s2, d flamego.Register, c string) *Multiply {
	return &Multiply{
		Statement: Statement{
			comment: c,
		},
		source1:     s1,
		source2:     s2,
		destination: d,
	}
}

func (a *Multiply) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Multiply) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Multiply) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Multiply) Instruction() flamego.Instruction {
	return isa.NewMultiply(a.source1, a.source2, a.destination)
}
