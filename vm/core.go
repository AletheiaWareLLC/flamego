package vm

import (
	"aletheiaware.com/flamego"
)

func NewCore(id int, processor flamego.Processor, cache flamego.Cache) *Core {
	return &Core{
		id:         id,
		processor:  processor,
		cache:      cache,
		lockHolder: -1,
	}
}

type Core struct {
	id        int
	processor flamego.Processor
	cache     flamego.Cache
	contexts  []flamego.Context

	next         int
	lockHolder   int
	requiresLock bool
	acquiredLock bool

	loadRegister0    uint64
	loadRegister1    uint64
	loadRegister2    uint64
	loadRegister3    uint64
	executeRegister0 uint64
	executeRegister1 uint64
	formatRegister0  uint64
	formatRegister1  uint64
}

func (c *Core) Id() int {
	return c.id
}

func (c *Core) Processor() flamego.Processor {
	return c.processor
}

func (c *Core) Context(index int) flamego.Context {
	return c.contexts[index]
}

func (c *Core) AddContext(x flamego.Context) {
	c.contexts = append(c.contexts, x)
}

// Index of the next context to fetch an instruction
func (c *Core) NextContext() int {
	return c.next
}

func (c *Core) LockHolder() int {
	return c.lockHolder
}

func (c *Core) RequiresLock() bool {
	return c.requiresLock
}

func (c *Core) AcquiredLock() bool {
	return c.acquiredLock
}

func (c *Core) SetAcquiredLock(acquired bool) {
	c.acquiredLock = acquired
}

func (c *Core) Cache() flamego.Cache {
	return c.cache
}

func (c *Core) LoadRegister0() uint64 {
	return c.loadRegister0
}

func (c *Core) LoadRegister1() uint64 {
	return c.loadRegister1
}

func (c *Core) LoadRegister2() uint64 {
	return c.loadRegister2
}

func (c *Core) LoadRegister3() uint64 {
	return c.loadRegister3
}

func (c *Core) ExecuteRegister0() uint64 {
	return c.executeRegister0
}

func (c *Core) ExecuteRegister1() uint64 {
	return c.executeRegister1
}

func (c *Core) FormatRegister0() uint64 {
	return c.formatRegister0
}

func (c *Core) FormatRegister1() uint64 {
	return c.formatRegister1
}

func (c *Core) Clock(cycle int) {
	// L2 Caches are 10 times slower
	if cycle%10 == 0 {
		c.cache.Clock(cycle / 10)
	}

	// Clock L1 Caches
	for _, c := range c.contexts {
		c.InstructionCache().Clock(cycle)
		c.DataCache().Clock(cycle)
	}

	// Run the pipeline in reverse so data flow in intermediate registers are not affected.
	c.context(7).RetireInstruction()
	c.context(6).StoreData(c.formatRegister0, c.formatRegister1)
	c.formatRegister0, c.formatRegister1 = c.context(5).FormatData(c.executeRegister0, c.executeRegister1)
	c.executeRegister0, c.executeRegister1 = c.context(4).ExecuteOperation(c.loadRegister0, c.loadRegister1, c.loadRegister2, c.loadRegister3)
	c.loadRegister0, c.loadRegister1, c.loadRegister2, c.loadRegister3 = c.context(3).LoadData()
	c.context(2).DecodeInstruction()
	c.context(1).LoadInstruction()
	c.context(0).FetchInstruction()
	c.next = (c.next + 1) % flamego.ContextCount

	// Update Hardware Lock
	if c.lockHolder == -1 {
		for i, x := range c.contexts {
			if x.RequiresLock() {
				if c.acquiredLock {
					x.SetAcquiredLock(true)
					c.lockHolder = i
				} else {
					c.requiresLock = true
				}
				break
			}
		}
	} else {
		if !c.contexts[c.lockHolder].RequiresLock() {
			c.contexts[c.lockHolder].SetAcquiredLock(false)
			c.requiresLock = false
			c.lockHolder = -1
		}
	}
}

func (c *Core) context(stage int) flamego.Context {
	index := c.next - stage
	if index < 0 {
		index += flamego.ContextCount
	}
	return c.contexts[index]
}
