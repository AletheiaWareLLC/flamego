package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type Call struct {
	AddressRegister flamego.Register
	success         bool
}

func NewCall(a flamego.Register) *Call {
	return &Call{
		AddressRegister: a,
	}
	// TODO does the current PC get put onto the stack?
}

func (i *Call) Load(x flamego.Context) (uint64, uint64, uint64) {
	i.success = true
	// Load Address Register
	a := x.ReadRegister(i.AddressRegister)
	// Load Program Start Register
	b := x.ReadRegister(flamego.RProgramStart)
	// Load Program Limit Register
	c := x.ReadRegister(flamego.RProgramLimit)
	return a, b, c
}

func (i *Call) Execute(x flamego.Context, a, b, c uint64) uint64 {
	if !x.IsInterrupted() {
		// Only add program start if not in an interrupt
		a += b
		if a > c {
			x.Error(flamego.InterruptProgramAccessError)
			i.success = false
			return 0
		}
	}
	return a
}

func (i *Call) Format(x flamego.Context, a uint64) uint64 {
	// Pass Through
	return a
}

func (i *Call) Store(x flamego.Context, a uint64) {
	if i.success {
		// Update Program Counter
		x.SetProgramCounter(a)
	}
}

func (i *Call) Retire(x flamego.Context) bool {
	// Do Nothing
	return true
}

func (i *Call) String() string {
	return fmt.Sprintf("call %s", i.AddressRegister)
}
