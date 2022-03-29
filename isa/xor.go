package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type Xor struct {
	Source1Register     flamego.Register
	Source2Register     flamego.Register
	DestinationRegister flamego.Register
}

func NewXor(s1, s2, d flamego.Register) *Xor {
	return &Xor{
		Source1Register:     s1,
		Source2Register:     s2,
		DestinationRegister: d,
	}
}

func (i *Xor) Load(x flamego.Context) (uint64, uint64, uint64) {
	// Load Source 1 Register
	a := x.ReadRegister(i.Source1Register)
	// Load Source 2 Register
	b := x.ReadRegister(i.Source2Register)
	return a, b, 0
}

func (i *Xor) Execute(x flamego.Context, a, b, c uint64) uint64 {
	return a ^ b
}

func (i *Xor) Format(x flamego.Context, a uint64) uint64 {
	return a
}

func (i *Xor) Store(x flamego.Context, a uint64) {
	// Write Destination Register
	x.WriteRegister(i.DestinationRegister, a)
}

func (i *Xor) Retire(x flamego.Context) bool {
	x.IncrementProgramCounter()
	return true
}

func (i *Xor) String() string {
	return fmt.Sprintf("xor %s %s %s", i.Source1Register, i.Source2Register, i.DestinationRegister)
}
