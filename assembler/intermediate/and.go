package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*And)(nil)
var _ Emittable = (*And)(nil)

type And struct {
	Statement
	source1     flamego.Register
	source2     flamego.Register
	destination flamego.Register
}

func NewAnd(s1, s2, d flamego.Register, c string) *And {
	return &And{
		Statement: Statement{
			comment: c,
		},
		source1:     s1,
		source2:     s2,
		destination: d,
	}
}

func (a *And) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *And) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *And) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *And) Instruction() flamego.Instruction {
	return isa.NewAnd(a.source1, a.source2, a.destination)
}
