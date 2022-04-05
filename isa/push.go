package isa

import (
	"aletheiaware.com/flamego"
	"encoding/binary"
	"fmt"
)

type Push struct {
	Mask    uint16
	success bool
	issued  bool
	index   uint8
}

func NewPush(m uint16) *Push {
	return &Push{
		Mask:  m,  // 16 bits representing r16..r31
		index: 15, // Start with r16 (MSB), interate up to r31 (LSB)
	}
}

func (i *Push) Load(x flamego.Context) (uint64, uint64, uint64, uint64) {
	i.success = true
	// Load Stack Pointer Register
	a := x.ReadRegister(flamego.RStackPointer)
	// Load Stack Limit Register
	b := x.ReadRegister(flamego.RStackLimit)
	var c uint64
	for ; i.index >= 0; i.index-- {
		m := (uint16(1) << i.index)
		if i.Mask&m != 0 {
			r := flamego.R31 - flamego.Register(i.index)
			// Load Source Register
			c = x.ReadRegister(r)
			// Clear bit from mask
			i.Mask &= ^m
			break
		}
	}
	return a, b, c, 0
}

func (i *Push) Execute(x flamego.Context, a, b, c, d uint64) (uint64, uint64) {
	if a >= b {
		x.Error(flamego.InterruptStackOverflowError)
		i.success = false
	} else if !i.issued {
		l1d := x.Core().DataCache()
		if l1d.IsBusy() || !l1d.IsFree() {
			i.success = false // Cache Unavailable
			return 0, 0
		}
		// Copy Data to Bus
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, c)
		for i := 0; i < 8; i++ {
			l1d.Bus().Write(i, buffer[i])
		}
		// Issue Write Request
		l1d.Write(a)
		i.issued = true
	}
	return a, 0
}

func (i *Push) Format(x flamego.Context, a, b uint64) (uint64, uint64) {
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
		l1d.Free() // Free Cache
	}
	return a, 0
}

func (i *Push) Store(x flamego.Context, a, b uint64) {
	if i.success {
		// Increment Stack Pointer
		p := a + flamego.DataSize
		x.WriteRegister(flamego.RStackPointer, p)
	}
}

func (i *Push) Retire(x flamego.Context) bool {
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

func (i *Push) String() string {
	return fmt.Sprintf("push %016b", i.Mask)
}
