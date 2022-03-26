package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

type Clear struct {
	Statement
	address flamego.Register
	offset  uint32
}

func NewClear(a flamego.Register, o uint32, c string) *Clear {
	return &Clear{
		Statement: Statement{
			comment: c,
		},
		address: a,
		offset:  o,
	}
}

func (a *Clear) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Clear) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Clear) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Clear) Instruction() flamego.Instruction {
	return isa.NewClear(a.address, a.offset)
}
