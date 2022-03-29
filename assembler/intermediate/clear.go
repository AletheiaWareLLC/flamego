package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*Clear)(nil)
var _ Emittable = (*Clear)(nil)
var _ Linkable = (*Clear)(nil)

type Clear struct {
	Statement
	address      flamego.Register
	labelName    string
	constantName string
	offset       uint32
	label        *Label
	constant     *Data
}

func NewClearWithLabel(a flamego.Register, l string, c string) *Clear {
	return &Clear{
		Statement: Statement{
			comment: c,
		},
		address:   a,
		labelName: l,
	}
}

func NewClearWithConstant(a flamego.Register, n string, c string) *Clear {
	return &Clear{
		Statement: Statement{
			comment: c,
		},
		address:      a,
		constantName: n,
	}
}

func NewClearWithOffset(a flamego.Register, o uint32, c string) *Clear {
	return &Clear{
		Statement: Statement{
			comment: c,
		},
		address: a,
		offset:  o,
	}
}

func (a *Clear) Link(l Linker) error {
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
	if a.label != nil {
		a.offset = uint32(a.label.AbsoluteAddress())
	} else if a.constant != nil {
		a.offset = uint32(a.constant.Value())
	}
	return isa.NewClear(a.address, a.offset)
}
