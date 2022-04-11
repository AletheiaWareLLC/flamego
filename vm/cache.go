package vm

import (
	"aletheiaware.com/flamego"
	"fmt"
)

func NewL1Cache(size int, lower flamego.Store) *Cache {
	return NewCache(size, flamego.LineWidthL1Cache, flamego.BusSize, flamego.OffsetBitsL1Cache, lower)
}

func NewL2Cache(size int, lower flamego.Store) *Cache {
	return NewCache(size, flamego.LineWidthL2Cache, flamego.BusSize, flamego.OffsetBitsL2Cache, lower)
}

func NewL3Cache(size int, lower flamego.Store) *Cache {
	return NewCache(size, flamego.LineWidthL3Cache, flamego.BusSize, flamego.OffsetBitsL3Cache, lower)
}

func NewCache(size, lineWidth, busWidth, offsetBits int, lower flamego.Store) *Cache {
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
		bus:        NewBus(busWidth),
		tagBits:    tagBits,
		indexBits:  indexBits,
		offsetBits: offsetBits,
		lower:      lower,
		isFree:     true,
	}
}

type Cache struct {
	size           int
	lineWidth      int
	lineCount      int
	lines          []*CacheLine
	bus            *Bus
	tagBits        int
	indexBits      int
	offsetBits     int
	isSuccessful   bool
	isBusy         bool
	isFree         bool
	address        uint64
	operation      flamego.CacheOperation
	lower          flamego.Store
	lowerAddress   uint64
	lowerOperation flamego.CacheOperation
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

func (c *Cache) Operation() flamego.CacheOperation {
	return c.operation
}

func (c *Cache) LowerAddress() uint64 {
	return c.lowerAddress
}

func (c *Cache) LowerOperation() flamego.CacheOperation {
	return c.lowerOperation
}

func (c *Cache) Clock(cycle int) {
	if c.lower.IsBusy() {
		// Do nothing
	} else {
		switch c.lowerOperation {
		case flamego.CacheNone:
			// Do nothing
		case flamego.CacheRead:
			if c.lower.IsSuccessful() {
				tag, index, offset := c.ParseAddress(c.lowerAddress)
				line := c.lines[index]
				if line.tag != tag {
					writeback := false
					start := 0
					for i := 0; i < c.lineWidth; i++ {
						if line.IsValid(i) && line.IsDirty(i) {
							writeback = true
							start = i
							break
						}
					}
					if writeback {
						// The cache line is valid and contains some dirty values
						// The read from lower is now discarded :(
						// the values in the bus will be overwritten with dirty values from line
						// and then a write will be issued to lower.
						// This may repeat until the cache line can be repurposed.

						// align start to data boundary
						for start%flamego.DataSize != 0 {
							start--
						}

						victim := c.CreateAddress(line.tag, index, uint64(start))

						lb := c.lower.Bus()
						for i := 0; i < lb.Size() && start < c.lineWidth; i, start = i+1, start+1 {
							if line.IsValid(start) {
								lb.Write(i, line.Read(start))
							} else {
								lb.SetValid(i, false)
							}
						}
						// Issue a write
						c.lowerAddress = victim
						c.lowerOperation = flamego.CacheWrite
						c.lower.Write(victim)
						break // Don't free lower
					} else {
						// Writeback unnecessary, line can repurposed
						line.tag = tag
						for i := 0; i < c.lineWidth; i++ {
							line.SetValid(i, false)
						}
					}
				}
				lb := c.lower.Bus()
				// Copy from lower bus into cache line
				for i, j := 0, int(offset); i < lb.Size() && j < c.lineWidth; i, j = i+1, j+1 {
					if lb.IsValid(i) {
						line.Write(j, lb.Read(i))
						line.SetDirty(j, false)
					}
				}
			} else {
				// If lower was unsuccessful it will get retried
			}
			c.lower.Free()
			c.lowerOperation = flamego.CacheNone
		case flamego.CacheWrite:
			if c.lower.IsSuccessful() {
				tag, index, offset := c.ParseAddress(c.lowerAddress)
				line := c.lines[index]
				if line.tag == tag {
					lb := c.lower.Bus()
					for i, j := 0, int(offset); i < lb.Size() && j < c.lineWidth; i, j = i+1, j+1 {
						if line.IsValid(j) && lb.IsValid(i) && line.Read(j) == lb.Read(i) {
							line.SetDirty(j, false)
						}
					}
				}
			} else {
				// If lower was unsuccessful it will get retried
			}
			c.lower.Free()
			c.lowerOperation = flamego.CacheNone
		default:
			panic(fmt.Errorf("Unrecognized Lower Cache Operation: %v", c.lowerOperation))
		}
	}

	if c.isBusy {
		tag, index, offset := c.ParseAddress(c.address)

		line := c.lines[index]
		c.isSuccessful = line.tag == tag

		switch c.operation {
		case flamego.CacheNone:
			// Do nothing
		case flamego.CacheRead:
			// Check all values are valid
			for i, j := 0, int(offset); i < c.bus.Size() && j < c.lineWidth; i, j = i+1, j+1 {
				if !line.IsValid(j) {
					c.isSuccessful = false
				}
			}
			if c.isSuccessful {
				// Copy values into bus
				for i, j := 0, int(offset); i < c.bus.Size() && j < c.lineWidth; i, j = i+1, j+1 {
					c.bus.Write(i, line.Read(j))
					c.bus.SetDirty(i, false)
				}
			} else {
				// Issue read request to lower store
				c.lowerRead(c.address)
			}
		case flamego.CacheWrite:
			if c.isSuccessful {
				for i, j := 0, int(offset); i < c.bus.Size() && j < c.lineWidth; i, j = i+1, j+1 {
					if !c.bus.IsValid(i) || !c.bus.IsDirty(i) {
						continue
					}
					c.bus.SetDirty(i, false)
					line.Write(j, c.bus.Read(i))
				}
			} else {
				writeback := false
				start := 0
				for i := 0; i < c.lineWidth; i++ {
					if line.IsValid(i) && line.IsDirty(i) {
						writeback = true
						start = i
						break
					}
				}
				if writeback {
					// align start to data boundary
					for start%flamego.DataSize != 0 {
						start--
					}

					victim := c.CreateAddress(line.tag, index, uint64(start))

					// Write back to lower
					c.lowerWrite(victim, line, start)
				} else {
					// Writeback unnecessary, line can repurposed
					line.tag = tag
					c.isSuccessful = true
					for i := 0; i < c.lineWidth; i++ {
						line.SetValid(i, false)
						line.SetDirty(i, false)
					}
					for i, j := 0, int(offset); i < c.bus.Size() && j < c.lineWidth; i, j = i+1, j+1 {
						if !c.bus.IsValid(i) || !c.bus.IsDirty(i) {
							continue
						}
						line.Write(j, c.bus.Read(i))
						c.bus.SetDirty(i, false)
					}
				}
			}
		case flamego.CacheClear:
			if c.isSuccessful {
				for i, j := 0, int(offset); i < c.bus.Size() && j < c.lineWidth; i, j = i+1, j+1 {
					line.SetValid(j, false)
				}
			}
			c.isSuccessful = true
		case flamego.CacheFlush:
			if c.isSuccessful {
				// Only flush if any of the data is dirty
				c.isSuccessful = false
				for i, j := 0, int(offset); i < c.bus.Size() && j < c.lineWidth; i, j = i+1, j+1 {
					if line.IsValid(j) && line.IsDirty(j) {
						c.isSuccessful = true
					}
				}
			}
			if c.isSuccessful {
				// Issue write request to lower store
				c.lowerWrite(c.address, line, int(offset))
				// Cache still contains dirty data until the lower store write is successful
				c.isSuccessful = false
			} else {
				// Cache doesn't contain the data for the given address, so nothing to flush, success!
				c.isSuccessful = true
			}
		default:
			panic(fmt.Errorf("Unrecognized Cache Operation: %v", c.operation))
		}
		c.isBusy = false
		c.operation = flamego.CacheNone
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
	c.operation = flamego.CacheRead
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
	c.operation = flamego.CacheWrite
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
	c.operation = flamego.CacheClear
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
	c.operation = flamego.CacheFlush
	c.address = address
}

func (c *Cache) lowerRead(address uint64) {
	if !c.lower.IsBusy() && c.lower.IsFree() {
		c.lowerAddress = address
		c.lowerOperation = flamego.CacheRead
		c.lower.Read(address)
	}
}

func (c *Cache) lowerWrite(address uint64, line *CacheLine, offset int) {
	if !c.lower.IsBusy() && c.lower.IsFree() {
		c.lowerAddress = address
		c.lowerOperation = flamego.CacheWrite
		// Copy values into bus
		lb := c.lower.Bus()
		for i := 0; i < lb.Size() && offset < c.lineWidth; i, offset = i+1, offset+1 {
			if line.IsValid(offset) && line.IsDirty(offset) {
				lb.Write(i, line.Read(offset))
			} else {
				lb.SetValid(i, false)
			}
		}
		c.lower.Write(address)
	}
}

func (c *Cache) ParseAddress(a uint64) (uint64, uint64, uint64) {
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

func (c *Cache) CreateAddress(tag, index, offset uint64) uint64 {
	return tag<<(c.indexBits+c.offsetBits) | index<<c.offsetBits | offset
}

type CacheLine struct {
	Bus
	tag uint64
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
