package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type Interrupt struct {
	Value flamego.InterruptValue
}

func NewInterrupt(value flamego.InterruptValue) *Interrupt {
	return &Interrupt{
		Value: value,
	}
}

func (i *Interrupt) Load(x flamego.Context) (uint64, uint64, uint64, uint64) {
	// Load Interrupt Vector Table
	return x.ReadRegister(flamego.RInterruptVectorTable), uint64(i.Value), 0, 0
}

func (i *Interrupt) Execute(x flamego.Context, a, b, c, d uint64) (uint64, uint64) {
	// Calculate address of Interrupt Service Routine by add interrupt value to Interrupt Vector Table
	return a + b, 0
}

func (i *Interrupt) Format(x flamego.Context, a, b uint64) (uint64, uint64) {
	// Do Nothing
	return a, 0
}

func (i *Interrupt) Store(x flamego.Context, a, b uint64) {
	// Jump to Interrupt Service Routine by updating the Program Counter
	x.SetProgramCounter(a)
}

func (i *Interrupt) Retire(x flamego.Context) bool {
	x.SetInterrupted(true)
	return true
}

func (i *Interrupt) String() string {
	return fmt.Sprintf("interrupt 0x%04x", uint16(i.Value))
}
