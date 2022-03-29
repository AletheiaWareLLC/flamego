package vm

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"encoding/binary"
)

func NewContext(id int, c *Core) *Context {
	return &Context{
		id:            id,
		core:          c,
		status:        "asleep",
		isAsleep:      true,
		nextInterrupt: -1,
	}
}

type Context struct {
	id        int
	core      *Core
	registers [flamego.RegisterCount]uint64

	status        string
	isValid       bool
	isAsleep      bool
	sleepCycles   int
	isInterrupted bool
	nextInterrupt flamego.InterruptValue
	isSignalled   bool
	requiresLock  bool
	acquiredLock  bool

	opcode      uint32
	instruction flamego.Instruction
}

func (x *Context) Id() int {
	return x.id
}

func (x *Context) Core() flamego.Core {
	return x.core
}

func (x *Context) Status() string {
	return x.status
}

func (x *Context) IsValid() bool {
	return x.isValid
}

func (x *Context) IsAsleep() bool {
	return x.isAsleep
}

func (x *Context) Sleep() {
	x.isAsleep = true
	x.status = "asleep"
}

func (x *Context) SleepCycles() int {
	return x.sleepCycles
}

func (x *Context) IsInterrupted() bool {
	return x.isInterrupted
}

func (x *Context) NextInterrupt() flamego.InterruptValue {
	return x.nextInterrupt
}

func (x *Context) SetInterrupted(i bool) {
	x.isInterrupted = i
}

func (x *Context) Error(value flamego.InterruptValue) {
	if x.isInterrupted {
		panic("Double Interrupt")
	}
	x.nextInterrupt = value
}

func (x *Context) Signal() {
	x.isSignalled = true
}

func (x *Context) IsSignalled() bool {
	return x.isSignalled
}

func (x *Context) Opcode() uint32 {
	return x.opcode
}

func (x *Context) Instruction() flamego.Instruction {
	return x.instruction
}

func (x *Context) RequiresLock() bool {
	return x.requiresLock
}

func (x *Context) SetRequiresLock(required bool) {
	x.requiresLock = required
}

func (x *Context) AcquiredLock() bool {
	return x.acquiredLock
}

func (x *Context) SetAcquiredLock(acquired bool) {
	x.acquiredLock = acquired
}

func (x *Context) FetchInstruction() {
	x.isValid = true
	if x.nextInterrupt >= 0 {
		x.status = "interrupted"
	} else if !x.isInterrupted && !x.acquiredLock && x.isSignalled {
		x.status = "signalled"
		x.nextInterrupt = flamego.InterruptSignal
		x.isAsleep = false
		x.isSignalled = false
		x.sleepCycles = 0
	} else if !x.isAsleep {
		pc := x.ReadRegister(flamego.RProgramCounter)
		if !x.isInterrupted {
			pc += x.ReadRegister(flamego.RProgramStart)
		}
		is := x.core.InstructionCache()
		if is.IsBusy() || !is.IsFree() {
			x.status = "cache busy"
			x.isValid = false
		} else {
			is.Read(pc)
			x.status = "fetched instruction"
		}
	} else {
		x.sleepCycles++
	}
}

func (x *Context) LoadInstruction() {
	if !x.isValid {
		return
	}
	if x.isAsleep {
		x.sleepCycles++
		return
	}
	if !x.isInterrupted && x.nextInterrupt != -1 {
		x.opcode = isa.Encode(isa.NewInterrupt(x.nextInterrupt))
		x.nextInterrupt = -1
	} else {
		is := x.core.InstructionCache()
		if is.IsBusy() {
			x.status = "cache busy"
			x.isValid = false
		} else {
			if is.IsSuccessful() {
				bus := is.Bus()
				x.opcode = binary.BigEndian.Uint32([]byte{
					bus.Read(0),
					bus.Read(1),
					bus.Read(2),
					bus.Read(3),
				})
				x.status = "loaded instruction"
			} else {
				x.status = "cache miss"
				x.isValid = false
			}
			is.Free()
		}
	}
}

func (x *Context) DecodeInstruction() {
	if !x.isValid {
		return
	}
	if x.isAsleep {
		x.sleepCycles++
		return
	}
	x.instruction = isa.Decode(x.opcode)
	x.status = "decoded instruction"
}

func (x *Context) LoadData() (a, b, c uint64) {
	if !x.isValid {
		return
	}
	if x.isAsleep {
		x.sleepCycles++
		return
	}
	a, b, c = x.instruction.Load(x)
	x.status = "loaded data"
	return
}

func (x *Context) ExecuteOperation(a, b, c uint64) (r uint64) {
	if !x.isValid {
		return
	}
	if x.isAsleep {
		x.sleepCycles++
		return
	}
	r = x.instruction.Execute(x, a, b, c)
	x.status = "executed operation"
	return
}

func (x *Context) FormatData(a uint64) (r uint64) {
	if !x.isValid {
		return
	}
	if x.isAsleep {
		x.sleepCycles++
		return
	}
	r = x.instruction.Format(x, a)
	x.status = "formatted data"
	return
}

func (x *Context) StoreData(a uint64) {
	if !x.isValid {
		return
	}
	if x.isAsleep {
		x.sleepCycles++
		return
	}
	x.instruction.Store(x, a)
	x.status = "stored data"
}

func (x *Context) RetireInstruction() {
	if !x.isValid {
		return
	}
	if x.isAsleep {
		x.sleepCycles++
		return
	}
	x.instruction.Retire(x)
	x.opcode = 0
	x.instruction = nil
	x.status = "retired instruction"
}

func (x *Context) ReadRegister(register flamego.Register) uint64 {
	if register < flamego.R0 || register > flamego.R31 {
		x.Error(flamego.InterruptRegisterAccessError)
		return 0
	}
	switch register {
	case flamego.R0:
		return 0
	case flamego.R1:
		return 1
	case flamego.R2:
		return uint64(x.core.Id())
	case flamego.R3:
		return uint64(x.id)
	default:
		return x.registers[register]
	}
}

func (x *Context) WriteRegister(register flamego.Register, value uint64) {
	if register < flamego.R0 || register > flamego.R31 {
		x.Error(flamego.InterruptRegisterAccessError)
		return
	}
	switch register {
	case flamego.R0, flamego.R1, flamego.R2, flamego.R3:
		x.Error(flamego.InterruptRegisterAccessError)
		return
	case flamego.R4, flamego.R5, flamego.R6, flamego.R7, flamego.R8, flamego.R9, flamego.R10, flamego.R11, flamego.R12, flamego.R13, flamego.R14, flamego.R15:
		if !x.isInterrupted {
			x.Error(flamego.InterruptRegisterAccessError)
			return
		}
		fallthrough
	default:
		x.registers[register] = value
	}
}

func (x *Context) IncrementProgramCounter() {
	x.SetProgramCounter(x.registers[flamego.RProgramCounter] + flamego.InstructionSize)
}

func (x *Context) SetProgramCounter(pc uint64) {
	x.registers[flamego.RProgramCounter] = pc
}
