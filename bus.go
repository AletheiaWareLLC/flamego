package flamego

const (
	// Unit: Bytes
	BusSize = 8
)

type Bus interface {
	Size() int
	IsValid(int) bool
	SetValid(int, bool)
	IsDirty(int) bool
	SetDirty(int, bool)
	Read(int) byte
	Write(int, byte)
	ReadFrom(Bus)
	WriteTo(Bus, int)
}
