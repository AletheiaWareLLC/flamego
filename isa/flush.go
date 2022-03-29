package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
	"log"
)

type Flush struct {
	AddressRegister flamego.Register
	Offset          uint32
	success         bool
	issuedL1D       bool
	issuedL2        bool
	flushedL1D      bool
	flushedL2       bool
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
	address := a + b
	if !i.issuedL1D {
		l1d := x.Core().DataCache()
		if l1d.IsBusy() || !l1d.IsFree() {
			log.Println("Execute: L1D Cache Unavailable")
			i.success = false // Cache Unavailable
			return 0
		}

		// Issue Flush Request
		l1d.Flush(address)
		i.issuedL1D = true
		log.Println("Execute: L1D Flush Issued")
	} else if !i.issuedL2 {
		l2 := x.Core().Processor().Cache()
		if l2.IsBusy() || !l2.IsFree() {
			log.Println("Execute: L2 Cache Unavailable")
			i.success = false // Cache Unavailable
			return 0
		}

		// Issue Flush Request
		l2.Flush(address)
		i.issuedL2 = true
		log.Println("Execute: L2 Flush Issued")
	}
	return 0
}

func (i *Flush) Format(x flamego.Context, a uint64) uint64 {
	if !i.success {
		return 0
	}
	if !i.flushedL1D {
		l1d := x.Core().DataCache()
		log.Println("L1D", l1d.IsBusy(), l1d.IsSuccessful())
		if l1d.IsBusy() {
			log.Println("Format: L1D Cache Busy")
			i.success = false
		} else if !l1d.IsSuccessful() {
			log.Println("Format: L1D Cache Miss")
			i.success = false
			i.issuedL1D = false // Reissue Request
			l1d.Free()          // Free Cache
		} else {
			i.flushedL1D = true
			l1d.Free() // Free Cache
			log.Println("Format: L1D Flush Successful")
		}
	} else if !i.flushedL2 {
		l2 := x.Core().Processor().Cache()
		log.Println("L2", l2.IsBusy(), l2.IsSuccessful())
		if l2.IsBusy() {
			log.Println("Format: L2 Cache Busy")
			i.success = false
		} else if !l2.IsSuccessful() {
			log.Println("Format: L2 Cache Miss")
			i.success = false
			i.issuedL2 = false // Reissue Request
			l2.Free()          // Free Cache
		} else {
			i.flushedL2 = true
			l2.Free() // Free Cache
			log.Println("Format: L2 Flush Successful")
		}
	}
	return 0
}

func (i *Flush) Store(x flamego.Context, a uint64) {
	// Do Nothing
}

func (i *Flush) Retire(x flamego.Context) bool {
	if i.success && i.flushedL1D && i.flushedL2 {
		x.IncrementProgramCounter()
		return true
	}
	return false
}

func (i *Flush) String() string {
	return fmt.Sprintf("flush %s 0x%x", i.AddressRegister, i.Offset)
}
