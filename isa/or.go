package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type Or struct {
	Source1Register     flamego.Register
	Source2Register     flamego.Register
	DestinationRegister flamego.Register
}

func NewOr(s1, s2, d flamego.Register) *Or {
	return &Or{
		Source1Register:     s1,
		Source2Register:     s2,
		DestinationRegister: d,
	}
}

func (i *Or) Load(x flamego.Context) (uint64, uint64, uint64, uint64) {
	// Load Source 1 Register
	a := x.ReadRegister(i.Source1Register)
	// Load Source 2 Register
	b := x.ReadRegister(i.Source2Register)
	return a, b, 0, 0
}

func (i *Or) Execute(x flamego.Context, a, b, c, d uint64) (uint64, uint64) {
	return a | b, 0
}

func (i *Or) Format(x flamego.Context, a, b uint64) (uint64, uint64) {
	return a, 0
}

func (i *Or) Store(x flamego.Context, a, b uint64) {
	// Write Destination Register
	x.WriteRegister(i.DestinationRegister, a)
}

func (i *Or) Retire(x flamego.Context) bool {
	x.IncrementProgramCounter()
	return true
}

func (i *Or) String() string {
	return fmt.Sprintf("or %s %s %s", i.Source1Register, i.Source2Register, i.DestinationRegister)
}
