package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type Modulo struct {
	Source1Register     flamego.Register
	Source2Register     flamego.Register
	DestinationRegister flamego.Register
}

func NewModulo(s1, s2, d flamego.Register) *Modulo {
	return &Modulo{
		Source1Register:     s1,
		Source2Register:     s2,
		DestinationRegister: d,
	}
}

func (i *Modulo) Load(x flamego.Context) (uint64, uint64, uint64, uint64) {
	// Load Source 1 Register
	a := x.ReadRegister(i.Source1Register)
	// Load Source 2 Register
	b := x.ReadRegister(i.Source2Register)
	return a, b, 0, 0
}

func (i *Modulo) Execute(x flamego.Context, a, b, c, d uint64) (uint64, uint64) {
	if b == 0 {
		x.Error(flamego.InterruptArithmeticError)
		return 0, 0
	}
	return a % b, 0
}

func (i *Modulo) Format(x flamego.Context, a, b uint64) (uint64, uint64) {
	return a, 0
}

func (i *Modulo) Store(x flamego.Context, a, b uint64) {
	// Write Destination Register
	x.WriteRegister(i.DestinationRegister, a)
}

func (i *Modulo) Retire(x flamego.Context) bool {
	x.IncrementProgramCounter()
	return true
}

func (i *Modulo) String() string {
	return fmt.Sprintf("modulo %s %s %s", i.Source1Register, i.Source2Register, i.DestinationRegister)
}
