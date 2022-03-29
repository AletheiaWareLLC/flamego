package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*Load)(nil)
var _ Emittable = (*Load)(nil)
var _ Linkable = (*Load)(nil)

type Load struct {
	Statement
	address      flamego.Register
	labelName    string
	constantName string
	offset       uint32
	destination  flamego.Register
	label        *Label
	constant     *Data
}

func NewLoadWithLabel(a flamego.Register, l string, d flamego.Register, c string) *Load {
	return &Load{
		Statement: Statement{
			comment: c,
		},
		address:     a,
		labelName:   l,
		destination: d,
	}
}

func NewLoadWithConstant(a flamego.Register, n string, d flamego.Register, c string) *Load {
	return &Load{
		Statement: Statement{
			comment: c,
		},
		address:      a,
		constantName: n,
		destination:  d,
	}
}

func NewLoadWithOffset(a flamego.Register, o uint32, d flamego.Register, c string) *Load {
	return &Load{
		Statement: Statement{
			comment: c,
		},
		address:     a,
		offset:      o,
		destination: d,
	}
}

func (a *Load) Link(l Linker) error {
	if a.labelName != "" {
		label, err := l.Label(a.labelName)
		if err != nil {
			return err
		}
		a.label = label
	} else if a.constantName != "" {
		constant, err := l.Constant(a.constantName)
		if err != nil {
			return err
		}
		a.constant = constant
	}
	return nil
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
	if a.label != nil {
		a.offset = uint32(a.label.AbsoluteAddress())
	} else if a.constant != nil {
		a.offset = uint32(a.constant.Value())
	}
	return isa.NewLoad(a.address, a.offset, a.destination)
}
