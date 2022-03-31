package vm

import (
	"aletheiaware.com/flamego"
	"log"
)

func NewProcessor(cache flamego.Cache, memory flamego.Memory) *Processor {
	return &Processor{
		cache:      cache,
		memory:     memory,
		lockHolder: -1,
	}
}

type Processor struct {
	cores      []flamego.Core
	cache      flamego.Cache
	memory     flamego.Memory
	devices    []flamego.Device
	halted     bool
	lockHolder int
}

func (p *Processor) Cache() flamego.Cache {
	return p.cache
}

func (p *Processor) Core(index int) flamego.Core {
	return p.cores[index]
}

func (p *Processor) AddCore(c flamego.Core) {
	p.cores = append(p.cores, c)
}

func (p *Processor) Device(index int) flamego.Device {
	return p.devices[index]
}

func (p *Processor) AddDevice(d flamego.Device) {
	p.devices = append(p.devices, d)
	d.SetOnSignal(p.Signal)
}

func (p *Processor) Halt() {
	log.Println("Processor Halted")
	p.halted = true
}

func (p *Processor) HasHalted() bool {
	return p.halted
}

func (p *Processor) LockHolder() int {
	return p.lockHolder
}

func (p *Processor) Signal(device int) {
	if device < flamego.CoreCount*flamego.ContextCount {
		// Signal Core
		core := device / flamego.CoreCount
		context := device % flamego.ContextCount
		p.cores[core].(*Core).contexts[context].Signal()
	} else {
		// Signal IO device
		p.devices[device-flamego.CoreCount*flamego.ContextCount].Signal()
	}
}

func (p *Processor) Clock(cycle int) {
	// IO Devices are 500 times slower
	if cycle%500 == 0 {
		for _, d := range p.devices {
			d.Clock(cycle / 500)
		}
	}

	// Main Memory is 100 times slower
	if cycle%100 == 0 {
		log.Println("Memory Clock")
		p.memory.Clock(cycle / 100)
	}

	// L2 Caches are 10 times slower
	if cycle%10 == 0 {
		log.Println("L2 Clock")
		p.cache.Clock(cycle / 10)
	}

	// Clock Each Core
	for _, c := range p.cores {
		c.Clock(cycle)
	}

	// Update Hardware Lock
	if p.lockHolder == -1 {
		for i, c := range p.cores {
			if c.RequiresLock() {
				c.SetAcquiredLock(true)
				p.lockHolder = i
				break
			}
		}
	} else {
		if !p.cores[p.lockHolder].RequiresLock() {
			p.cores[p.lockHolder].SetAcquiredLock(false)
			p.lockHolder = -1
		}
	}
}
