package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type Multiply struct {
	Source1Register     flamego.Register
	Source2Register     flamego.Register
	DestinationRegister flamego.Register
}

func NewMultiply(s1, s2, d flamego.Register) *Multiply {
	return &Multiply{
		Source1Register:     s1,
		Source2Register:     s2,
		DestinationRegister: d,
	}
}

func (i *Multiply) Load(x flamego.Context) (uint64, uint64, uint64, uint64) {
	// Load Source 1 Register
	a := x.ReadRegister(i.Source1Register)
	// Load Source 2 Register
	b := x.ReadRegister(i.Source2Register)
	return a, b, 0, 0
}

func (i *Multiply) Execute(x flamego.Context, a, b, c, d uint64) (uint64, uint64) {
	return a * b, 0
}

func (i *Multiply) Format(x flamego.Context, a, b uint64) (uint64, uint64) {
	return a, 0
}

func (i *Multiply) Store(x flamego.Context, a, b uint64) {
	// Write Destination Register
	x.WriteRegister(i.DestinationRegister, a)
}

func (i *Multiply) Retire(x flamego.Context) bool {
	x.IncrementProgramCounter()
	return true
}

func (i *Multiply) String() string {
	return fmt.Sprintf("multiply %s %s %s", i.Source1Register, i.Source2Register, i.DestinationRegister)
}
