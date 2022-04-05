package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*Pop)(nil)
var _ Emittable = (*Pop)(nil)

type Pop struct {
	Statement
	mask uint16
}

func NewPop(m uint16, c string) *Pop {
	return &Pop{
		Statement: Statement{
			comment: c,
		},
		mask: m,
	}
}

func (a *Pop) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Pop) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Pop) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Pop) Instruction() flamego.Instruction {
	return isa.NewPop(a.mask)
}
