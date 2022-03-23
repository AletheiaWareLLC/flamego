package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type Not struct {
	SourceRegister      flamego.Register
	DestinationRegister flamego.Register
}

func NewNot(s, d flamego.Register) *Not {
	return &Not{
		SourceRegister:      s,
		DestinationRegister: d,
	}
}

func (i *Not) Load(x flamego.Context) (uint64, uint64, uint64) {
	// Load Source Register
	a := x.ReadRegister(i.SourceRegister)
	return a, 0, 0
}

func (i *Not) Execute(x flamego.Context, a, b, c uint64) uint64 {
	return ^a
}

func (i *Not) Format(x flamego.Context, a uint64) uint64 {
	return a
}

func (i *Not) Store(x flamego.Context, a uint64) {
	// Write Destination Register
	x.WriteRegister(i.DestinationRegister, a)
}

func (i *Not) Retire(x flamego.Context) {
	x.IncrementProgramCounter()
}

func (i *Not) String() string {
	return fmt.Sprintf("not %s %s", i.SourceRegister, i.DestinationRegister)
}
