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
		// TODO
		/*
		   Jump
		   Call
		   Return
		*/
	})
	t.Run("DataMovement", func(t *testing.T) {
		// TODO
		/*
		   LoadC
		   Load
		   Store
		   Clear
		   Flush
		   Push
		   Pop
		*/
	})
	t.Run("Bitwise", func(t *testing.T) {
		// TODO
	})
	t.Run("Arithmetic", func(t *testing.T) {
		// TODO
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
}
