package isa

import (
	"aletheiaware.com/flamego"
	"fmt"
)

const (
	Width1Bit  = 0x1
	Width2Bit  = 0x3
	Width4Bit  = 0xF
	Width5Bit  = 0x1F
	Width8Bit  = 0xFF
	Width10Bit = 0x3FF
	Width16Bit = 0xFFFF
	Width17Bit = 0x1FFFF
	Width22Bit = 0x3FFFFF
	Width26Bit = 0x3FFFFFF

	WidthRegister = Width5Bit
)

func Encode(instruction flamego.Instruction) uint32 {
	switch i := instruction.(type) {
	case *LoadC:
		c := uint32(i.Constant & Width26Bit)
		return (1 << 31) | (c << 5) | uint32(i.DestinationRegister)
	case *Jump:
		cc := uint32(i.ConditionCode & Width2Bit)
		b := uint32(0)
		if i.Direction {
			b = 1
		}
		o := uint32(i.Offset & Width22Bit)
		return (1 << 30) | (cc << 28) | (b << 27) | (o << 5) | uint32(i.ConditionRegister)
	case *Load:
		return (1 << 29) | (i.Offset << 10) | (uint32(i.AddressRegister) << 5) | uint32(i.DestinationRegister)
	case *Store:
		return (1 << 29) | (1 << 27) | (i.Offset << 10) | (uint32(i.AddressRegister) << 5) | uint32(i.SourceRegister)
	case *Clear:
		return (1 << 29) | (2 << 27) | (i.Offset << 10) | (uint32(i.AddressRegister) << 5)
	case *Flush:
		return (1 << 29) | (3 << 27) | (i.Offset << 10) | (uint32(i.AddressRegister) << 5)
	case *Not:
		return (1 << 28) | (uint32(i.SourceRegister) << 5) | uint32(i.DestinationRegister)
	case *And:
		return (1 << 28) | (1 << 24) | (uint32(i.Source2Register) << 10) | (uint32(i.Source1Register) << 5) | uint32(i.DestinationRegister)
	case *Or:
		return (1 << 28) | (2 << 24) | (uint32(i.Source2Register) << 10) | (uint32(i.Source1Register) << 5) | uint32(i.DestinationRegister)
	case *Xor:
		return (1 << 28) | (3 << 24) | (uint32(i.Source2Register) << 10) | (uint32(i.Source1Register) << 5) | uint32(i.DestinationRegister)
	case *LeftShift:
		return (1 << 28) | (4 << 24) | (uint32(i.Source2Register) << 10) | (uint32(i.Source1Register) << 5) | uint32(i.DestinationRegister)
	case *RightShift:
		return (1 << 28) | (5 << 24) | (uint32(i.Source2Register) << 10) | (uint32(i.Source1Register) << 5) | uint32(i.DestinationRegister)
	case *Add:
		return (1 << 28) | (8 << 24) | (uint32(i.Source2Register) << 10) | (uint32(i.Source1Register) << 5) | uint32(i.DestinationRegister)
	case *Subtract:
		return (1 << 28) | (9 << 24) | (uint32(i.Source2Register) << 10) | (uint32(i.Source1Register) << 5) | uint32(i.DestinationRegister)
	case *Multiply:
		return (1 << 28) | (10 << 24) | (uint32(i.Source2Register) << 10) | (uint32(i.Source1Register) << 5) | uint32(i.DestinationRegister)
	case *Divide:
		return (1 << 28) | (11 << 24) | (uint32(i.Source2Register) << 10) | (uint32(i.Source1Register) << 5) | uint32(i.DestinationRegister)
	case *Modulo:
		return (1 << 28) | (12 << 24) | (uint32(i.Source2Register) << 10) | (uint32(i.Source1Register) << 5) | uint32(i.DestinationRegister)
	case *Push:
		return (1 << 26) | uint32(i.Mask)
	case *Pop:
		return (1 << 26) | (1 << 25) | uint32(i.Mask)
	case *Call:
		return (1 << 25) | uint32(i.AddressRegister)
	case *Return:
		return (1 << 25) | (1 << 24)
	case *Halt:
		return (1 << 24)
	case *Noop:
		return (1 << 24) | (1 << 20)
	case *Sleep:
		return (1 << 24) | (2 << 20)
	case *Signal:
		return (1 << 24) | (3 << 20) | uint32(i.DeviceIdRegister)
	case *Lock:
		return (1 << 24) | (4 << 20)
	case *Unlock:
		return (1 << 24) | (5 << 20)
	case *Interrupt:
		return (1 << 24) | (6 << 20) | (uint32(i.Value) & Width8Bit)
	case *Uninterrupt:
		return (1 << 24) | (7 << 20) | uint32(i.AddressRegister)
	}
	panic(fmt.Sprintf("Unrecognize Instruction: %+v\n", instruction))
	return 0
}

func Decode(opcode uint32) flamego.Instruction {
	if (opcode >> 31) == 0x1 {
		c := (opcode >> 5) & Width26Bit
		r := flamego.Register(opcode & WidthRegister)
		return NewLoadC(c, r)
	} else if (opcode >> 30) == 0x1 {
		cc := JumpConditionCode((opcode >> 28) & Width2Bit)
		b := JumpForward
		if ((opcode >> 27) & Width1Bit) == 0x1 {
			b = JumpBackward
		}
		o := (opcode >> 5) & Width22Bit
		r := flamego.Register(opcode & WidthRegister)
		return NewJump(cc, b, o, r)
	} else if (opcode >> 29) == 0x1 {
		o := (opcode >> 10) & Width17Bit
		a := flamego.Register((opcode >> 5) & WidthRegister)
		r := flamego.Register(opcode & WidthRegister)
		switch (opcode >> 27) & Width2Bit {
		case 0:
			return NewLoad(a, o, r)
		case 1:
			return NewStore(a, o, r)
		case 2:
			return NewClear(a, o)
		case 3:
			return NewFlush(a, o)
		}
	} else if (opcode >> 28) == 0x1 {
		s2 := flamego.Register((opcode >> 10) & WidthRegister)
		s1 := flamego.Register((opcode >> 5) & WidthRegister)
		d := flamego.Register(opcode & WidthRegister)
		switch (opcode >> 24) & Width4Bit {
		case 0:
			return NewNot(s1, d)
		case 1:
			return NewAnd(s1, s2, d)
		case 2:
			return NewOr(s1, s2, d)
		case 3:
			return NewXor(s1, s2, d)
		case 4:
			return NewLeftShift(s1, s2, d)
		case 5:
			return NewRightShift(s1, s2, d)
		case 8:
			return NewAdd(s1, s2, d)
		case 9:
			return NewSubtract(s1, s2, d)
		case 10:
			return NewMultiply(s1, s2, d)
		case 11:
			return NewDivide(s1, s2, d)
		case 12:
			return NewModulo(s1, s2, d)
		}
	} else if (opcode >> 26) == 0x1 {
		m := uint16(opcode & Width16Bit)
		if (opcode>>25)&Width1Bit == 0x1 {
			return NewPop(m)
		}
		return NewPush(m)
	} else if (opcode >> 25) == 0x1 {
		if (opcode>>24)&Width1Bit == 0x1 {
			return NewReturn()
		}
		return NewCall(flamego.Register(opcode & WidthRegister))
	} else if (opcode >> 24) == 0x1 {
		switch (opcode >> 20) & Width4Bit {
		case 0:
			return NewHalt()
		case 1:
			return NewNoop()
		case 2:
			return NewSleep()
		case 3:
			return NewSignal(flamego.Register(opcode & WidthRegister))
		case 4:
			return NewLock()
		case 5:
			return NewUnlock()
		case 6:
			return NewInterrupt(flamego.InterruptValue(opcode & Width8Bit))
		case 7:
			return NewUninterrupt(flamego.Register(opcode & WidthRegister))
		}
	}
	panic(fmt.Sprintf("Unrecognized Opcode: 0x%016x %032b\n", uint32(opcode), uint32(opcode)))
	return NewInterrupt(flamego.InterruptUnsupportedOperationError)
}
