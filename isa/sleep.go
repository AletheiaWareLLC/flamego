package isa

import (
	"aletheiaware.com/flamego"
)

type Sleep struct {
}

func NewSleep() *Sleep {
	return &Sleep{}
}

func (i *Sleep) Load(x flamego.Context) (uint64, uint64, uint64) {
	// Do Nothing
	return 0, 0, 0
}

func (i *Sleep) Execute(x flamego.Context, a, b, c uint64) uint64 {
	x.Sleep()
	return 0
}

func (i *Sleep) Format(x flamego.Context, a uint64) uint64 {
	// Do Nothing
	return 0
}

func (i *Sleep) Store(x flamego.Context, a uint64) {
	// Do Nothing
}

func (i *Sleep) Retire(x flamego.Context) bool {
	// Do Nothing
	return true
}

func (i *Sleep) String() string {
	return "sleep"
}
