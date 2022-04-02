package intermediate

import (
	"fmt"
)

type Align struct {
	Statement
	value uint64
}

func NewAlign(v uint64, c string) *Align {
	return &Align{
		Statement: Statement{
			comment: c,
		},
		value: v,
	}
}

func (a *Align) Value() uint64 {
	return a.value
}

func (a *Align) EmittedSize() uint32 {
	v := uint32(a.value)
	if a.address > v {
		fmt.Println("Cannot align to smaller address")
	}
	return v - a.address
}

func (a *Align) Emit() []byte {
	return make([]byte, a.EmittedSize())
}

func (a *Align) String() string {
	return fmt.Sprintf("align 0x%016x%s", a.value, a.Statement.String())
}
