package flamego

const CoreCount = 8

type Core interface {
	Clockable

	Id() int
	Processor() Processor
	Context(int) Context
	Cache() Cache

	AddContext(Context)

	RequiresLock() bool
	SetAcquiredLock(bool)
}
