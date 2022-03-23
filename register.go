package flamego

import (
	"fmt"
)

type Register uint8

// RO: Read-Only
// PR: Privileged - Read-Only except by Interrupt Service Routine
// GP: General Purpose

const (
	R0  Register = iota // RO, Always 0
	R1                  // RO, Always 1
	R2                  // RO, Core Identifier (0-7)
	R3                  // RO, Context Identifier (0-7)
	R4                  // PR, Interrupt Vector Table
	R5                  // PR, Process Identifier
	R6                  // PR, Program Counter (offset from Program Start)
	R7                  // PR, Program Start
	R8                  // PR, Program Limit
	R9                  // PR, Stack Pointer
	R10                 // PR, Stack Start
	R11                 // PR, Stack Limit
	R12                 // PR, Data Start
	R13                 // PR, Data Limit
	R14                 // PR, Reserved
	R15                 // PR, Reserved
	R16                 // GP
	R17                 // GP
	R18                 // GP
	R19                 // GP
	R20                 // GP
	R21                 // GP
	R22                 // GP
	R23                 // GP
	R24                 // GP
	R25                 // GP
	R26                 // GP
	R27                 // GP
	R28                 // GP
	R29                 // GP
	R30                 // GP
	R31                 // GP
)

const (
	RegisterCount = 32

	RCoreIdentifier       = R2
	RContextIdentifier    = R3
	RInterruptVectorTable = R4
	RProcessIdentifier    = R5
	RProgramCounter       = R6
	RProgramStart         = R7
	RProgramLimit         = R8
	RStackPointer         = R9
	RStackStart           = R10
	RStackLimit           = R11
	RDataStart            = R12
	RDataLimit            = R13
)

func (r Register) String() string {
	return fmt.Sprintf("r%d", r)
}
