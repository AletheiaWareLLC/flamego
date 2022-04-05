package isa

import (
	"aletheiaware.com/flamego"
	"encoding/binary"
	"fmt"
)

type Pop struct {
	Mask    uint16
	success bool
	issued  bool
	index   uint8
}

func NewPop(m uint16) *Pop {
	return &Pop{
		Mask:  m,
		index: 0, // Start with r31 (LSB), interate down to r16 (MSB)
	}
}

func (i *Pop) Load(x flamego.Context) (uint64, uint64, uint64, uint64) {
	i.success = true
	// Load Stack Pointer Register
	a := x.ReadRegister(flamego.RStackPointer)
	// Load Stack Start Register
	b := x.ReadRegister(flamego.RStackStart)
	return a, b, 0, 0
}

func (i *Pop) Execute(x flamego.Context, a, b, c, d uint64) (uint64, uint64) {
	// Decrement Stack Pointer
	a -= flamego.DataSize

	if a < b {
		x.Error(flamego.InterruptStackUnderflowError)
		i.success = false
	} else if !i.issued {
		l1d := x.Core().DataCache()
		if l1d.IsBusy() || !l1d.IsFree() {
			i.success = false // Cache Unavailable
			return 0, 0
		}
		// Issue Read Request
		l1d.Read(a)
		i.issued = true
	}
	return a, 0
}

func (i *Pop) Format(x flamego.Context, a, b uint64) (uint64, uint64) {
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
		b := binary.BigEndian.Uint64(buffer)
		return a, b
	}
	return a, 0
}

func (i *Pop) Store(x flamego.Context, a, b uint64) {
	if !i.success {
		return
	}
	// Increment Stack Pointer
	x.WriteRegister(flamego.RStackPointer, a)
	for ; i.index < 16; i.index++ {
		m := (uint16(1) << i.index)
		if i.Mask&m != 0 {
			r := flamego.R31 - flamego.Register(i.index)
			// Save Popped Value
			x.WriteRegister(r, b)
			// Clear bit from mask
			i.Mask &= ^m
			break
		}
	}
}

func (i *Pop) Retire(x flamego.Context) bool {
	if i.success {
		if i.Mask == 0 {
			x.IncrementProgramCounter()
			return true
		} else {
			// Reset issued flag
			i.issued = false
		}
	}
	return false
}

func (i *Pop) String() string {
	return fmt.Sprintf("pop %016b", i.Mask)
}
