package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type Signal struct {
	DeviceIdRegister flamego.Register
}

func NewSignal(r flamego.Register) *Signal {
	return &Signal{
		DeviceIdRegister: r,
	}
}

func (i *Signal) Load(x flamego.Context) (uint64, uint64, uint64) {
	// Load Device Address
	return x.ReadRegister(i.DeviceIdRegister), 0, 0
}

func (i *Signal) Execute(x flamego.Context, a, b, c uint64) uint64 {
	x.Core().Processor().Signal(int(a))
	return 0
}

func (i *Signal) Format(x flamego.Context, a uint64) uint64 {
	// Do Nothing
	return 0
}

func (i *Signal) Store(x flamego.Context, a uint64) {
	// Do Nothing
}

func (i *Signal) Retire(x flamego.Context) {
	// Do Nothing
}

func (i *Signal) String() string {
	return fmt.Sprintf("signal %s", i.DeviceIdRegister)
}
