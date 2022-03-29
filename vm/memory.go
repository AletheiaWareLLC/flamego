package vm

import (
	"aletheiaware.com/flamego"
	"fmt"
	"io"
)

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
		return fmt.Sprintf("Unrecognized Memory Operation: %T", o)
	}
}

func NewMemory(size int) *Memory {
	return &Memory{
		size:   size,
		bus:    NewBus(BusSizeMemory),
		data:   make([]byte, size),
		isFree: true,
	}
}

type Memory struct {
	size         int
	bus          *Bus
	data         []byte
	address      uint64
	isSuccessful bool
	isBusy       bool
	isFree       bool
	operation    MemoryOperation
}

func (m *Memory) Size() int {
	return m.size
}

func (m *Memory) Bus() flamego.Bus {
	return m.bus
}

func (m *Memory) Data() []byte {
	return m.data
}

func (m *Memory) Address() uint64 {
	return m.address
}

func (m *Memory) IsBusy() bool {
	return m.isBusy
}

func (m *Memory) IsFree() bool {
	return m.isFree
}

func (m *Memory) Free() {
	m.isFree = true
}

func (m *Memory) IsSuccessful() bool {
	return m.isSuccessful
}

func (m *Memory) Operation() MemoryOperation {
	return m.operation
}

func (m *Memory) Read(address uint64) {
	if address < 0 || address > uint64(m.size) {
		panic("Memory access error")
	}
	if m.isBusy {
		panic("Memory already busy")
	}
	m.isSuccessful = false
	m.isBusy = true
	m.isFree = false
	m.operation = MemoryRead
	m.address = address
}

func (m *Memory) Write(address uint64) {
	if address < 0 || address > uint64(m.size) {
		panic("Memory access error")
	}
	if m.isBusy {
		panic("Memory already busy")
	}
	m.isSuccessful = false
	m.isBusy = true
	m.isFree = false
	m.operation = MemoryWrite
	m.address = address
}

func (m *Memory) Clock(cycle int) {
	if m.isBusy {
		for i := 0; i < m.bus.Size(); i++ {
			switch m.operation {
			case MemoryNone:
				// Do nothing
			case MemoryRead:
				m.bus.Write(i, m.data[m.address+uint64(i)])
			case MemoryWrite:
				if m.bus.IsDirty(i) {
					m.data[m.address+uint64(i)] = m.bus.Read(i)
				}
			default:
				panic(fmt.Errorf("Unrecognized Memory Operation: %T", m.operation))
			}
			m.bus.SetDirty(i, false)
		}
		m.isSuccessful = true
		m.isBusy = false
		m.operation = MemoryNone
	}
}

func (m *Memory) Load(r io.Reader) (int, error) { //heysup? like ketchup, but made of hey
	d, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}
	return copy(m.data, d), nil
}

func (m *Memory) Set(address uint64, data []byte) {
	for i := 0; i < len(data); i++ {
		m.data[address+uint64(i)] = data[i]
	}
}
