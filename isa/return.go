package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type Return struct {
	AddressRegister flamego.Register
	success         bool
}

func NewReturn(a flamego.Register) *Return {
	return &Return{
		AddressRegister: a,
	}
}

func (i *Return) Load(x flamego.Context) (uint64, uint64, uint64) {
	i.success = true
	// Load Address Register
	a := x.ReadRegister(i.AddressRegister)
	// Load Program Start Register
	b := x.ReadRegister(flamego.RProgramStart)
	// Load Program Limit Register
	c := x.ReadRegister(flamego.RProgramLimit)
	return a, b, c
}

func (i *Return) Execute(x flamego.Context, a, b, c uint64) uint64 {
	if a > c {
		x.Error(flamego.InterruptProgramAccessError)
		i.success = false
		return 0
	}
	if !x.IsInterrupted() {
		// Only add program start if not in an interrupt
		a += b
	}
	return a
}

func (i *Return) Format(x flamego.Context, a uint64) uint64 {
	// Pass Through
	return a
}

func (i *Return) Store(x flamego.Context, a uint64) {
	if i.success {
		// Update Program Counter
		x.SetProgramCounter(a)
	}
}

func (i *Return) Retire(x flamego.Context) {
	// Do Nothing
}

func (i *Return) String() string {
	return fmt.Sprintf("return %s", i.AddressRegister)
}
