package flamego

const (
	// Unit: Bytes
	BusSizeL1Cache = 8
	BusSizeL2Cache = 64
	BusSizeMemory  = 512
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
