package flamego

type Store interface {
	Clockable

	Bus() Bus

	IsBusy() bool
	IsFree() bool
	Free()
	IsSuccessful() bool

	Read(uint64)
	Write(uint64)
}
