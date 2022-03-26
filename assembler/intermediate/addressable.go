package intermediate

type Addressable interface {
	AbsoluteAddress() uint32
	RelativeAddress(uint32) int32
	SetAddress(uint32)
}
