package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type Flush struct {
	AddressRegister flamego.Register
	Offset          uint32
	success         bool
	issuedL1D       bool
	issuedL2        bool
	issuedL3        bool
	flushedL1D      bool
	flushedL2       bool
	flushedL3       bool
}

func NewFlush(a flamego.Register, o uint32) *Flush {
	return &Flush{
		AddressRegister: a,
		Offset:          o,
	}
}

func (i *Flush) Load(x flamego.Context) (uint64, uint64, uint64, uint64) {
	i.success = true
	// Load Base Register
	a := x.ReadRegister(i.AddressRegister)
	// Load Offset
	b := uint64(i.Offset)
	return a, b, 0, 0
}

func (i *Flush) Execute(x flamego.Context, a, b, c, d uint64) (uint64, uint64) {
	address := a + b
	if !i.issuedL1D {
		l1d := x.DataCache()
		if l1d.IsBusy() || !l1d.IsFree() {
			i.success = false // Cache Unavailable
			return 0, 0
		}

		// Issue Flush Request
		l1d.Flush(address)
		i.issuedL1D = true
	} else if !i.issuedL2 && i.flushedL1D {
		l2 := x.Core().Cache()
		if l2.IsBusy() || !l2.IsFree() {
			i.success = false // Cache Unavailable
			return 0, 0
		}

		// Issue Flush Request
		l2.Flush(address)
		i.issuedL2 = true
	} else if !i.issuedL3 && i.flushedL2 {
		l3 := x.Core().Processor().Cache()
		if l3.IsBusy() || !l3.IsFree() {
			i.success = false // Cache Unavailable
			return 0, 0
		}

		// Issue Flush Request
		l3.Flush(address)
		i.issuedL3 = true
	}
	return 0, 0
}

func (i *Flush) Format(x flamego.Context, a, b uint64) (uint64, uint64) {
	if !i.success {
		return 0, 0
	}
	if !i.flushedL1D {
		l1d := x.DataCache()
		if l1d.IsBusy() {
			i.success = false
		} else if !l1d.IsSuccessful() {
			i.success = false
			i.issuedL1D = false // Reissue Request
			l1d.Free()          // Free Cache
		} else {
			i.flushedL1D = true
			l1d.Free() // Free Cache
		}
	} else if !i.flushedL2 {
		l2 := x.Core().Cache()
		if l2.IsBusy() {
			i.success = false
		} else if !l2.IsSuccessful() {
			i.success = false
			i.issuedL2 = false // Reissue Request
			l2.Free()          // Free Cache
		} else {
			i.flushedL2 = true
			l2.Free() // Free Cache
		}
	} else if !i.flushedL3 {
		l3 := x.Core().Processor().Cache()
		if l3.IsBusy() {
			i.success = false
		} else if !l3.IsSuccessful() {
			i.success = false
			i.issuedL3 = false // Reissue Request
			l3.Free()          // Free Cache
		} else {
			i.flushedL3 = true
			l3.Free() // Free Cache
		}
	}
	return 0, 0
}

func (i *Flush) Store(x flamego.Context, a, b uint64) {
	// Do Nothing
}

func (i *Flush) Retire(x flamego.Context) bool {
	if i.success && i.flushedL1D && i.flushedL2 && i.flushedL3 {
		x.IncrementProgramCounter()
		return true
	}
	return false
}

func (i *Flush) String() string {
	return fmt.Sprintf("flush %s 0x%x", i.AddressRegister, i.Offset)
}
