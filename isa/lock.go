package isa

import (
	"aletheiaware.com/flamego"
)

type Lock struct {
}

func NewLock() *Lock {
	return &Lock{}
}

func (i *Lock) Load(x flamego.Context) (uint64, uint64, uint64) {
	// Do Nothing
	return 0, 0, 0
}

func (i *Lock) Execute(x flamego.Context, a, b, c uint64) uint64 {
	x.SetRequiresLock(true)
	return 0
}

func (i *Lock) Format(x flamego.Context, a uint64) uint64 {
	// Do Nothing
	return 0
}

func (i *Lock) Store(x flamego.Context, a uint64) {
	// Do Nothing
}

func (i *Lock) Retire(x flamego.Context) {
	if x.AcquiredLock() {
		x.IncrementProgramCounter()
	}
}

func (i *Lock) String() string {
	return "lock"
}
