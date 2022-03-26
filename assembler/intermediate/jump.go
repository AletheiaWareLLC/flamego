package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*Jump)(nil)
var _ Emittable = (*Jump)(nil)
var _ Linkable = (*Jump)(nil)

type Jump struct {
	Statement
	condition   isa.JumpConditionCode
	destination string
	register    flamego.Register
	label       *Label
}

func NewJump(cc isa.JumpConditionCode, d string, r flamego.Register, c string) *Jump {
	return &Jump{
		Statement: Statement{
			comment: c,
		},
		condition:   cc,
		destination: d,
		register:    r,
	}
}

func (a *Jump) Link(l Linker) error {
	label, err := l.Label(a.destination)
	if err != nil {
		return err
	}
	a.label = label
	return nil
}

func (a *Jump) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Jump) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Jump) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Jump) Instruction() flamego.Instruction {
	o := a.label.RelativeAddress(a.address)
	d := isa.JumpForward
	if o < 0 {
		d = isa.JumpBackward
		o = -o
	}
	return isa.NewJump(a.condition, d, uint32(o&0x3FFFFF), a.register)
}
