package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type Uninterrupt struct {
	AddressRegister flamego.Register
}

func NewUninterrupt(r flamego.Register) *Uninterrupt {
	return &Uninterrupt{
		AddressRegister: r,
	}
}

func (i *Uninterrupt) Load(x flamego.Context) (uint64, uint64, uint64) {
	// Load Return Address
	return x.ReadRegister(i.AddressRegister), 0, 0
}

func (i *Uninterrupt) Execute(x flamego.Context, a, b, c uint64) uint64 {
	// Do Nothing
	return a
}

func (i *Uninterrupt) Format(x flamego.Context, a uint64) uint64 {
	// Do Nothing
	return a
}

func (i *Uninterrupt) Store(x flamego.Context, a uint64) {
	// Jump out of interrupt by updating the program counter
	x.SetProgramCounter(a)
}

func (i *Uninterrupt) Retire(x flamego.Context) bool {
	x.SetInterrupted(false)
	return true
}

func (i *Uninterrupt) String() string {
	return fmt.Sprintf("uninterrupt %s", i.AddressRegister)
}
