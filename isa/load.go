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
	issued              bool
}

func NewLoad(a flamego.Register, o uint32, r flamego.Register) *Load {
	return &Load{
		AddressRegister:     a,
		Offset:              o,
		DestinationRegister: r,
	}
}

func (i *Load) Load(x flamego.Context) (uint64, uint64, uint64, uint64) {
	i.success = true
	// Load Base Register
	a := x.ReadRegister(i.AddressRegister)
	// Load Address
	b := uint64(i.Offset)
	return a, b, 0, 0
}

func (i *Load) Execute(x flamego.Context, a, b, c, d uint64) (uint64, uint64) {
	if !i.issued {
		l1d := x.Core().DataCache()
		if l1d.IsBusy() || !l1d.IsFree() {
			i.success = false // Cache Unavailable
			return 0, 0
		}
		// Issue Read Request
		l1d.Read(a + b)
		i.issued = true
	}
	return 0, 0
}

func (i *Load) Format(x flamego.Context, a, b uint64) (uint64, uint64) {
	if !i.success {
		return 0, 0
	}
	l1d := x.Core().DataCache()
	if l1d.IsBusy() {
		i.success = false
	} else if !l1d.IsSuccessful() {
		i.success = false
		i.issued = false // Reissue Request
		l1d.Free()       // Free Cache
	} else {
		// Copy Data from Bus
		buffer := make([]byte, 8)
		for i := 0; i < 8; i++ {
			buffer[i] = l1d.Bus().Read(i)
		}
		l1d.Free() // Free Cache
		return binary.BigEndian.Uint64(buffer), 0
	}
	return 0, 0
}

func (i *Load) Store(x flamego.Context, a, b uint64) {
	if !i.success {
		return
	}
	// Write Destination Register
	x.WriteRegister(i.DestinationRegister, a)
}

func (i *Load) Retire(x flamego.Context) bool {
	if i.success {
		x.IncrementProgramCounter()
		return true
	}
	return false
}

func (i *Load) String() string {
	return fmt.Sprintf("load %s 0x%x %s", i.AddressRegister, i.Offset, i.DestinationRegister)
}
