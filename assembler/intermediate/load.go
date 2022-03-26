package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*Load)(nil)
var _ Emittable = (*Load)(nil)

type Load struct {
	Statement
	address     flamego.Register
	offset      uint32
	destination flamego.Register
}

func NewLoad(a flamego.Register, o uint32, d flamego.Register, c string) *Load {
	return &Load{
		Statement: Statement{
			comment: c,
		},
		address:     a,
		offset:      o,
		destination: d,
	}
}

func (a *Load) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Load) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Load) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Load) Instruction() flamego.Instruction {
	return isa.NewLoad(a.address, a.offset, a.destination)
}
