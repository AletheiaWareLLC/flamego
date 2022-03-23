package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type Pop struct {
	DestinationRegister flamego.Register
	success             bool
}

func NewPop(a flamego.Register) *Pop {
	return &Pop{
		DestinationRegister: a,
	}
}

func (i *Pop) Load(x flamego.Context) (uint64, uint64, uint64) {
	i.success = true
	// Load Stack Pointer Register
	a := x.ReadRegister(flamego.RStackPointer)
	// Load Stack Start Register
	b := x.ReadRegister(flamego.RStackStart)
	return a, b, 0
}

func (i *Pop) Execute(x flamego.Context, a, b, c uint64) uint64 {
	a -= flamego.DataSize
	if a < b {
		x.Error(flamego.InterruptStackUnderflowError)
		i.success = false
	} else {
		// TODO write c to memory[a]
	}
	return a
}

func (i *Pop) Format(x flamego.Context, a uint64) uint64 {
	// Pass Through
	return a
}

func (i *Pop) Store(x flamego.Context, a uint64) {
	if i.success {
		// Increment Stack Pointer
		x.WriteRegister(flamego.RStackPointer, a)
	}
}

func (i *Pop) Retire(x flamego.Context) {
	if i.success {
		x.IncrementProgramCounter()
	}
}

func (i *Pop) String() string {
	return fmt.Sprintf("pop %s", i.DestinationRegister)
}
