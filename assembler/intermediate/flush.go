package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*Flush)(nil)
var _ Emittable = (*Flush)(nil)
var _ Linkable = (*Flush)(nil)

type Flush struct {
	Statement
	address      flamego.Register
	labelName    string
	constantName string
	offset       uint32
	label        *Label
	constant     *Data
}

func NewFlushWithLabel(a flamego.Register, l string, c string) *Flush {
	return &Flush{
		Statement: Statement{
			comment: c,
		},
		address:   a,
		labelName: l,
	}
}

func NewFlushWithConstant(a flamego.Register, n string, c string) *Flush {
	return &Flush{
		Statement: Statement{
			comment: c,
		},
		address:      a,
		constantName: n,
	}
}

func NewFlushWithOffset(a flamego.Register, o uint32, c string) *Flush {
	return &Flush{
		Statement: Statement{
			comment: c,
		},
		address: a,
		offset:  o,
	}
}

func (a *Flush) Link(l Linker) error {
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
	if a.label != nil {
		a.offset = uint32(a.label.AbsoluteAddress())
	} else if a.constant != nil {
		a.offset = uint32(a.constant.Value())
	}
	return isa.NewFlush(a.address, a.offset)
}
