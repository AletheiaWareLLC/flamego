package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*Xor)(nil)
var _ Emittable = (*Xor)(nil)

type Xor struct {
	Statement
	source1     flamego.Register
	source2     flamego.Register
	destination flamego.Register
}

func NewXor(s1, s2, d flamego.Register, c string) *Xor {
	return &Xor{
		Statement: Statement{
			comment: c,
		},
		source1:     s1,
		source2:     s2,
		destination: d,
	}
}

func (a *Xor) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Xor) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Xor) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Xor) Instruction() flamego.Instruction {
	return isa.NewXor(a.source1, a.source2, a.destination)
}
