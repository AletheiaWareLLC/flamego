package intermediate

type Statement struct {
	address uint32
	comment string
}

func (s *Statement) AbsoluteAddress() uint32 {
	return s.address
}

func (s *Statement) RelativeAddress(a uint32) int32 {
	return int32(s.address) - int32(a)
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
