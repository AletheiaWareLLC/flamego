package vm

import (
	"aletheiaware.com/flamego"
	"fmt"
)

func NewL1Cache(size int, lower flamego.Store) *Cache {
	return NewCache(size, flamego.LineWidthL1Cache, flamego.BusSizeL1Cache, flamego.OffsetBitsL1Cache, lower)
}

func NewL2Cache(size int, lower flamego.Store) *Cache {
	return NewCache(size, flamego.LineWidthL2Cache, flamego.BusSizeL2Cache, flamego.OffsetBitsL2Cache, lower)
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
	size           int
	lineWidth      int
	lineCount      int
	lines          []*CacheLine
	busWidth       int
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
				if offset != 0 {
					panic("Unaligned Cache Update")
				}
				line := c.lines[index]
				if line.tag != tag {
					writeback := false
					for i := 0; i < c.lineWidth; i++ {
						if line.IsValid(i) && line.IsDirty(i) {
							writeback = true
						}
					}
					if writeback {
						victim := c.CreateAddress(line.tag, index, 0)
						// Swap values with values in bus
						lb := c.lower.Bus()
						for i := 0; i < c.lineWidth; i++ {
							n := lb.Read(i) // New
							nv := lb.IsValid(i)
							o := line.Read(i) // Old
							ov := line.IsValid(i)
							od := line.IsDirty(i)
							if ov && od {
								lb.Write(i, o)
							} else {
								lb.SetValid(i, false)
								lb.SetDirty(i, false)
							}
							if nv {
								line.Write(i, n)
							} else {
								line.SetValid(i, false)
							}
							line.SetDirty(i, false)
						}
						line.tag = tag
						// Issue a write
						c.lowerAddress = victim
						c.lowerOperation = flamego.CacheWrite
						c.lower.Write(victim)
						break // Don't free lower
					}
				}
				lb := c.lower.Bus()
				// Copy from lower bus into cache line
				for i := 0; i < c.lineWidth; i++ {
					if lb.IsValid(i) {
						line.Write(i, lb.Read(i))
						line.SetDirty(i, false)
					}
				}
				line.tag = tag
			} else {
				// If lower was unsuccessful it will get retried
			}
			c.lower.Free()
			c.lowerOperation = flamego.CacheNone
		case flamego.CacheWrite:
			if c.lower.IsSuccessful() {
				tag, index, offset := c.ParseAddress(c.lowerAddress)
				if offset != 0 {
					panic("Unaligned Cache Update")
				}
				line := c.lines[index]
				if line.tag == tag {
					lb := c.lower.Bus()
					for i := 0; i < c.lineWidth; i++ {
						if line.IsValid(i) && lb.IsValid(i) && line.Read(i) == lb.Read(i) {
							line.SetDirty(i, false)
						}
					}
				}
			} else {
				// TODO what if lower write was unsuccessful?
				panic("Not Yet Implemented")
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
		c.isSuccessful = false
		if line.tag == tag {
			c.isSuccessful = true
		}

		switch c.operation {
		case flamego.CacheNone:
			// Do nothing
		case flamego.CacheRead:
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
				a := int(offset)
				for i := 0; i < c.busWidth && a < c.lineWidth; i, a = i+1, a+1 {
					if line.IsValid(i) {
						c.bus.Write(i, line.Read(a))
					} else {
						c.bus.SetValid(i, false)
					}
					c.bus.SetDirty(i, false)
				}
			} else {
				// Issue read request to lower store
				// Clearing all offset bits read whole aligned line
				c.lowerRead((c.address >> c.offsetBits) << c.offsetBits)
			}
		case flamego.CacheWrite:
			if c.isSuccessful {
				a := int(offset)
				for i := 0; i < c.busWidth && a < c.lineWidth; i, a = i+1, a+1 {
					if !c.bus.IsValid(i) || !c.bus.IsDirty(i) {
						continue
					}
					c.bus.SetDirty(i, false)
					line.Write(a, c.bus.Read(i))
				}
			} else {
				victim := c.CreateAddress(line.tag, index, 0)
				writeback := false
				for i := 0; i < c.lineWidth; i++ {
					if line.IsValid(i) && line.IsDirty(i) {
						writeback = true
						break
					}
				}
				if writeback {
					// Write back to lower
					c.lowerWrite(victim, line)
				} else {
					// Repurpose line for tag
					line.tag = tag
					c.isSuccessful = true
					a := int(offset)
					for i := 0; i < a; i++ {
						line.SetValid(i, false)
						line.SetDirty(i, false)
					}
					for i := 0; i < c.busWidth && a < c.lineWidth; i, a = i+1, a+1 {
						if !c.bus.IsValid(i) || !c.bus.IsDirty(i) {
							line.SetValid(i, false)
							line.SetDirty(i, false)
							continue
						}
						c.bus.SetDirty(i, false)
						line.Write(a, c.bus.Read(i))
					}
				}
			}
		case flamego.CacheClear:
			if c.isSuccessful {
				a := int(offset)
				for i := 0; i < flamego.DataSize && a < c.lineWidth; i, a = i+1, a+1 {
					line.SetValid(a, false)
				}
			}
			c.isSuccessful = true
		case flamego.CacheFlush:
			if c.isSuccessful {
				// Only flush if any of the data is dirty
				c.isSuccessful = false
				a := int(offset)
				for i := 0; i < flamego.DataSize && a < c.lineWidth; i, a = i+1, a+1 {
					if line.IsValid(a) && line.IsDirty(a) {
						c.isSuccessful = true
					}
				}
			}
			if c.isSuccessful {
				// Issue write request to lower store
				// Clearing all offset bits to flush whole aligned line
				c.isSuccessful = c.lowerWrite((c.address>>c.offsetBits)<<c.offsetBits, line)
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

func (c *Cache) lowerRead(address uint64) bool {
	if !c.lower.IsBusy() && c.lower.IsFree() {
		c.lowerAddress = address
		c.lowerOperation = flamego.CacheRead
		c.lower.Read(address)
		return true
	}
	return false
}

func (c *Cache) lowerWrite(address uint64, line *CacheLine) bool {
	if !c.lower.IsBusy() && c.lower.IsFree() {
		c.lowerAddress = address
		c.lowerOperation = flamego.CacheWrite
		// Copy values into bus
		lb := c.lower.Bus()
		i := 0
		for ; i < lb.Size() && i < c.lineWidth; i++ {
			if line.IsValid(i) && line.IsDirty(i) {
				lb.Write(i, line.Read(i))
			} else {
				lb.SetValid(i, false)
				lb.SetDirty(i, false)
			}
		}
		for ; i < c.lineWidth; i++ {
			lb.SetValid(i, false)
			lb.SetDirty(i, false)
		}
		c.lower.Write(address)
		return true
	}
	return false
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
