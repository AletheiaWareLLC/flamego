package isa

import (
	"aletheiaware.com/flamego"
)

type Noop struct {
}

func NewNoop() *Noop {
	return &Noop{}
}

func (i *Noop) Load(x flamego.Context) (uint64, uint64, uint64, uint64) {
	// Do Nothing
	return 0, 0, 0, 0
}

func (i *Noop) Execute(x flamego.Context, a, b, c, d uint64) (uint64, uint64) {
	// Do Nothing
	return 0, 0
}

func (i *Noop) Format(x flamego.Context, a, b uint64) (uint64, uint64) {
	// Do Nothing
	return 0, 0
}

func (i *Noop) Store(x flamego.Context, a, b uint64) {
	// Do Nothing
}

func (i *Noop) Retire(x flamego.Context) bool {
	x.IncrementProgramCounter()
	return true
}

func (i *Noop) String() string {
	return "noop"
}
