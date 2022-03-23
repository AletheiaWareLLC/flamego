package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type Flush struct {
	AddressRegister flamego.Register
	Offset          uint32
	success         bool
}

func NewFlush(a flamego.Register, o uint32) *Flush {
	return &Flush{
		AddressRegister: a,
		Offset:          o,
	}
}

func (i *Flush) Load(x flamego.Context) (uint64, uint64, uint64) {
	i.success = true
	// Load Base Register
	a := x.ReadRegister(i.AddressRegister)
	// Load Offset
	b := uint64(i.Offset)
	return a, b, 0
}

func (i *Flush) Execute(x flamego.Context, a, b, c uint64) uint64 {
	d := x.Core().DataCache()
	if d.IsBusy() {
		i.success = false // Cache Unavailable
		return 0
	}
	// Issue Flush Request
	d.Flush(a + b)
	return 0
}

func (i *Flush) Format(x flamego.Context, a uint64) uint64 {
	d := x.Core().DataCache()
	if d.IsBusy() || !d.IsSuccessful() {
		i.success = false // Cache Miss
		return 0
	}
	return 0
}

func (i *Flush) Store(x flamego.Context, a uint64) {
	// Do Nothing
}

func (i *Flush) Retire(x flamego.Context) {
	if !i.success {
		return
	}
	// Only proceed if successfull
	x.IncrementProgramCounter()
}

func (i *Flush) String() string {
	return fmt.Sprintf("flush %s 0x%x", i.AddressRegister, i.Offset)
}
