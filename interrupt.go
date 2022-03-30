package flamego

import (
	"fmt"
)

type InterruptValue int16

const (
	InterruptSignal InterruptValue = iota
	InterruptBreakpoint
	InterruptUnsupportedOperationError
	InterruptArithmeticError
	InterruptRegisterAccessError
	InterruptMemoryAccessError
	InterruptProgramAccessError
	InterruptStackOverflowError
	InterruptStackUnderflowError
)

const InterruptCount = 9

func (i InterruptValue) String() string {
	return fmt.Sprintf("Interrupt 0x%04x", uint16(i))
}
