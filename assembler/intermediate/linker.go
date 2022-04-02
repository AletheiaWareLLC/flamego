package intermediate

type Linker interface {
	Addressables() []Addressable
	Constant(string) (*Data, error)
	Constants() map[string]*Data
	Label(string) (*Label, error)
	Labels() map[string]*Label
}
