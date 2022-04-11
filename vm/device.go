package vm

import (
	"aletheiaware.com/flamego"
	"encoding/binary"
	"fmt"
	"log"
)

const (
	ReadCommand       = flamego.MemoryOperation(3)
	ReadDeviceAddress = flamego.MemoryOperation(4)
	ReadMemoryAddress = flamego.MemoryOperation(5)
)

func NewDevice(m flamego.Memory, o uint64) *Device {
	return &Device{
		memory:       m,
		memoryOffset: o,
		operations:   make(map[flamego.DeviceOperation]func() error),
	}
}

type Device struct {
	memory          flamego.Memory
	memoryOffset    uint64
	memoryOperation flamego.MemoryOperation
	isBusy          bool
	operations      map[flamego.DeviceOperation]func() error
	operation       flamego.DeviceOperation
	command         uint64
	controller      int
	parameter       uint64
	deviceAddress   uint64
	memoryAddress   uint64
	OnMemoryRead    func() error
	OnMemoryWrite   func() error
	OnSignal        func(int)
}

func (d *Device) MemoryOffset() uint64 {
	return d.memoryOffset
}

func (d *Device) MemoryOperation() flamego.MemoryOperation {
	return d.memoryOperation
}

func (d *Device) IsBusy() bool {
	return d.isBusy
}

func (d *Device) AddOperation(o flamego.DeviceOperation, f func() error) {
	d.operations[o] = f
}

func (d *Device) Operation() flamego.DeviceOperation {
	return d.operation
}

func (d *Device) Command() uint64 {
	return d.command
}

func (d *Device) Controller() int {
	return d.controller
}

func (d *Device) Parameter() uint64 {
	return d.parameter
}

func (d *Device) DeviceAddress() uint64 {
	return d.deviceAddress
}

func (d *Device) MemoryAddress() uint64 {
	return d.memoryAddress
}

func (d *Device) SetOnSignal(s func(int)) {
	d.OnSignal = s
}

func (d *Device) Signal() {
	if d.isBusy {
		panic("Device Already Busy")
	}
	d.isBusy = true
	d.operation = flamego.DeviceNone
}

func (d *Device) Clock(cycle int) {
	if d.memory.IsBusy() {
		// Do nothing
	} else {
		switch d.memoryOperation {
		case flamego.MemoryNone:
			// Do nothing
		case ReadCommand:
			if d.memory.IsSuccessful() {
				d.CopyCommand()
				d.ReadDeviceAddress()
			} else {
				panic("Not Yet Implemented")
			}
			return
		case ReadDeviceAddress:
			if d.memory.IsSuccessful() {
				d.CopyDeviceAddress()
				d.ReadMemoryAddress()
			} else {
				panic("Not Yet Implemented")
			}
			return
		case ReadMemoryAddress:
			if d.memory.IsSuccessful() {
				d.CopyMemoryAddress()
			} else {
				panic("Not Yet Implemented")
			}
			d.memory.Free()
		case flamego.MemoryRead:
			if d.memory.IsSuccessful() {
				if f := d.OnMemoryRead; f != nil {
					if err := f(); err != nil {
						panic(err)
					}
				}
			} else {
				panic("Not Yet Implemented")
			}
			d.memory.Free()
		case flamego.MemoryWrite:
			if d.memory.IsSuccessful() {
				if f := d.OnMemoryWrite; f != nil {
					if err := f(); err != nil {
						panic(err)
					}
				}
			} else {
				panic("Not Yet Implemented")
			}
			d.memory.Free()
		default:
			panic(fmt.Errorf("Unrecognized Memory Operation: %v", d.memoryOperation))
		}
		d.memoryOperation = flamego.MemoryNone
	}
	if d.isBusy {
		switch d.operation {
		case flamego.DeviceNone:
			if !d.memory.IsBusy() && d.memory.IsFree() {
				d.ReadCommand()
			}
		default:
			f, ok := d.operations[d.operation]
			if !ok {
				panic(fmt.Errorf("Unrecognized Device Operation: %v", d.operation))
			}
			if err := f(); err != nil {
				panic(err)
			}
		}
	}
}

func (d *Device) ReadCommand() {
	// Read command from memory
	log.Println("Loading Command")
	d.memoryOperation = ReadCommand
	d.memory.Read(d.memoryOffset + flamego.DataSize*0)
}

func (d *Device) ReadDeviceAddress() {
	// Read device address from memory
	log.Println("Loading Device Address")
	d.memoryOperation = ReadDeviceAddress
	d.memory.Read(d.memoryOffset + flamego.DataSize*1)
}

func (d *Device) ReadMemoryAddress() {
	// Read memory address from memory
	log.Println("Loading Memory Address")
	d.memoryOperation = ReadMemoryAddress
	d.memory.Read(d.memoryOffset + flamego.DataSize*2)
}

func (d *Device) CopyCommand() {
	// Copy command from memory bus
	d.command = d.CopyData()
	// Split command into controller and operation
	d.controller = int((d.command >> 56) & 0xff)
	d.operation = flamego.DeviceOperation((d.command >> 48) & 0xff)
	d.parameter = d.command & 0xffffffffffff
	log.Println("Controller:", d.controller)
	log.Println("Operation:", d.operation)
	log.Println("Parameter:", d.parameter)
}

func (d *Device) CopyDeviceAddress() {
	// Copy device address from memory bus
	d.deviceAddress = d.CopyData()
	log.Println("Device Address:", d.deviceAddress)
}

func (d *Device) CopyMemoryAddress() {
	// Copy memory address from memory bus
	d.memoryAddress = d.CopyData()
	log.Println("Memory Address:", d.memoryAddress)
}

func (d *Device) CopyData() uint64 {
	mb := d.memory.Bus()
	buffer := make([]byte, flamego.DataSize)
	for i := 0; i < flamego.DataSize; i++ {
		buffer[i] = mb.Read(i)
	}
	return binary.BigEndian.Uint64(buffer[:])
}

func (d *Device) SignalController() {
	log.Println("Signalling Controller:", d.controller)
	d.OnSignal(d.controller)
}
