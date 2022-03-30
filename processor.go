package flamego

type Processor interface {
	Clockable

	Cache() Cache
	Core(int) Core
	AddCore(Core)
	Device(int) Device
	AddDevice(Device)

	Halt()
	HasHalted() bool

	Signal(int)
}
