package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*Return)(nil)
var _ Emittable = (*Return)(nil)

type Return struct {
	Statement
}

func NewReturn(c string) *Return {
	return &Return{
		Statement: Statement{
			comment: c,
		},
	}
}

func (a *Return) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Return) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Return) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Return) Instruction() flamego.Instruction {
	return isa.NewReturn()
}
