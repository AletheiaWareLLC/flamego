package flamego

type Device interface {
	Store
	Signal()
}
