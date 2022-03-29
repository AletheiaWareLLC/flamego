package isa

import (
	"aletheiaware.com/flamego"
)

type Unlock struct {
}

func NewUnlock() *Unlock {
	return &Unlock{}
}

func (i *Unlock) Load(x flamego.Context) (uint64, uint64, uint64) {
	// Do Nothing
	return 0, 0, 0
}

func (i *Unlock) Execute(x flamego.Context, a, b, c uint64) uint64 {
	x.SetRequiresLock(false)
	return 0
}

func (i *Unlock) Format(x flamego.Context, a uint64) uint64 {
	// Do Nothing
	return 0
}

func (i *Unlock) Store(x flamego.Context, a uint64) {
	// Do Nothing
}

func (i *Unlock) Retire(x flamego.Context) bool {
	if !x.AcquiredLock() {
		x.IncrementProgramCounter()
		return true
	}
	return false
}

func (i *Unlock) String() string {
	return "unlock"
}
