package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*Push)(nil)
var _ Emittable = (*Push)(nil)

type Push struct {
	Statement
	mask uint16
}

func NewPush(m uint16, c string) *Push {
	return &Push{
		Statement: Statement{
			comment: c,
		},
		mask: m,
	}
}

func (a *Push) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Push) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Push) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Push) Instruction() flamego.Instruction {
	return isa.NewPush(a.mask)
}
