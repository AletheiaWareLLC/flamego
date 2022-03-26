package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type LoadC struct {
	Constant            uint32
	DestinationRegister flamego.Register
}

func NewLoadC(c uint32, r flamego.Register) *LoadC {
	return &LoadC{
		Constant:            c,
		DestinationRegister: r,
	}
}

func (i *LoadC) Load(x flamego.Context) (uint64, uint64, uint64) {
	// Load Constant
	return uint64(i.Constant), 0, 0
}

func (i *LoadC) Execute(x flamego.Context, a, b, c uint64) uint64 {
	// Pass Through
	return a
}

func (i *LoadC) Format(x flamego.Context, a uint64) uint64 {
	// Pass Through
	return a
}

func (i *LoadC) Store(x flamego.Context, a uint64) {
	// Write Destination Register
	x.WriteRegister(i.DestinationRegister, a)
}

func (i *LoadC) Retire(x flamego.Context) {
	x.IncrementProgramCounter()
}

func (i *LoadC) String() string {
	return fmt.Sprintf("loadc 0x%x %s", i.Constant, i.DestinationRegister)
}
