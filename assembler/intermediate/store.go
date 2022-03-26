package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*Store)(nil)
var _ Emittable = (*Store)(nil)

type Store struct {
	Statement
	address flamego.Register
	offset  uint32
	source  flamego.Register
}

func NewStore(a flamego.Register, o uint32, s flamego.Register, c string) *Store {
	return &Store{
		Statement: Statement{
			comment: c,
		},
		address: a,
		offset:  o,
		source:  s,
	}
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
	return isa.NewStore(a.address, a.offset, a.source)
}
