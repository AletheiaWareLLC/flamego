package isa

import (
	"aletheiaware.com/flamego"
	"encoding/binary"
	"fmt"
	"log"
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

func (i *Load) Load(x flamego.Context) (uint64, uint64, uint64) {
	i.success = true
	// Load Base Register
	a := x.ReadRegister(i.AddressRegister)
	// Load Address
	b := uint64(i.Offset)
	return a, b, 0
}

func (i *Load) Execute(x flamego.Context, a, b, c uint64) uint64 {
	if !i.issued {
		l1d := x.Core().DataCache()
		if l1d.IsBusy() || !l1d.IsFree() {
			log.Println("Execute: L1D Cache Unavailable")
			i.success = false // Cache Unavailable
			return 0
		}
		// Issue Read Request
		l1d.Read(a + b)
		i.issued = true
		log.Println("Execute: L1D Load Issued")
	}
	return 0
}

func (i *Load) Format(x flamego.Context, a uint64) uint64 {
	if !i.success {
		return 0
	}
	l1d := x.Core().DataCache()
	if l1d.IsBusy() {
		log.Println("Format: L1D Cache Busy")
		i.success = false
	} else if !l1d.IsSuccessful() {
		log.Println("Format: L1D Cache Miss")
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
		log.Println("Format: L1D Load Successful")
		return binary.BigEndian.Uint64(buffer)
	}
	return 0
}

func (i *Load) Store(x flamego.Context, a uint64) {
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
