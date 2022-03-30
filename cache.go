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

type Cache interface {
	Store
	Clear(uint64)
	Flush(uint64)
}
