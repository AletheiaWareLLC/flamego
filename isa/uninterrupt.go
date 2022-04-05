package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type Uninterrupt struct {
	AddressRegister flamego.Register
	success         bool
}

func NewUninterrupt(r flamego.Register) *Uninterrupt {
	return &Uninterrupt{
		AddressRegister: r,
	}
}

func (i *Uninterrupt) Load(x flamego.Context) (uint64, uint64, uint64, uint64) {
	i.success = true
	// Load Return Address
	return x.ReadRegister(i.AddressRegister), 0, 0, 0
}

func (i *Uninterrupt) Execute(x flamego.Context, a, b, c, d uint64) (uint64, uint64) {
	if !x.IsInterrupted() {
		// Uinterrupt only allowed in an interrupt
		x.Error(flamego.InterruptUnsupportedOperationError)
		i.success = false
		return 0, 0
	}
	return a, 0
}

func (i *Uninterrupt) Format(x flamego.Context, a, b uint64) (uint64, uint64) {
	// Do Nothing
	return a, 0
}

func (i *Uninterrupt) Store(x flamego.Context, a, b uint64) {
	if i.success {
		// Jump out of interrupt by updating the program counter
		x.SetProgramCounter(a)
	}
}

func (i *Uninterrupt) Retire(x flamego.Context) bool {
	if i.success {
		x.SetInterrupted(false)
	}
	return true
}

func (i *Uninterrupt) String() string {
	return fmt.Sprintf("uninterrupt %s", i.AddressRegister)
}
