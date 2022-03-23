package isa

import (
	"aletheiaware.com/flamego"
)

type Halt struct {
}

func NewHalt() *Halt {
	return &Halt{}
}

func (i *Halt) Load(x flamego.Context) (uint64, uint64, uint64) {
	// Do Nothing
	return 0, 0, 0
}

func (i *Halt) Execute(x flamego.Context, a, b, c uint64) uint64 {
	x.Core().Processor().Halt()
	return 0
}

func (i *Halt) Format(x flamego.Context, a uint64) uint64 {
	// Do Nothing
	return 0
}

func (i *Halt) Store(x flamego.Context, a uint64) {
	// Do Nothing
}

func (i *Halt) Retire(x flamego.Context) {
	// Do Nothing
}

func (i *Halt) String() string {
	return "halt"
}
