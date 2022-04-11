package isa

import (
	"aletheiaware.com/flamego"
	"encoding/binary"
)

type Return struct {
	success bool
	issued  bool
}

func NewReturn() *Return {
	return &Return{}
}

func (i *Return) Load(x flamego.Context) (uint64, uint64, uint64, uint64) {
	i.success = true
	// Load Stack Pointer Register
	a := x.ReadRegister(flamego.RStackPointer)
	// Load Stack Start Register
	b := x.ReadRegister(flamego.RStackStart)
	return a, b, 0, 0
}

func (i *Return) Execute(x flamego.Context, a, b, c, d uint64) (uint64, uint64) {
	// Decrement Stack Pointer
	a -= flamego.DataSize

	if a < b {
		x.Error(flamego.InterruptStackUnderflowError)
		i.success = false
	} else if !i.issued {
		l1d := x.DataCache()
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

func (i *Return) Format(x flamego.Context, a, b uint64) (uint64, uint64) {
	if !i.success {
		return 0, 0
	}
	l1d := x.DataCache()
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

func (i *Return) Store(x flamego.Context, a, b uint64) {
	if i.success {
		// Update Stack Pointer
		x.WriteRegister(flamego.RStackPointer, a)
		// Update Program Counter
		x.WriteRegister(flamego.RProgramCounter, b)
	}
}

func (i *Return) Retire(x flamego.Context) bool {
	return i.success
}

func (i *Return) String() string {
	return "return"
}
