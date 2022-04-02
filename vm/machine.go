package vm

import (
	"aletheiaware.com/flamego"
	"log"
)

type Machine struct {
	Processor *Processor
	Memory    *Memory

	Tick int
}

func NewMachine() *Machine {
	memory := NewMemory(flamego.SizeMemory)
	l2Cache := NewL2Cache(flamego.SizeL2Cache, memory)
	processor := NewProcessor(l2Cache, memory)
	for i := 0; i < flamego.CoreCount; i++ {
		l1ICache := NewL1Cache(flamego.SizeL1Cache, l2Cache)
		l1DCache := NewL1Cache(flamego.SizeL1Cache, l2Cache)
		core := NewCore(i, processor, l1ICache, l1DCache)
		processor.AddCore(core)
		for j := 0; j < flamego.ContextCount; j++ {
			core.AddContext(NewContext(j, core))
		}
	}
	return &Machine{
		Processor: processor,
		Memory:    memory,
	}
}

func (m *Machine) Clock() {
	if m.Processor.HasHalted() {
		log.Println("Processor Halted")
		return
	}

	// Tick Processor
	m.Processor.Clock(m.Tick)

	m.Tick++
}
