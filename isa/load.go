package isa

import (
	"aletheiaware.com/flamego"
	"encoding/binary"
	"fmt"
)

type Load struct {
	AddressRegister     flamego.Register
	Offset              uint32
	DestinationRegister flamego.Register
	success             bool
}

func NewLoad(a flamego.Register, o uint32, r flamego.Register) *Load {
	return &Load{
		AddressRegister:     a,
		Offset:              o,
		DestinationRegister: r,
	}
}

func (i *Load) Load(x flamego.Context) (uint64, uint64, uint64) {
	i.success = true
	// Load Base Register
	a := x.ReadRegister(i.AddressRegister)
	// Load Address
	b := uint64(i.Offset)
	return a, b, 0
}

func (i *Load) Execute(x flamego.Context, a, b, c uint64) uint64 {
	d := x.Core().DataCache()
	if d.IsBusy() {
		i.success = false // Cache Unavailable
		return 0
	}
	// Issue Read Request
	d.Read(a + b)
	return 0
}

func (i *Load) Format(x flamego.Context, a uint64) uint64 {
	if !i.success {
		return 0
	}
	d := x.Core().DataCache()
	if d.IsBusy() || !d.IsSuccessful() {
		i.success = false // Cache Miss
		return 0
	}
	// Copy Data from Bus
	buffer := make([]byte, 8)
	for i := 0; i < 8; i++ {
		buffer[i] = d.Bus().Read(i)
	}
	return binary.BigEndian.Uint64(buffer)
}

func (i *Load) Store(x flamego.Context, a uint64) {
	if !i.success {
		return
	}
	// Write Destination Register
	x.WriteRegister(i.DestinationRegister, a)
}

func (i *Load) Retire(x flamego.Context) {
	if !i.success {
		return
	}
	// Only proceed if successfull
	x.IncrementProgramCounter()
}

func (i *Load) String() string {
	return fmt.Sprintf("load %s 0x%x %s", i.AddressRegister, i.Offset, i.DestinationRegister)
}
