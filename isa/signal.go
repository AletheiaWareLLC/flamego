package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

type Signal struct {
	DeviceIdRegister flamego.Register
	success          bool
}

func NewSignal(r flamego.Register) *Signal {
	return &Signal{
		DeviceIdRegister: r,
	}
}

func (i *Signal) Load(x flamego.Context) (uint64, uint64, uint64, uint64) {
	i.success = true
	// Load Device Address
	return x.ReadRegister(i.DeviceIdRegister), 0, 0, 0
}

func (i *Signal) Execute(x flamego.Context, a, b, c, d uint64) (uint64, uint64) {
	if !x.IsInterrupted() {
		// Signal only allowed in an interrupt
		x.Error(flamego.InterruptUnsupportedOperationError)
		i.success = false
		return 0, 0
	}
	x.Core().Processor().Signal(int(a))
	return 0, 0
}

func (i *Signal) Format(x flamego.Context, a, b uint64) (uint64, uint64) {
	// Do Nothing
	return 0, 0
}

func (i *Signal) Store(x flamego.Context, a, b uint64) {
	// Do Nothing
}

func (i *Signal) Retire(x flamego.Context) bool {
	if i.success {
		x.IncrementProgramCounter()
	}
	return true
}

func (i *Signal) String() string {
	return fmt.Sprintf("signal %s", i.DeviceIdRegister)
}
