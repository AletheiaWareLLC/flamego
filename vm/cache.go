package vm

import (
	"aletheiaware.com/flamego"
	"fmt"
)

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
		return fmt.Sprintf("Unrecognized Cache Operation: %T", o)
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

func NewL1Cache(size int, lower flamego.Store) *Cache {
	return newCache(size, LineWidthL1Cache, BusSizeL1Cache, OffsetBitsL1Cache, lower)
}

func NewL2Cache(size int, lower flamego.Store) *Cache {
	return newCache(size, LineWidthL2Cache, BusSizeL2Cache, OffsetBitsL2Cache, lower)
}

func newCache(size, lineWidth, busWidth, offsetBits int, lower flamego.Store) *Cache {
	lineCount := size / lineWidth
	lines := make([]*CacheLine, lineCount)
	for i := 0; i < lineCount; i++ {
		lines[i] = NewCacheLine(lineWidth)
	}
	indexBits := 0
	for ; (1 << indexBits) < lineCount; indexBits++ {
	}
	tagBits := 64 - indexBits - offsetBits
	return &Cache{
		size:       size,
		lineWidth:  lineWidth,
		lineCount:  lineCount,
		lines:      lines,
		busWidth:   busWidth,
		bus:        NewBus(busWidth),
		tagBits:    tagBits,
		indexBits:  indexBits,
		offsetBits: offsetBits,
		lower:      lower,
		isFree:     true,
	}
}

type Cache struct {
	size              int
	lineWidth         int
	lineCount         int
	lines             []*CacheLine
	busWidth          int
	bus               *Bus
	tagBits           int
	indexBits         int
	offsetBits        int
	lower             flamego.Store
	address           uint64
	isSuccessful      bool
	isBusy            bool
	isFree            bool
	operation         CacheOperation
	lowerAddress      uint64
	lowerOperation    CacheOperation
	recentlyUsedCount int
}

func (c *Cache) Size() int {
	return c.size
}

func (c *Cache) LineWidth() int {
	return c.lineWidth
}

func (c *Cache) Lines() []*CacheLine {
	return c.lines
}

func (c *Cache) Bus() flamego.Bus {
	return c.bus
}

func (c *Cache) TagBits() int {
	return c.tagBits
}

func (c *Cache) IndexBits() int {
	return c.indexBits
}

func (c *Cache) OffsetBits() int {
	return c.offsetBits
}

func (c *Cache) IsBusy() bool {
	return c.isBusy
}

func (c *Cache) IsFree() bool {
	return c.isFree
}

func (c *Cache) Free() {
	c.isFree = true
}

func (c *Cache) IsSuccessful() bool {
	return c.isSuccessful
}

func (c *Cache) Address() uint64 {
	return c.address
}

func (c *Cache) Operation() CacheOperation {
	return c.operation
}

func (c *Cache) LowerAddress() uint64 {
	return c.lowerAddress
}

func (c *Cache) LowerOperation() CacheOperation {
	return c.lowerOperation
}

func (c *Cache) RecentlyUsedCount() int {
	return c.recentlyUsedCount
}

func (c *Cache) Clock(cycle int) {
	if c.lower.IsBusy() {
		// Do nothing
	} else {
		switch c.lowerOperation {
		case CacheNone:
			// Do nothing
		case CacheRead:
			if c.lower.IsSuccessful() {
				tag, index, offset := c.parseAddress(c.lowerAddress)
				if offset != 0 {
					panic("Unaligned Cache Update")
				}
				line := c.lines[index]
				if line.tag == tag {
					c.setRecentlyUsed(index)
				} else {
					// TODO Cache line already in use, need to write out to lower store but the bus is currently full of data from previous lower read
					panic("Not Yet Implemented")
				}
				lb := c.lower.Bus()
				// Copy from lower bus into cache line
				for i := 0; i < c.lineWidth; i++ {
					line.Write(i, lb.Read(i))
					line.SetDirty(i, false)
				}
				line.tag = tag
			} else {
				// If lower was unsuccessful it will get retried
			}
			c.lower.Free()
		case CacheWrite:
			if c.lower.IsSuccessful() {
				tag, index, offset := c.parseAddress(c.lowerAddress)
				if offset != 0 {
					panic("Unaligned Cache Update")
				}
				line := c.lines[index]
				if line.tag != tag {
					panic("Lower write doesn't match cache line tag")
				}
				lb := c.lower.Bus()
				for i := 0; i < c.lineWidth; i++ {
					if line.IsValid(i) && lb.IsValid(i) && line.Read(i) == lb.Read(i) {
						line.SetDirty(i, false)
					}
				}
			} else {
				// TODO what if lower write was unsuccessful?
				panic("Not Yet Implemented")
			}
			c.lower.Free()
		default:
			panic(fmt.Errorf("Unrecognized Lower Cache Operation: %T", c.lowerOperation))
		}
		c.lowerOperation = CacheNone
	}

	if c.isBusy {
		tag, index, offset := c.parseAddress(c.address)

		line := c.lines[index]
		c.isSuccessful = false
		if line.tag == tag {
			c.isSuccessful = true
		}

		switch c.operation {
		case CacheNone:
			// Do nothing
		case CacheRead:
			// Check all values are valid
			// Align offset to bus width
			start := int(offset & (^uint64(c.busWidth - 1)))
			for i := 0; i < c.busWidth; i++ {
				if !line.IsValid(i + start) {
					c.isSuccessful = false
				}
			}
			if c.isSuccessful {
				// Copy values into bus
				for i := 0; i < c.busWidth; i++ {
					var d byte
					if a := i + int(offset); a < line.Size() {
						d = line.Read(a)
					}
					c.bus.Write(i, d)
					c.bus.SetDirty(i, false)
				}
				c.setRecentlyUsed(index)
			} else {
				// Issue read request to lower store
				// Clearing all offset bits read whole aligned line
				c.lowerRead((c.address >> c.offsetBits) << c.offsetBits)
			}
		case CacheWrite:
			if c.isSuccessful {
				for i := 0; i < c.busWidth; i++ {
					if !c.bus.IsDirty(i) {
						continue
					}
					v := c.bus.Read(i)
					a := i + int(offset)
					o := line.Read(a)
					line.Write(a, v)
					line.SetDirty(a, o != v) // Dirty if changed
				}
				c.setRecentlyUsed(index)
			} else {
				panic("Not Yet Implemented")
				// TODO pick a line to evict
				//  - if any values are dirty, write them to lower store
				//  - if no values are dirty, reassign line to this tag
			}
		case CacheClear:
			if c.isSuccessful {
				for i := 0; i < 8; i++ {
					line.SetValid(int(offset)+i, false)
				}
				c.setRecentlyUsed(index)
			}
			c.isSuccessful = true
		case CacheFlush:
			if c.isSuccessful {
				// Only flush if any of the data is dirty
				c.isSuccessful = false
				for i := 0; i < 8; i++ {
					if line.IsDirty(int(offset) + i) {
						c.isSuccessful = true
					}
				}
			}
			if c.isSuccessful {
				// Issue write request to lower store
				// Clearing all offset bits to flush whole aligned line
				c.isSuccessful = c.lowerWrite((c.address>>c.offsetBits)<<c.offsetBits, line)
				c.setRecentlyUsed(index)
			} else {
				c.isSuccessful = true
			}
		default:
			panic(fmt.Errorf("Unrecognized Cache Operation: %T", c.operation))
		}
		c.isBusy = false
		c.operation = CacheNone
	}
}

func (c *Cache) Read(address uint64) {
	if address < 0 {
		panic("Cache access error")
	}
	if c.isBusy {
		panic("Cache already busy")
	}
	c.isSuccessful = false
	c.isBusy = true
	c.isFree = false
	c.operation = CacheRead
	c.address = address
}

func (c *Cache) Write(address uint64) {
	if address < 0 {
		panic("Cache access error")
	}
	if c.isBusy {
		panic("Cache already busy")
	}
	c.isSuccessful = false
	c.isBusy = true
	c.isFree = false
	c.operation = CacheWrite
	c.address = address
}

func (c *Cache) Clear(address uint64) {
	if address < 0 {
		panic("Cache access error")
	}
	if c.isBusy {
		panic("Cache already busy")
	}
	c.isSuccessful = false
	c.isBusy = true
	c.isFree = false
	c.operation = CacheClear
	c.address = address
}

func (c *Cache) Flush(address uint64) {
	if address < 0 {
		panic("Cache access error")
	}
	if c.isBusy {
		panic("Cache already busy")
	}
	c.isSuccessful = false
	c.isBusy = true
	c.isFree = false
	c.operation = CacheFlush
	c.address = address
}

func (c *Cache) setRecentlyUsed(index uint64) {
	line := c.lines[index]
	c.recentlyUsedCount++
	if c.recentlyUsedCount > c.lineCount/2 {
		c.recentlyUsedCount = 0
		for _, l := range c.lines {
			l.recentlyUsed = false
		}
	}
	line.recentlyUsed = true
}

func (c *Cache) lowerRead(address uint64) {
	if !c.lower.IsBusy() && c.lower.IsFree() {
		c.lowerAddress = address
		c.lowerOperation = CacheRead
		c.lower.Read(address)
	}
}

func (c *Cache) lowerWrite(address uint64, line *CacheLine) bool {
	if !c.lower.IsBusy() && c.lower.IsFree() {
		c.lowerAddress = address
		c.lowerOperation = CacheWrite
		// Copy values into bus
		b := c.lower.Bus()
		for i := 0; i < b.Size(); i++ {
			if line.IsDirty(i) {
				b.Write(i, line.Read(i))
				b.SetDirty(i, true)
				line.SetDirty(i, false)
			}
		}
		c.lower.Write(address)
		return true
	}
	return false
}

func (c *Cache) parseAddress(a uint64) (uint64, uint64, uint64) {
	tagMask := (uint64(1) << c.tagBits) - 1
	indexMask := (uint64(1) << c.indexBits) - 1
	offsetMask := (uint64(1) << c.offsetBits) - 1
	offset := a & offsetMask
	a >>= c.offsetBits
	index := a & indexMask
	a >>= c.indexBits
	tag := a & tagMask
	return tag, index, offset
}

type CacheLine struct {
	Bus
	tag          uint64
	recentlyUsed bool
}

func NewCacheLine(size int) *CacheLine {
	return &CacheLine{
		Bus: Bus{
			size:  size,
			valid: make([]bool, size),
			dirty: make([]bool, size),
			data:  make([]byte, size),
		},
	}
}

func (c *CacheLine) Tag() uint64 {
	return c.tag
}
