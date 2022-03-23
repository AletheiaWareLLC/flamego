package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type Clear struct {
	AddressRegister flamego.Register
	Offset          uint32
	success         bool
}

func NewClear(a flamego.Register, o uint32) *Clear {
	return &Clear{
		AddressRegister: a,
		Offset:          o,
	}
}

func (i *Clear) Load(x flamego.Context) (uint64, uint64, uint64) {
	i.success = true
	// Load Base Register
	a := x.ReadRegister(i.AddressRegister)
	// Load Offset
	b := uint64(i.Offset)
	return a, b, 0
}

func (i *Clear) Execute(x flamego.Context, a, b, c uint64) uint64 {
	d := x.Core().DataCache()
	if d.IsBusy() {
		i.success = false // Cache Unavailable
		return 0
	}
	// Issue Clear Request
	d.Clear(a + b)
	return 0
}

func (i *Clear) Format(x flamego.Context, a uint64) uint64 {
	if !i.success {
		return 0
	}
	d := x.Core().DataCache()
	if d.IsBusy() || !d.IsSuccessful() {
		i.success = false // Cache Miss
		return 0
	}
	return 0
}

func (i *Clear) Store(x flamego.Context, a uint64) {
	// Do Nothing
}

func (i *Clear) Retire(x flamego.Context) {
	if !i.success {
		return
	}
	// Only proceed if successfull
	x.IncrementProgramCounter()
}

func (i *Clear) String() string {
	return fmt.Sprintf("clear %s 0x%x", i.AddressRegister, i.Offset)
}
