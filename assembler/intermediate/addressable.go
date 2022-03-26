package intermediate

type Addressable interface {
	AbsoluteAddress() uint32
	RelativeAddress(uint32) uint32
	SetAddress(uint32)
}
