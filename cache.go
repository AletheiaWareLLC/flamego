package flamego

import "fmt"

type CacheOperation uint8

const (
	CacheNone CacheOperation = iota
	CacheRead
	CacheWrite
	CacheClear
	CacheFlush
)

func (o CacheOperation) String() string {
	switch o {
	case CacheNone:
		return "-"
	case CacheRead:
		return "Read"
	case CacheWrite:
		return "Write"
	case CacheClear:
		return "Clear"
	case CacheFlush:
		return "Flush"
	default:
		return fmt.Sprintf("Unrecognized Cache Operation: %d", o)
	}
}

const (
	// Unit: Bytes
	LineWidthL1Cache = 64
	LineWidthL2Cache = 512

	// Unit: Bits
	OffsetBitsL1Cache = 6
	OffsetBitsL2Cache = 9
)

type Cache interface {
	Store
	Clear(uint64)
	Flush(uint64)
}
