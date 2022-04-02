package isa

import (
	"aletheiaware.com/flamego"
)

type Lock struct {
	success bool
}

func NewLock() *Lock {
	return &Lock{}
}

func (i *Lock) Load(x flamego.Context) (uint64, uint64, uint64) {
	i.success = true
	// Do Nothing
	return 0, 0, 0
}

func (i *Lock) Execute(x flamego.Context, a, b, c uint64) uint64 {
	if !x.IsInterrupted() {
		// Hardware Lock acquirable only in an interrupt
		x.Error(flamego.InterruptUnsupportedOperationError)
		i.success = false
		return 0
	}
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

func (i *Lock) Retire(x flamego.Context) bool {
	if i.success && x.AcquiredLock() {
		x.IncrementProgramCounter()
		return true
	}
	return false
}

func (i *Lock) String() string {
	return "lock"
}
