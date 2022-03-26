package intermediate

type Linkable interface {
	Link(Linker) error
}
