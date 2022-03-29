package flamego

import (
	"fmt"
)

const InstructionSize = 4

type Instruction interface {
	fmt.Stringer
	Load(Context) (uint64, uint64, uint64)
	Execute(Context, uint64, uint64, uint64) uint64
	Format(Context, uint64) uint64
	Store(Context, uint64)
	Retire(Context) bool
}
