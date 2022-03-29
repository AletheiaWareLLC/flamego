package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*Store)(nil)
var _ Emittable = (*Store)(nil)
var _ Linkable = (*Store)(nil)

type Store struct {
	Statement
	address      flamego.Register
	labelName    string
	constantName string
	offset       uint32
	source       flamego.Register
	label        *Label
	constant     *Data
}

func NewStoreWithLabel(a flamego.Register, l string, s flamego.Register, c string) *Store {
	return &Store{
		Statement: Statement{
			comment: c,
		},
		address:   a,
		labelName: l,
		source:    s,
	}
}

func NewStoreWithConstant(a flamego.Register, n string, s flamego.Register, c string) *Store {
	return &Store{
		Statement: Statement{
			comment: c,
		},
		address:      a,
		constantName: n,
		source:       s,
	}
}

func NewStoreWithOffset(a flamego.Register, o uint32, s flamego.Register, c string) *Store {
	return &Store{
		Statement: Statement{
			comment: c,
		},
		address: a,
		offset:  o,
		source:  s,
	}
}

func (a *Store) Link(l Linker) error {
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

func (a *Store) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Store) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Store) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Store) Instruction() flamego.Instruction {
	if a.label != nil {
		a.offset = uint32(a.label.AbsoluteAddress())
	} else if a.constant != nil {
		a.offset = uint32(a.constant.Value())
	}
	return isa.NewStore(a.address, a.offset, a.source)
}
