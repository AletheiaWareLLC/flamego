package vm

import (
	"aletheiaware.com/flamego"
)

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
	isRead            bool
	lowerReadPending  bool
	lowerWritePending bool
	lowerAddress      uint64
	recentlyUsedCount int
}

func (c *Cache) Size() int {
	return c.size
}

func (c *Cache) Lines() []*CacheLine {
	return c.lines
}

func (c *Cache) Bus() flamego.Bus {
	return c.bus
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

func (c *Cache) Clock(cycle int) {
	if !c.lower.IsBusy() {
		if c.lowerReadPending {
			c.lowerReadPending = false
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
			}
			c.Free()
		} else if c.lowerWritePending {
			c.lowerWritePending = false
			c.Free()
		}
	}

	if c.isBusy {
		tag, index, offset := c.parseAddress(c.address)

		line := c.lines[index]
		c.isSuccessful = false
		if line.tag == tag {
			c.isSuccessful = true
			c.setRecentlyUsed(index)
		}

		if c.isRead {
			// Check all values are valid
			for i := 0; i < c.busWidth; i++ {
				if !line.IsValid(i + int(offset)) {
					c.isSuccessful = false
				}
			}
			if c.isSuccessful {
				// Copy values into bus
				for i := 0; i < c.busWidth; i++ {
					c.bus.Write(i, line.Read(i+int(offset)))
					c.bus.SetDirty(i, false)
				}
			} else {
				// Issue read request to lower store
				c.lowerRead((c.address >> c.offsetBits) << c.offsetBits)
			}
		} else {
			// TODO
			panic("Not Yet Implemented")
		}
		c.isBusy = false
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
	c.isRead = true
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
	c.isRead = false
	c.address = address
}

func (c *Cache) Clear(address uint64) {
	// TODO
	panic("Not Yet Implemented")
}

func (c *Cache) Flush(address uint64) bool {
	// TODO
	panic("Not Yet Implemented")
	return true
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
		c.lowerReadPending = true
		c.lowerAddress = address
		c.lower.Read(address)
	}
}

func (c *Cache) lowerWrite(address uint64) {
	if !c.lower.IsBusy() && c.lower.IsFree() {
		c.lowerWritePending = true
		c.lowerAddress = address
		// TODO Copy values into lower bus
		panic("Not Yet Implemented")
		c.lower.Write(address)
	}
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
