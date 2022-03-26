package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*Modulo)(nil)
var _ Emittable = (*Modulo)(nil)

type Modulo struct {
	Statement
	source1     flamego.Register
	source2     flamego.Register
	destination flamego.Register
}

func NewModulo(s1, s2, d flamego.Register, c string) *Modulo {
	return &Modulo{
		Statement: Statement{
			comment: c,
		},
		source1:     s1,
		source2:     s2,
		destination: d,
	}
}

func (a *Modulo) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Modulo) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Modulo) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Modulo) Instruction() flamego.Instruction {
	return isa.NewModulo(a.source1, a.source2, a.destination)
}
