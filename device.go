package flamego

import "fmt"

type DeviceOperation uint32

const (
	DeviceNone DeviceOperation = iota
	DeviceEnable
	DeviceDisable
	DeviceRead
	DeviceWrite
)

func (o DeviceOperation) String() string {
	switch o {
	case DeviceNone:
		return "-"
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

const DeviceControlBlockSize = 32

const (
	DeviceCommandEnable uint32 = iota
	DeviceCommandDisable
	DeviceCommandRead
	DeviceCommandWrite
)

const (
	DeviceOffsetCommand uint32 = iota
	DeviceOffsetParameter
	DeviceOffsetDeviceAddress
	DeviceOffsetMemoryAddress
)

type Device interface {
	Clockable

	Signal()
}
