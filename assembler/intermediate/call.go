package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*Call)(nil)
var _ Emittable = (*Call)(nil)

type Call struct {
	Statement
	register flamego.Register
}

func NewCall(r flamego.Register, c string) *Call {
	return &Call{
		Statement: Statement{
			comment: c,
		},
		register: r,
	}
}

func (a *Call) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Call) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Call) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Call) Instruction() flamego.Instruction {
	return isa.NewCall(a.register)
}
