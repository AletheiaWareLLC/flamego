package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*LoadC)(nil)
var _ Emittable = (*LoadC)(nil)
var _ Linkable = (*LoadC)(nil)

type LoadC struct {
	Statement
	labelName    string
	constantName string
	number       uint32
	register     flamego.Register
	label        *Label
	constant     *Data
}

func NewLoadCWithLabel(l string, r flamego.Register, c string) *LoadC {
	return &LoadC{
		Statement: Statement{
			comment: c,
		},
		labelName: l,
		register:  r,
	}
}

func NewLoadCWithConstant(n string, r flamego.Register, c string) *LoadC {
	return &LoadC{
		Statement: Statement{
			comment: c,
		},
		constantName: n,
		register:     r,
	}
}

func NewLoadCWithNumber(n uint32, r flamego.Register, c string) *LoadC {
	return &LoadC{
		Statement: Statement{
			comment: c,
		},
		number:   n,
		register: r,
	}
}

func (a *LoadC) Link(l Linker) error {
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

func (a *LoadC) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *LoadC) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *LoadC) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *LoadC) Instruction() flamego.Instruction {
	if a.label != nil {
		a.number = uint32(a.label.AbsoluteAddress())
	} else if a.constant != nil {
		a.number = uint32(a.constant.Value())
	}
	return isa.NewLoadC(a.number, a.register)
}
