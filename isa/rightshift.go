package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type RightShift struct {
	Source1Register     flamego.Register
	Source2Register     flamego.Register
	DestinationRegister flamego.Register
}

func NewRightShift(s1, s2, d flamego.Register) *RightShift {
	return &RightShift{
		Source1Register:     s1,
		Source2Register:     s2,
		DestinationRegister: d,
	}
}

func (i *RightShift) Load(x flamego.Context) (uint64, uint64, uint64) {
	// Load Source 1 Register
	a := x.ReadRegister(i.Source1Register)
	// Load Source 2 Register
	b := x.ReadRegister(i.Source2Register)
	return a, b, 0
}

func (i *RightShift) Execute(x flamego.Context, a, b, c uint64) uint64 {
	return a >> b
}

func (i *RightShift) Format(x flamego.Context, a uint64) uint64 {
	return a
}

func (i *RightShift) Store(x flamego.Context, a uint64) {
	// Write Destination Register
	x.WriteRegister(i.DestinationRegister, a)
}

func (i *RightShift) Retire(x flamego.Context) {
	x.IncrementProgramCounter()
}

func (i *RightShift) String() string {
	return fmt.Sprintf("rshift %s %s %s", i.Source1Register, i.Source2Register, i.DestinationRegister)
}
