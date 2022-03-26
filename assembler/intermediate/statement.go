package intermediate

type Statement struct {
	address uint32
	comment string
}

func (s *Statement) AbsoluteAddress() uint32 {
	return s.address
}

func (s *Statement) RelativeAddress(a uint32) uint32 {
	return s.address - a
}

func (s *Statement) SetAddress(a uint32) {
	s.address = a
}

func (s *Statement) String() string {
	if s.comment != "" {
		return " // " + s.comment
	}
	return ""
}
