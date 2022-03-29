package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type JumpConditionCode uint8

const (
	JumpEZ JumpConditionCode = iota
	JumpNZ
	JumpLE
	JumpLZ
)

func (i JumpConditionCode) String() string {
	switch i {
	case JumpEZ:
		return "ez"
	case JumpNZ:
		return "nz"
	case JumpLE:
		return "le"
	case JumpLZ:
		return "lz"
	}
	return "Unrecognized Jump Condition Code"
}

type JumpDirection bool

const (
	JumpForward  JumpDirection = false
	JumpBackward JumpDirection = true
)

func (i JumpDirection) String() string {
	if i {
		return "-"
	}
	return "+"
}

type Jump struct {
	ConditionCode     JumpConditionCode
	Direction         JumpDirection
	Offset            uint32
	ConditionRegister flamego.Register
}

func NewJump(cc JumpConditionCode, d JumpDirection, o uint32, r flamego.Register) *Jump {
	return &Jump{
		ConditionCode:     cc,
		Direction:         d,
		Offset:            o,
		ConditionRegister: r,
	}
}

func (i *Jump) Load(x flamego.Context) (uint64, uint64, uint64) {
	// Load Program Counter
	a := x.ReadRegister(flamego.RProgramCounter)
	// Load Condition Registers
	b := x.ReadRegister(i.ConditionRegister)
	return a, b, 0
}

func (i *Jump) Execute(x flamego.Context, a, b, c uint64) uint64 {
	jump := false
	switch i.ConditionCode {
	case JumpEZ:
		jump = b == 0
	case JumpNZ:
		jump = b != 0
	case JumpLE:
		jump = b <= 0
	case JumpLZ:
		jump = b < 0
	}
	if jump {
		offset := uint64(i.Offset)
		switch i.Direction {
		case JumpForward:
			return a + offset
		case JumpBackward:
			return a - offset
		}
	}
	return a + flamego.InstructionSize
}

func (i *Jump) Format(x flamego.Context, a uint64) uint64 {
	// Pass Through
	return a
}

func (i *Jump) Store(x flamego.Context, a uint64) {
	// Update Program Counter
	x.SetProgramCounter(a)
}

func (i *Jump) Retire(x flamego.Context) bool {
	// Do Nothing
	return true
}

func (i *Jump) String() string {
	return fmt.Sprintf("j%s %s %s0x%x", i.ConditionCode, i.ConditionRegister, i.Direction, i.Offset)
}
