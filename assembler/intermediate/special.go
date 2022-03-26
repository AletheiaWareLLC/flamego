package intermediate

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

var _ Addressable = (*Halt)(nil)
var _ Emittable = (*Halt)(nil)

type Halt struct {
	Statement
}

func NewHalt(c string) *Halt {
	return &Halt{
		Statement: Statement{
			comment: c,
		},
	}
}

func (a *Halt) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Halt) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Halt) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Halt) Instruction() flamego.Instruction {
	return isa.NewHalt()
}

var _ Addressable = (*Noop)(nil)
var _ Emittable = (*Noop)(nil)

type Noop struct {
	Statement
}

func NewNoop(c string) *Noop {
	return &Noop{
		Statement: Statement{
			comment: c,
		},
	}
}

func (a *Noop) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Noop) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Noop) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Noop) Instruction() flamego.Instruction {
	return isa.NewNoop()
}

var _ Addressable = (*Sleep)(nil)
var _ Emittable = (*Sleep)(nil)

type Sleep struct {
	Statement
}

func NewSleep(c string) *Sleep {
	return &Sleep{
		Statement: Statement{
			comment: c,
		},
	}
}

func (a *Sleep) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Sleep) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Sleep) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Sleep) Instruction() flamego.Instruction {
	return isa.NewSleep()
}

var _ Addressable = (*Signal)(nil)
var _ Emittable = (*Signal)(nil)

type Signal struct {
	Statement
	register flamego.Register
}

func NewSignal(r flamego.Register, c string) *Signal {
	return &Signal{
		Statement: Statement{
			comment: c,
		},
		register: r,
	}
}

func (a *Signal) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Signal) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Signal) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Signal) Instruction() flamego.Instruction {
	return isa.NewSignal(a.register)
}

var _ Addressable = (*Lock)(nil)
var _ Emittable = (*Lock)(nil)

type Lock struct {
	Statement
}

func NewLock(c string) *Lock {
	return &Lock{
		Statement: Statement{
			comment: c,
		},
	}
}

func (a *Lock) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Lock) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Lock) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Lock) Instruction() flamego.Instruction {
	return isa.NewLock()
}

var _ Addressable = (*Unlock)(nil)
var _ Emittable = (*Unlock)(nil)

type Unlock struct {
	Statement
}

func NewUnlock(c string) *Unlock {
	return &Unlock{
		Statement: Statement{
			comment: c,
		},
	}
}

func (a *Unlock) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Unlock) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Unlock) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Unlock) Instruction() flamego.Instruction {
	return isa.NewUnlock()
}

var _ Addressable = (*Interrupt)(nil)
var _ Emittable = (*Interrupt)(nil)

type Interrupt struct {
	Statement
	value flamego.InterruptValue
}

func NewInterrupt(v flamego.InterruptValue, c string) *Interrupt {
	return &Interrupt{
		Statement: Statement{
			comment: c,
		},
		value: v,
	}
}

func (a *Interrupt) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Interrupt) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Interrupt) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Interrupt) Instruction() flamego.Instruction {
	return isa.NewInterrupt(a.value)
}

var _ Addressable = (*Uninterrupt)(nil)
var _ Emittable = (*Uninterrupt)(nil)

type Uninterrupt struct {
	Statement
	register flamego.Register
}

func NewUninterrupt(r flamego.Register, c string) *Uninterrupt {
	return &Uninterrupt{
		Statement: Statement{
			comment: c,
		},
		register: r,
	}
}

func (a *Uninterrupt) String() string {
	return a.Instruction().String() + a.Statement.String()
}

func (a *Uninterrupt) Emit() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, isa.Encode(a.Instruction()))
	return buffer
}

func (a *Uninterrupt) EmittedSize() uint32 {
	return flamego.InstructionSize
}

func (a *Uninterrupt) Instruction() flamego.Instruction {
	return isa.NewUninterrupt(a.register)
}
