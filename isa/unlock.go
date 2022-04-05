package isa

import (
	"aletheiaware.com/flamego"
)

type Unlock struct {
	success bool
}

func NewUnlock() *Unlock {
	return &Unlock{}
}

func (i *Unlock) Load(x flamego.Context) (uint64, uint64, uint64, uint64) {
	i.success = true
	// Do Nothing
	return 0, 0, 0, 0
}

func (i *Unlock) Execute(x flamego.Context, a, b, c, d uint64) (uint64, uint64) {
	if !x.IsInterrupted() {
		// Hardware Lock only releasable in an interrupt
		x.Error(flamego.InterruptUnsupportedOperationError)
		i.success = false
		return 0, 0
	}
	x.SetRequiresLock(false)
	return 0, 0
}

func (i *Unlock) Format(x flamego.Context, a, b uint64) (uint64, uint64) {
	// Do Nothing
	return 0, 0
}

func (i *Unlock) Store(x flamego.Context, a, b uint64) {
	// Do Nothing
}

func (i *Unlock) Retire(x flamego.Context) bool {
	if i.success && !x.AcquiredLock() {
		x.IncrementProgramCounter()
		return true
	}
	return false
}

func (i *Unlock) String() string {
	return "unlock"
}
