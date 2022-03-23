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
}

func NewStore(a flamego.Register, o uint32, r flamego.Register) *Store {
	return &Store{
		AddressRegister: a,
		Offset:          o,
		SourceRegister:  r,
	}
}

func (i *Store) Load(x flamego.Context) (uint64, uint64, uint64) {
	i.success = true
	// Load Base Register
	a := x.ReadRegister(i.AddressRegister)
	// Load Offset
	b := uint64(i.Offset)
	// Load Source Register
	c := x.ReadRegister(i.SourceRegister)
	return a, b, c
}

func (i *Store) Execute(x flamego.Context, a, b, c uint64) uint64 {
	d := x.Core().DataCache()
	if d.IsBusy() {
		i.success = false // Cache Unavailable
		return 0
	}
	// Copy Data to Bus
	buffer := make([]byte, 8)
	binary.BigEndian.PutUint64(buffer, c)
	for i := 0; i < 8; i++ {
		d.Bus().Write(i, buffer[i])
	}
	// Issue Write Request
	d.Write(a + b)
	return 0
}

func (i *Store) Format(x flamego.Context, a uint64) uint64 {
	d := x.Core().DataCache()
	if d.IsBusy() || !d.IsSuccessful() {
		i.success = false // Cache Miss
		return 0
	}
	return 0
}

func (i *Store) Store(x flamego.Context, a uint64) {
	// Do Nothing
}

func (i *Store) Retire(x flamego.Context) {
	if !i.success {
		return
	}
	// Only proceed if successfull
	x.IncrementProgramCounter()
}

func (i *Store) String() string {
	return fmt.Sprintf("store %s 0x%x %s", i.AddressRegister, i.Offset, i.SourceRegister)
}
