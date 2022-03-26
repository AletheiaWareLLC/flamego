package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

type Flush struct {
	Statement
	address flamego.Register
	offset  uint32
}

func NewFlush(a flamego.Register, o uint32, c string) *Flush {
	return &Flush{
		Statement: Statement{
			comment: c,
		},
		address: a,
		offset:  o,
	}
}

func (a *Flush) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Flush) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Flush) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Flush) Instruction() flamego.Instruction {
	return isa.NewFlush(a.address, a.offset)
}
