package flamego

type Device interface {
	Clockable

	Signal()
}
