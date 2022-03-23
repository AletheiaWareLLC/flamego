package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type LeftShift struct {
	Source1Register     flamego.Register
	Source2Register     flamego.Register
	DestinationRegister flamego.Register
}

func NewLeftShift(s1, s2, d flamego.Register) *LeftShift {
	return &LeftShift{
		Source1Register:     s1,
		Source2Register:     s2,
		DestinationRegister: d,
	}
}

func (i *LeftShift) Load(x flamego.Context) (uint64, uint64, uint64) {
	// Load Source 1 Register
	a := x.ReadRegister(i.Source1Register)
	// Load Source 2 Register
	b := x.ReadRegister(i.Source2Register)
	return a, b, 0
}

func (i *LeftShift) Execute(x flamego.Context, a, b, c uint64) uint64 {
	return a << b
}

func (i *LeftShift) Format(x flamego.Context, a uint64) uint64 {
	return a
}

func (i *LeftShift) Store(x flamego.Context, a uint64) {
	// Write Destination Register
	x.WriteRegister(i.DestinationRegister, a)
}

func (i *LeftShift) Retire(x flamego.Context) {
	x.IncrementProgramCounter()
}

func (i *LeftShift) String() string {
	return fmt.Sprintf("lshift %s %s %s", i.Source1Register, i.Source2Register, i.DestinationRegister)
}
