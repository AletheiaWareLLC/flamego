package flamego

import "fmt"

type MemoryOperation uint8

const (
	MemoryNone MemoryOperation = iota
	MemoryRead
	MemoryWrite
)

func (o MemoryOperation) String() string {
	switch o {
	case MemoryNone:
		return "-"
	case MemoryRead:
		return "Read"
	case MemoryWrite:
		return "Write"
	default:
		return fmt.Sprintf("Unrecognized Memory Operation: %d", o)
	}
}

type Memory interface {
	Store
}
