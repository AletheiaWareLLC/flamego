package isa

import (
	"aletheiaware.com/flamego"
	"encoding/binary"
	"fmt"
)

type Call struct {
	AddressRegister flamego.Register
	success         bool
	issued          bool
}

func NewCall(a flamego.Register) *Call {
	return &Call{
		AddressRegister: a,
	}
}

func (i *Call) Load(x flamego.Context) (uint64, uint64, uint64, uint64) {
	i.success = true
	// Load Stack Pointer Register
	a := x.ReadRegister(flamego.RStackPointer)
	// Load Stack Limit Register
	b := x.ReadRegister(flamego.RStackLimit)
	// Load Program Counter Register
	c := x.ReadRegister(flamego.RProgramCounter)
	// Load Address Register
	d := x.ReadRegister(i.AddressRegister)
	return a, b, c, d
}

func (i *Call) Execute(x flamego.Context, a, b, c, d uint64) (uint64, uint64) {
	// Calculate address of next instruction
	c += flamego.InstructionSize
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
	return a, d
}

func (i *Call) Format(x flamego.Context, a, b uint64) (uint64, uint64) {
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
	// Pass Through
	return a, b
}

func (i *Call) Store(x flamego.Context, a, b uint64) {
	if i.success {
		// Update Stack Pointer
		p := a + flamego.DataSize
		x.WriteRegister(flamego.RStackPointer, p)
		// Update Program Counter
		x.WriteRegister(flamego.RProgramCounter, b)
	}
}

func (i *Call) Retire(x flamego.Context) bool {
	// Do Nothing
	return true
}

func (i *Call) String() string {
	return fmt.Sprintf("call %s", i.AddressRegister)
}
