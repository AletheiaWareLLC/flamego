package isa_test

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/isa"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestEncoding(t *testing.T) {
	t.Run("Special", func(t *testing.T) {
		t.Run("Halt", func(t *testing.T) {
			opcode := isa.Encode(isa.NewHalt())
			assert.Equal(t, "00000001000000000000000000000000", fmt.Sprintf("%032b", opcode))
		})
		t.Run("Noop", func(t *testing.T) {
			opcode := isa.Encode(isa.NewNoop())
			assert.Equal(t, "00000001000100000000000000000000", fmt.Sprintf("%032b", opcode))
		})
		t.Run("Sleep", func(t *testing.T) {
			opcode := isa.Encode(isa.NewSleep())
			assert.Equal(t, "00000001001000000000000000000000", fmt.Sprintf("%032b", opcode))
		})
		t.Run("Signal", func(t *testing.T) {
			opcode := isa.Encode(isa.NewSignal(flamego.R31))
			assert.Equal(t, "00000001001100000000000000011111", fmt.Sprintf("%032b", opcode))
		})
		t.Run("Lock", func(t *testing.T) {
			opcode := isa.Encode(isa.NewLock())
			assert.Equal(t, "00000001010000000000000000000000", fmt.Sprintf("%032b", opcode))
		})
		t.Run("Unlock", func(t *testing.T) {
			opcode := isa.Encode(isa.NewUnlock())
			assert.Equal(t, "00000001010100000000000000000000", fmt.Sprintf("%032b", opcode))
		})
		t.Run("Interrupt", func(t *testing.T) {
			opcode := isa.Encode(isa.NewInterrupt(flamego.InterruptBreakpoint))
			assert.Equal(t, "00000001011000000000000000000001", fmt.Sprintf("%032b", opcode))
		})
		t.Run("Uninterrupt", func(t *testing.T) {
			opcode := isa.Encode(isa.NewUninterrupt(flamego.R31))
			assert.Equal(t, "00000001011100000000000000011111", fmt.Sprintf("%032b", opcode))
		})
	})
	t.Run("ControlFlow", func(t *testing.T) {
		t.Run("Jump", func(t *testing.T) {
			t.Run("Forward", func(t *testing.T) {
				t.Run("EZ", func(t *testing.T) {
					opcode := isa.Encode(isa.NewJump(isa.JumpEZ, isa.JumpForward, 22, flamego.R31))
					assert.Equal(t, "01000000000000000000001011011111", fmt.Sprintf("%032b", opcode))
				})
				t.Run("NZ", func(t *testing.T) {
					opcode := isa.Encode(isa.NewJump(isa.JumpNZ, isa.JumpForward, 22, flamego.R31))
					assert.Equal(t, "01010000000000000000001011011111", fmt.Sprintf("%032b", opcode))
				})
				t.Run("LE", func(t *testing.T) {
					opcode := isa.Encode(isa.NewJump(isa.JumpLE, isa.JumpForward, 22, flamego.R31))
					assert.Equal(t, "01100000000000000000001011011111", fmt.Sprintf("%032b", opcode))
				})
				t.Run("LZ", func(t *testing.T) {
					opcode := isa.Encode(isa.NewJump(isa.JumpLZ, isa.JumpForward, 22, flamego.R31))
					assert.Equal(t, "01110000000000000000001011011111", fmt.Sprintf("%032b", opcode))
				})
			})
			t.Run("Backward", func(t *testing.T) {
				t.Run("EZ", func(t *testing.T) {
					opcode := isa.Encode(isa.NewJump(isa.JumpEZ, isa.JumpBackward, 22, flamego.R31))
					assert.Equal(t, "01001000000000000000001011011111", fmt.Sprintf("%032b", opcode))
				})
				t.Run("NZ", func(t *testing.T) {
					opcode := isa.Encode(isa.NewJump(isa.JumpNZ, isa.JumpBackward, 22, flamego.R31))
					assert.Equal(t, "01011000000000000000001011011111", fmt.Sprintf("%032b", opcode))
				})
				t.Run("LE", func(t *testing.T) {
					opcode := isa.Encode(isa.NewJump(isa.JumpLE, isa.JumpBackward, 22, flamego.R31))
					assert.Equal(t, "01101000000000000000001011011111", fmt.Sprintf("%032b", opcode))
				})
				t.Run("LZ", func(t *testing.T) {
					opcode := isa.Encode(isa.NewJump(isa.JumpLZ, isa.JumpBackward, 22, flamego.R31))
					assert.Equal(t, "01111000000000000000001011011111", fmt.Sprintf("%032b", opcode))
				})
			})
		})
		t.Run("Call", func(t *testing.T) {
			opcode := isa.Encode(isa.NewCall(flamego.R31))
			assert.Equal(t, "00000010000000000000000000011111", fmt.Sprintf("%032b", opcode))
		})
		t.Run("Return", func(t *testing.T) {
			opcode := isa.Encode(isa.NewReturn(flamego.R31))
			assert.Equal(t, "00000011000000000000000000011111", fmt.Sprintf("%032b", opcode))
		})
	})
	t.Run("DataMovement", func(t *testing.T) {
		t.Run("LoadC", func(t *testing.T) {
			opcode := isa.Encode(isa.NewLoadC(22, flamego.R31))
			assert.Equal(t, "10000000000000000000001011011111", fmt.Sprintf("%032b", opcode))
		})
		t.Run("Load", func(t *testing.T) {
			opcode := isa.Encode(isa.NewLoad(flamego.R30, 22, flamego.R31))
			assert.Equal(t, "00100000000000000101101111011111", fmt.Sprintf("%032b", opcode))
		})
		t.Run("Store", func(t *testing.T) {
			opcode := isa.Encode(isa.NewStore(flamego.R30, 22, flamego.R31))
			assert.Equal(t, "00101000000000000101101111011111", fmt.Sprintf("%032b", opcode))
		})
		t.Run("Clear", func(t *testing.T) {
			opcode := isa.Encode(isa.NewClear(flamego.R30, 22))
			assert.Equal(t, "00110000000000000101101111000000", fmt.Sprintf("%032b", opcode))
		})
		t.Run("Flush", func(t *testing.T) {
			opcode := isa.Encode(isa.NewFlush(flamego.R30, 22))
			assert.Equal(t, "00111000000000000101101111000000", fmt.Sprintf("%032b", opcode))
		})
		t.Run("Push", func(t *testing.T) {
			opcode := isa.Encode(isa.NewPush(flamego.R31))
			assert.Equal(t, "00000100000000000000000000011111", fmt.Sprintf("%032b", opcode))
		})
		t.Run("Pop", func(t *testing.T) {
			opcode := isa.Encode(isa.NewPop(flamego.R31))
			assert.Equal(t, "00000110000000000000000000011111", fmt.Sprintf("%032b", opcode))
		})
	})
	t.Run("Bitwise", func(t *testing.T) {
		t.Run("Not", func(t *testing.T) {
			opcode := isa.Encode(isa.NewNot(flamego.R29, flamego.R31))
			assert.Equal(t, "00010000000000000000001110111111", fmt.Sprintf("%032b", opcode))
		})
		t.Run("And", func(t *testing.T) {
			opcode := isa.Encode(isa.NewAnd(flamego.R29, flamego.R30, flamego.R31))
			assert.Equal(t, "00010001000000000111101110111111", fmt.Sprintf("%032b", opcode))
		})
		t.Run("Or", func(t *testing.T) {
			opcode := isa.Encode(isa.NewOr(flamego.R29, flamego.R30, flamego.R31))
			assert.Equal(t, "00010010000000000111101110111111", fmt.Sprintf("%032b", opcode))
		})
		t.Run("Xor", func(t *testing.T) {
			opcode := isa.Encode(isa.NewXor(flamego.R29, flamego.R30, flamego.R31))
			assert.Equal(t, "00010011000000000111101110111111", fmt.Sprintf("%032b", opcode))
		})
		t.Run("LeftShift", func(t *testing.T) {
			opcode := isa.Encode(isa.NewLeftShift(flamego.R29, flamego.R30, flamego.R31))
			assert.Equal(t, "00010100000000000111101110111111", fmt.Sprintf("%032b", opcode))
		})
		t.Run("RightShift", func(t *testing.T) {
			opcode := isa.Encode(isa.NewRightShift(flamego.R29, flamego.R30, flamego.R31))
			assert.Equal(t, "00010101000000000111101110111111", fmt.Sprintf("%032b", opcode))
		})
	})
	t.Run("Arithmetic", func(t *testing.T) {
		t.Run("Add", func(t *testing.T) {
			opcode := isa.Encode(isa.NewAdd(flamego.R29, flamego.R30, flamego.R31))
			assert.Equal(t, "00011000000000000111101110111111", fmt.Sprintf("%032b", opcode))
		})
		t.Run("Subtract", func(t *testing.T) {
			opcode := isa.Encode(isa.NewSubtract(flamego.R29, flamego.R30, flamego.R31))
			assert.Equal(t, "00011001000000000111101110111111", fmt.Sprintf("%032b", opcode))
		})
		t.Run("Multiply", func(t *testing.T) {
			opcode := isa.Encode(isa.NewMultiply(flamego.R29, flamego.R30, flamego.R31))
			assert.Equal(t, "00011010000000000111101110111111", fmt.Sprintf("%032b", opcode))
		})
		t.Run("Divide", func(t *testing.T) {
			opcode := isa.Encode(isa.NewDivide(flamego.R29, flamego.R30, flamego.R31))
			assert.Equal(t, "00011011000000000111101110111111", fmt.Sprintf("%032b", opcode))
		})
		t.Run("Modulo", func(t *testing.T) {
			opcode := isa.Encode(isa.NewModulo(flamego.R29, flamego.R30, flamego.R31))
			assert.Equal(t, "00011100000000000111101110111111", fmt.Sprintf("%032b", opcode))
		})
	})
}

func TestDecoding(t *testing.T) {
	t.Run("Special", func(t *testing.T) {
		t.Run("Halt", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00000001000000000000000000000000", 2, 32)
			assert.NoError(t, err)
			_, ok := isa.Decode(uint32(opcode)).(*isa.Halt)
			assert.True(t, ok)
		})
		t.Run("Noop", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00000001000100000000000000000000", 2, 32)
			assert.NoError(t, err)
			_, ok := isa.Decode(uint32(opcode)).(*isa.Noop)
			assert.True(t, ok)
		})
		t.Run("Sleep", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00000001001000000000000000000000", 2, 32)
			assert.NoError(t, err)
			_, ok := isa.Decode(uint32(opcode)).(*isa.Sleep)
			assert.True(t, ok)
		})
		t.Run("Signal", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00000001001100000000000000011111", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.Signal)
			assert.True(t, ok)
			assert.Equal(t, flamego.R31, inst.DeviceIdRegister)
		})
		t.Run("Lock", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00000001010000000000000000000000", 2, 32)
			assert.NoError(t, err)
			_, ok := isa.Decode(uint32(opcode)).(*isa.Lock)
			assert.True(t, ok)
		})
		t.Run("Unlock", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00000001010100000000000000000000", 2, 32)
			assert.NoError(t, err)
			_, ok := isa.Decode(uint32(opcode)).(*isa.Unlock)
			assert.True(t, ok)
		})
		t.Run("Interrupt", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00000001011000000000000000000001", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.Interrupt)
			assert.True(t, ok)
			assert.Equal(t, flamego.InterruptBreakpoint, inst.Value)
		})
		t.Run("Uninterrupt", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00000001011100000000000000011111", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.Uninterrupt)
			assert.True(t, ok)
			assert.Equal(t, flamego.R31, inst.AddressRegister)
		})
	})
	t.Run("ControlFlow", func(t *testing.T) {
		t.Run("Jump", func(t *testing.T) {
			t.Run("Forward", func(t *testing.T) {
				t.Run("EZ", func(t *testing.T) {
					opcode, err := strconv.ParseUint("01000000000000000000001011011111", 2, 32)
					assert.NoError(t, err)
					inst, ok := isa.Decode(uint32(opcode)).(*isa.Jump)
					assert.True(t, ok)
					assert.Equal(t, isa.JumpEZ, inst.ConditionCode)
					assert.Equal(t, isa.JumpForward, inst.Direction)
					assert.Equal(t, uint32(22), inst.Offset)
					assert.Equal(t, flamego.R31, inst.ConditionRegister)
				})
				t.Run("NZ", func(t *testing.T) {
					opcode, err := strconv.ParseUint("01010000000000000000001011011111", 2, 32)
					assert.NoError(t, err)
					inst, ok := isa.Decode(uint32(opcode)).(*isa.Jump)
					assert.True(t, ok)
					assert.Equal(t, isa.JumpNZ, inst.ConditionCode)
					assert.Equal(t, isa.JumpForward, inst.Direction)
					assert.Equal(t, uint32(22), inst.Offset)
					assert.Equal(t, flamego.R31, inst.ConditionRegister)
				})
				t.Run("LE", func(t *testing.T) {
					opcode, err := strconv.ParseUint("01100000000000000000001011011111", 2, 32)
					assert.NoError(t, err)
					inst, ok := isa.Decode(uint32(opcode)).(*isa.Jump)
					assert.True(t, ok)
					assert.Equal(t, isa.JumpLE, inst.ConditionCode)
					assert.Equal(t, isa.JumpForward, inst.Direction)
					assert.Equal(t, uint32(22), inst.Offset)
					assert.Equal(t, flamego.R31, inst.ConditionRegister)
				})
				t.Run("LZ", func(t *testing.T) {
					opcode, err := strconv.ParseUint("01110000000000000000001011011111", 2, 32)
					assert.NoError(t, err)
					inst, ok := isa.Decode(uint32(opcode)).(*isa.Jump)
					assert.True(t, ok)
					assert.Equal(t, isa.JumpLZ, inst.ConditionCode)
					assert.Equal(t, isa.JumpForward, inst.Direction)
					assert.Equal(t, uint32(22), inst.Offset)
					assert.Equal(t, flamego.R31, inst.ConditionRegister)
				})
			})
			t.Run("Backward", func(t *testing.T) {
				t.Run("EZ", func(t *testing.T) {
					opcode, err := strconv.ParseUint("01001000000000000000001011011111", 2, 32)
					assert.NoError(t, err)
					inst, ok := isa.Decode(uint32(opcode)).(*isa.Jump)
					assert.True(t, ok)
					assert.Equal(t, isa.JumpEZ, inst.ConditionCode)
					assert.Equal(t, isa.JumpBackward, inst.Direction)
					assert.Equal(t, uint32(22), inst.Offset)
					assert.Equal(t, flamego.R31, inst.ConditionRegister)
				})
				t.Run("NZ", func(t *testing.T) {
					opcode, err := strconv.ParseUint("01011000000000000000001011011111", 2, 32)
					assert.NoError(t, err)
					inst, ok := isa.Decode(uint32(opcode)).(*isa.Jump)
					assert.True(t, ok)
					assert.Equal(t, isa.JumpNZ, inst.ConditionCode)
					assert.Equal(t, isa.JumpBackward, inst.Direction)
					assert.Equal(t, uint32(22), inst.Offset)
					assert.Equal(t, flamego.R31, inst.ConditionRegister)
				})
				t.Run("LE", func(t *testing.T) {
					opcode, err := strconv.ParseUint("01101000000000000000001011011111", 2, 32)
					assert.NoError(t, err)
					inst, ok := isa.Decode(uint32(opcode)).(*isa.Jump)
					assert.True(t, ok)
					assert.Equal(t, isa.JumpLE, inst.ConditionCode)
					assert.Equal(t, isa.JumpBackward, inst.Direction)
					assert.Equal(t, uint32(22), inst.Offset)
					assert.Equal(t, flamego.R31, inst.ConditionRegister)
				})
				t.Run("LZ", func(t *testing.T) {
					opcode, err := strconv.ParseUint("01111000000000000000001011011111", 2, 32)
					assert.NoError(t, err)
					inst, ok := isa.Decode(uint32(opcode)).(*isa.Jump)
					assert.True(t, ok)
					assert.Equal(t, isa.JumpLZ, inst.ConditionCode)
					assert.Equal(t, isa.JumpBackward, inst.Direction)
					assert.Equal(t, uint32(22), inst.Offset)
					assert.Equal(t, flamego.R31, inst.ConditionRegister)
				})
			})
		})
		t.Run("Call", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00000010000000000000000000011111", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.Call)
			assert.True(t, ok)
			assert.Equal(t, flamego.R31, inst.AddressRegister)
		})
		t.Run("Return", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00000011000000000000000000011111", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.Return)
			assert.True(t, ok)
			assert.Equal(t, flamego.R31, inst.AddressRegister)
		})
	})
	t.Run("DataMovement", func(t *testing.T) {
		t.Run("LoadC", func(t *testing.T) {
			opcode, err := strconv.ParseUint("10000000000000000000001011011111", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.LoadC)
			assert.True(t, ok)
			assert.Equal(t, flamego.R31, inst.DestinationRegister)
			assert.Equal(t, uint32(22), inst.Constant)
		})
		t.Run("Load", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00100000000000000101101111011111", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.Load)
			assert.True(t, ok)
			assert.Equal(t, flamego.R30, inst.AddressRegister)
			assert.Equal(t, uint32(22), inst.Offset)
			assert.Equal(t, flamego.R31, inst.DestinationRegister)
		})
		t.Run("Store", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00101000000000000101101111011111", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.Store)
			assert.True(t, ok)
			assert.Equal(t, flamego.R30, inst.AddressRegister)
			assert.Equal(t, uint32(22), inst.Offset)
			assert.Equal(t, flamego.R31, inst.SourceRegister)
		})
		t.Run("Clear", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00110000000000000101101111000000", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.Clear)
			assert.True(t, ok)
			assert.Equal(t, flamego.R30, inst.AddressRegister)
			assert.Equal(t, uint32(22), inst.Offset)
		})
		t.Run("Flush", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00111000000000000101101111000000", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.Flush)
			assert.True(t, ok)
			assert.Equal(t, flamego.R30, inst.AddressRegister)
			assert.Equal(t, uint32(22), inst.Offset)
		})
		t.Run("Push", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00000100000000000000000000011111", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.Push)
			assert.True(t, ok)
			assert.Equal(t, flamego.R31, inst.SourceRegister)
		})
		t.Run("Pop", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00000110000000000000000000011111", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.Pop)
			assert.True(t, ok)
			assert.Equal(t, flamego.R31, inst.DestinationRegister)
		})
	})
	t.Run("Bitwise", func(t *testing.T) {
		t.Run("Not", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00010000000000000000001110111111", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.Not)
			assert.True(t, ok)
			assert.Equal(t, flamego.R29, inst.SourceRegister)
			assert.Equal(t, flamego.R31, inst.DestinationRegister)
		})
		t.Run("And", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00010001000000000111101110111111", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.And)
			assert.True(t, ok)
			assert.Equal(t, flamego.R29, inst.Source1Register)
			assert.Equal(t, flamego.R30, inst.Source2Register)
			assert.Equal(t, flamego.R31, inst.DestinationRegister)
		})
		t.Run("Or", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00010010000000000111101110111111", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.Or)
			assert.True(t, ok)
			assert.Equal(t, flamego.R29, inst.Source1Register)
			assert.Equal(t, flamego.R30, inst.Source2Register)
			assert.Equal(t, flamego.R31, inst.DestinationRegister)
		})
		t.Run("Xor", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00010011000000000111101110111111", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.Xor)
			assert.True(t, ok)
			assert.Equal(t, flamego.R29, inst.Source1Register)
			assert.Equal(t, flamego.R30, inst.Source2Register)
			assert.Equal(t, flamego.R31, inst.DestinationRegister)
		})
		t.Run("LeftShift", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00010100000000000111101110111111", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.LeftShift)
			assert.True(t, ok)
			assert.Equal(t, flamego.R29, inst.Source1Register)
			assert.Equal(t, flamego.R30, inst.Source2Register)
			assert.Equal(t, flamego.R31, inst.DestinationRegister)
		})
		t.Run("RightShift", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00010101000000000111101110111111", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.RightShift)
			assert.True(t, ok)
			assert.Equal(t, flamego.R29, inst.Source1Register)
			assert.Equal(t, flamego.R30, inst.Source2Register)
			assert.Equal(t, flamego.R31, inst.DestinationRegister)
		})
	})
	t.Run("Arithmetic", func(t *testing.T) {
		t.Run("Add", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00011000000000000111101110111111", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.Add)
			assert.True(t, ok)
			assert.Equal(t, flamego.R29, inst.Source1Register)
			assert.Equal(t, flamego.R30, inst.Source2Register)
			assert.Equal(t, flamego.R31, inst.DestinationRegister)
		})
		t.Run("Subtract", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00011001000000000111101110111111", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.Subtract)
			assert.True(t, ok)
			assert.Equal(t, flamego.R29, inst.Source1Register)
			assert.Equal(t, flamego.R30, inst.Source2Register)
			assert.Equal(t, flamego.R31, inst.DestinationRegister)
		})
		t.Run("Multiply", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00011010000000000111101110111111", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.Multiply)
			assert.True(t, ok)
			assert.Equal(t, flamego.R29, inst.Source1Register)
			assert.Equal(t, flamego.R30, inst.Source2Register)
			assert.Equal(t, flamego.R31, inst.DestinationRegister)
		})
		t.Run("Divide", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00011011000000000111101110111111", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.Divide)
			assert.True(t, ok)
			assert.Equal(t, flamego.R29, inst.Source1Register)
			assert.Equal(t, flamego.R30, inst.Source2Register)
			assert.Equal(t, flamego.R31, inst.DestinationRegister)
		})
		t.Run("Modulo", func(t *testing.T) {
			opcode, err := strconv.ParseUint("00011100000000000111101110111111", 2, 32)
			assert.NoError(t, err)
			inst, ok := isa.Decode(uint32(opcode)).(*isa.Modulo)
			assert.True(t, ok)
			assert.Equal(t, flamego.R29, inst.Source1Register)
			assert.Equal(t, flamego.R30, inst.Source2Register)
			assert.Equal(t, flamego.R31, inst.DestinationRegister)
		})
	})
}
