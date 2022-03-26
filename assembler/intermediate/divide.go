package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*Divide)(nil)
var _ Emittable = (*Divide)(nil)

type Divide struct {
	Statement
	source1     flamego.Register
	source2     flamego.Register
	destination flamego.Register
}

func NewDivide(s1, s2, d flamego.Register, c string) *Divide {
	return &Divide{
		Statement: Statement{
			comment: c,
		},
		source1:     s1,
		source2:     s2,
		destination: d,
	}
}

func (a *Divide) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Divide) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Divide) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Divide) Instruction() flamego.Instruction {
	return isa.NewDivide(a.source1, a.source2, a.destination)
}
