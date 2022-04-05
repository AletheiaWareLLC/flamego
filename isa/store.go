package isa

import (
	"aletheiaware.com/flamego"
	"encoding/binary"
	"fmt"
)

type Store struct {
	AddressRegister flamego.Register
	Offset          uint32
	SourceRegister  flamego.Register
	success         bool
	issued          bool
}

func NewStore(a flamego.Register, o uint32, r flamego.Register) *Store {
	return &Store{
		AddressRegister: a,
		Offset:          o,
		SourceRegister:  r,
	}
}

func (i *Store) Load(x flamego.Context) (uint64, uint64, uint64, uint64) {
	i.success = true
	// Load Base Register
	a := x.ReadRegister(i.AddressRegister)
	// Load Offset
	b := uint64(i.Offset)
	// Load Source Register
	c := x.ReadRegister(i.SourceRegister)
	return a, b, c, 0
}

func (i *Store) Execute(x flamego.Context, a, b, c, d uint64) (uint64, uint64) {
	if !i.issued {
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
		l1d.Write(a + b)
		i.issued = true
	}
	return 0, 0
}

func (i *Store) Format(x flamego.Context, a, b uint64) (uint64, uint64) {
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
	return 0, 0
}

func (i *Store) Store(x flamego.Context, a, b uint64) {
	// Do Nothing
}

func (i *Store) Retire(x flamego.Context) bool {
	if i.success {
		x.IncrementProgramCounter()
		return true
	}
	return false
}

func (i *Store) String() string {
	return fmt.Sprintf("store %s 0x%x %s", i.AddressRegister, i.Offset, i.SourceRegister)
}
