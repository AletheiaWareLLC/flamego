package vm

import (
	"aletheiaware.com/flamego"
)

func NewBus(size int) *Bus {
	return &Bus{
		size:  size,
		valid: make([]bool, size),
		dirty: make([]bool, size),
		data:  make([]byte, size),
	}
}

type Bus struct {
	size  int
	valid []bool
	dirty []bool
	data  []byte
}

func (b *Bus) Size() int {
	return b.size
}

func (b *Bus) IsValid(offset int) bool {
	return b.valid[offset]
}

func (b *Bus) SetValid(offset int, valid bool) {
	b.valid[offset] = valid
}

func (b *Bus) IsDirty(offset int) bool {
	return b.dirty[offset]
}

func (b *Bus) SetDirty(offset int, dirty bool) {
	b.dirty[offset] = dirty
}

func (b *Bus) Data() []byte {
	return b.data
}

func (b *Bus) Read(offset int) byte {
	return b.data[offset]
}

func (b *Bus) Write(offset int, value byte) {
	b.valid[offset] = true
	b.dirty[offset] = true
	b.data[offset] = value
}

func (b *Bus) ReadFrom(bus flamego.Bus) {
	for i := 0; i < b.size; i++ {
		b.valid[i] = bus.IsValid(i)
		b.dirty[i] = bus.IsDirty(i)
		b.data[i] = bus.Read(i)
	}
}

func (b *Bus) WriteTo(bus flamego.Bus, offset int) {
	for i := 0; i < b.size; i++ {
		if b.valid[i] && b.dirty[i] {
			bus.Write(i+offset, b.data[i])
			b.dirty[i] = false
		}
	}
}
