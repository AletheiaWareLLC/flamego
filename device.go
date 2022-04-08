package flamego

import "fmt"

type DeviceOperation uint32

const (
	DeviceNone DeviceOperation = iota
	DeviceStatus
	DeviceEnable
	DeviceDisable
	DeviceRead
	DeviceWrite
)

func (o DeviceOperation) String() string {
	switch o {
	case DeviceNone:
		return "-"
	case DeviceStatus:
		return "Status"
	case DeviceEnable:
		return "Enable"
	case DeviceDisable:
		return "Disable"
	case DeviceRead:
		return "Read"
	case DeviceWrite:
		return "Write"
	default:
		return fmt.Sprintf("Unrecognized Device Operation: %d", o)
	}
}

const (
	DeviceControlBlockAddress = 512
	DeviceControlBlockSize    = 32
)

const (
	DeviceOffsetCommand uint32 = iota
	DeviceOffsetDeviceAddress
	DeviceOffsetMemoryAddress
	DeviceOffsetParameter
)

type Device interface {
	Clockable

	Signal()

	SetOnSignal(func(int))
}
