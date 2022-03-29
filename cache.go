package flamego

type Cache interface {
	Store
	Clear(uint64)
	Flush(uint64)
}
