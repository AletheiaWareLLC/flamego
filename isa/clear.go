package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type Clear struct {
	AddressRegister flamego.Register
	Offset          uint32
	success         bool
	issuedL1I       bool
	issuedL1D       bool
	issuedL2        bool
	clearedL1I      bool
	clearedL1D      bool
	clearedL2       bool
}

func NewClear(a flamego.Register, o uint32) *Clear {
	return &Clear{
		AddressRegister: a,
		Offset:          o,
	}
}

func (i *Clear) Load(x flamego.Context) (uint64, uint64, uint64, uint64) {
	i.success = true
	// Load Base Register
	a := x.ReadRegister(i.AddressRegister)
	// Load Offset
	b := uint64(i.Offset)
	return a, b, 0, 0
}

func (i *Clear) Execute(x flamego.Context, a, b, c, d uint64) (uint64, uint64) {
	address := a + b
	if !i.issuedL1I {
		l1i := x.Core().InstructionCache()
		if l1i.IsBusy() || !l1i.IsFree() {
			i.success = false // Cache Unavailable
			return 0, 0
		}

		// Issue Clear Request
		l1i.Clear(address)
		i.issuedL1I = true
	} else if !i.issuedL1D {
		l1d := x.Core().DataCache()
		if l1d.IsBusy() || !l1d.IsFree() {
			i.success = false // Cache Unavailable
			return 0, 0
		}

		// Issue Clear Request
		l1d.Clear(address)
		i.issuedL1D = true
	} else if !i.issuedL2 {
		l2 := x.Core().Processor().Cache()
		if l2.IsBusy() || !l2.IsFree() {
			i.success = false // Cache Unavailable
			return 0, 0
		}

		// Issue Clear Request
		l2.Clear(address)
		i.issuedL2 = true
	}
	return 0, 0
}

func (i *Clear) Format(x flamego.Context, a, b uint64) (uint64, uint64) {
	if !i.success {
		return 0, 0
	}
	if !i.clearedL1I {
		l1i := x.Core().InstructionCache()
		if l1i.IsBusy() {
			i.success = false
		} else if !l1i.IsSuccessful() {
			i.success = false
			i.issuedL1I = false // Reissue Request
			l1i.Free()          // Free Cache
		} else {
			i.clearedL1I = true
			l1i.Free() // Free Cache
		}
	} else if !i.clearedL1D {
		l1d := x.Core().DataCache()
		if l1d.IsBusy() {
			i.success = false
		} else if !l1d.IsSuccessful() {
			i.success = false
			i.issuedL1D = false // Reissue Request
			l1d.Free()          // Free Cache
		} else {
			i.clearedL1D = true
			l1d.Free() // Free Cache
		}
	} else if !i.clearedL2 {
		l2 := x.Core().Processor().Cache()
		if l2.IsBusy() {
			i.success = false
		} else if !l2.IsSuccessful() {
			i.success = false
			i.issuedL2 = false // Reissue Request
			l2.Free()          // Free Cache
		} else {
			i.clearedL2 = true
			l2.Free() // Free Cache
		}
	}
	return 0, 0
}

func (i *Clear) Store(x flamego.Context, a, b uint64) {
	// Do Nothing
}

func (i *Clear) Retire(x flamego.Context) bool {
	if i.success && i.clearedL1I && i.clearedL1D && i.clearedL2 {
		x.IncrementProgramCounter()
		return true
	}
	return false
}

func (i *Clear) String() string {
	return fmt.Sprintf("clear %s 0x%x", i.AddressRegister, i.Offset)
}
