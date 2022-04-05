package isa

import (
	"aletheiaware.com/flamego"
)

type Halt struct {
	success bool
}

func NewHalt() *Halt {
	return &Halt{}
}

func (i *Halt) Load(x flamego.Context) (uint64, uint64, uint64, uint64) {
	i.success = true
	// Do Nothing
	return 0, 0, 0, 0
}

func (i *Halt) Execute(x flamego.Context, a, b, c, d uint64) (uint64, uint64) {
	if !x.IsInterrupted() {
		// Halt only allowed in an interrupt
		x.Error(flamego.InterruptUnsupportedOperationError)
		i.success = false
		return 0, 0
	}
	x.Core().Processor().Halt()
	return 0, 0
}

func (i *Halt) Format(x flamego.Context, a, b uint64) (uint64, uint64) {
	// Do Nothing
	return 0, 0
}

func (i *Halt) Store(x flamego.Context, a, b uint64) {
	// Do Nothing
}

func (i *Halt) Retire(x flamego.Context) bool {
	// Do Nothing
	return true
}

func (i *Halt) String() string {
	return "halt"
}
