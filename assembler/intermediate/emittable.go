package intermediate

type Emittable interface {
	Addressable
	EmittedSize() uint32
	Emit() []byte
}
