package intermediate

type Linker interface {
	Constant(string) (*Data, error)
	Label(string) (*Label, error)
}
