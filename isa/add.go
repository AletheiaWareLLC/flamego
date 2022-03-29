package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type Add struct {
	Source1Register     flamego.Register
	Source2Register     flamego.Register
	DestinationRegister flamego.Register
}

func NewAdd(s1, s2, d flamego.Register) *Add {
	return &Add{
		Source1Register:     s1,
		Source2Register:     s2,
		DestinationRegister: d,
	}
}

func (i *Add) Load(x flamego.Context) (uint64, uint64, uint64) {
	// Load Source 1 Register
	a := x.ReadRegister(i.Source1Register)
	// Load Source 2 Register
	b := x.ReadRegister(i.Source2Register)
	return a, b, 0
}

func (i *Add) Execute(x flamego.Context, a, b, c uint64) uint64 {
	return a + b
}

func (i *Add) Format(x flamego.Context, a uint64) uint64 {
	return a
}

func (i *Add) Store(x flamego.Context, a uint64) {
	// Write Destination Register
	x.WriteRegister(i.DestinationRegister, a)
}

func (i *Add) Retire(x flamego.Context) bool {
	x.IncrementProgramCounter()
	return true
}

func (i *Add) String() string {
	return fmt.Sprintf("add %s %s %s", i.Source1Register, i.Source2Register, i.DestinationRegister)
}
