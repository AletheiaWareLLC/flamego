package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type Push struct {
	SourceRegister flamego.Register
	success        bool
}

func NewPush(a flamego.Register) *Push {
	return &Push{
		SourceRegister: a,
	}
}

func (i *Push) Load(x flamego.Context) (uint64, uint64, uint64) {
	i.success = true
	// Load Stack Pointer Register
	a := x.ReadRegister(flamego.RStackPointer)
	// Load Stack Limit Register
	b := x.ReadRegister(flamego.RStackLimit)
	// Load Source Register
	c := x.ReadRegister(i.SourceRegister)
	return a, b, c
}

func (i *Push) Execute(x flamego.Context, a, b, c uint64) uint64 {
	if a > b {
		x.Error(flamego.InterruptStackOverflowError)
		i.success = false
	} else {
		// TODO write c to memory[a]
	}
	return a
}

func (i *Push) Format(x flamego.Context, a uint64) uint64 {
	// Pass Through
	return a
}

func (i *Push) Store(x flamego.Context, a uint64) {
	if i.success {
		// Increment Stack Pointer
		x.WriteRegister(flamego.RStackPointer, a+flamego.DataSize)
	}
}

func (i *Push) Retire(x flamego.Context) bool {
	if i.success {
		x.IncrementProgramCounter()
		return true
	}
	return false
}

func (i *Push) String() string {
	return fmt.Sprintf("push %s", i.SourceRegister)
}
