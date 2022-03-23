package flamego

type Processor interface {
	Clockable

	Cache() Cache
	Core(int) Core
	AddCore(Core)
	AddDevice(Device)

	Halt()
	HasHalted() bool

	Signal(int)
}
