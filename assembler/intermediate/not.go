package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*Not)(nil)
var _ Emittable = (*Not)(nil)

type Not struct {
	Statement
	source      flamego.Register
	destination flamego.Register
}

func NewNot(s, d flamego.Register, c string) *Not {
	return &Not{
		Statement: Statement{
			comment: c,
		},
		source:      s,
		destination: d,
	}
}

func (a *Not) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Not) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Not) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Not) Instruction() flamego.Instruction {
	return isa.NewNot(a.source, a.destination)
}
