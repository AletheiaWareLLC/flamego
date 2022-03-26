package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*Or)(nil)
var _ Emittable = (*Or)(nil)

type Or struct {
	Statement
	source1     flamego.Register
	source2     flamego.Register
	destination flamego.Register
}

func NewOr(s1, s2, d flamego.Register, c string) *Or {
	return &Or{
		Statement: Statement{
			comment: c,
		},
		source1:     s1,
		source2:     s2,
		destination: d,
	}
}

func (a *Or) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Or) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Or) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Or) Instruction() flamego.Instruction {
	return isa.NewOr(a.source1, a.source2, a.destination)
}
