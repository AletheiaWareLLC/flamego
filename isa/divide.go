package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type Divide struct {
	Source1Register     flamego.Register
	Source2Register     flamego.Register
	DestinationRegister flamego.Register
}

func NewDivide(s1, s2, d flamego.Register) *Divide {
	return &Divide{
		Source1Register:     s1,
		Source2Register:     s2,
		DestinationRegister: d,
	}
}

func (i *Divide) Load(x flamego.Context) (uint64, uint64, uint64) {
	// Load Source 1 Register
	a := x.ReadRegister(i.Source1Register)
	// Load Source 2 Register
	b := x.ReadRegister(i.Source2Register)
	return a, b, 0
}

func (i *Divide) Execute(x flamego.Context, a, b, c uint64) uint64 {
	if b == 0 {
		x.Error(flamego.InterruptArithmeticError)
		return 0
	}
	return a / b
}

func (i *Divide) Format(x flamego.Context, a uint64) uint64 {
	return a
}

func (i *Divide) Store(x flamego.Context, a uint64) {
	// Write Destination Register
	x.WriteRegister(i.DestinationRegister, a)
}

func (i *Divide) Retire(x flamego.Context) bool {
	x.IncrementProgramCounter()
	return true
}

func (i *Divide) String() string {
	return fmt.Sprintf("divide %s %s %s", i.Source1Register, i.Source2Register, i.DestinationRegister)
}
