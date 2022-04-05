package flamego

const ContextCount = 8

type Context interface {
	Id() int

	Core() Core

	IsValid() bool
	IsAsleep() bool
	Sleep()
	Error(InterruptValue)
	IsInterrupted() bool
	SetInterrupted(bool)
	Signal()
	IsSignalled() bool

	RequiresLock() bool
	SetRequiresLock(bool)
	AcquiredLock() bool
	SetAcquiredLock(bool)

	FetchInstruction()
	LoadInstruction()
	DecodeInstruction()
	LoadData() (uint64, uint64, uint64, uint64)
	ExecuteOperation(uint64, uint64, uint64, uint64) (uint64, uint64)
	FormatData(uint64, uint64) (uint64, uint64)
	StoreData(uint64, uint64)
	RetireInstruction()

	ReadRegister(Register) uint64
	WriteRegister(Register, uint64)
	IncrementProgramCounter()
	SetProgramCounter(uint64)
}
