package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type Subtract struct {
	Source1Register     flamego.Register
	Source2Register     flamego.Register
	DestinationRegister flamego.Register
}

func NewSubtract(s1, s2, d flamego.Register) *Subtract {
	return &Subtract{
		Source1Register:     s1,
		Source2Register:     s2,
		DestinationRegister: d,
	}
}

func (i *Subtract) Load(x flamego.Context) (uint64, uint64, uint64) {
	// Load Source 1 Register
	a := x.ReadRegister(i.Source1Register)
	// Load Source 2 Register
	b := x.ReadRegister(i.Source2Register)
	return a, b, 0
}

func (i *Subtract) Execute(x flamego.Context, a, b, c uint64) uint64 {
	return a - b
}

func (i *Subtract) Format(x flamego.Context, a uint64) uint64 {
	return a
}

func (i *Subtract) Store(x flamego.Context, a uint64) {
	// Write Destination Register
	x.WriteRegister(i.DestinationRegister, a)
}

func (i *Subtract) Retire(x flamego.Context) {
	x.IncrementProgramCounter()
}

func (i *Subtract) String() string {
	return fmt.Sprintf("subtract %s %s %s", i.Source1Register, i.Source2Register, i.DestinationRegister)
}
