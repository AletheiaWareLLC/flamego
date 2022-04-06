package vm

import (
	"aletheiaware.com/flamego"
	"encoding/binary"
	"fmt"
	"log"
)

const ReadControlBlock = flamego.MemoryOperation(3)

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
	deviceAddress   uint64
	memoryAddress   uint64
	parameter       uint64
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

func (d *Device) DeviceAddress() uint64 {
	return d.deviceAddress
}

func (d *Device) MemoryAddress() uint64 {
	return d.memoryAddress
}

func (d *Device) Parameter() uint64 {
	return d.parameter
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
		case ReadControlBlock:
			if d.memory.IsSuccessful() {
				d.CopyControlBlock()
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
			d.LoadControlBlock()
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

func (d *Device) LoadControlBlock() {
	// Read control block from memory
	if !d.memory.IsBusy() && d.memory.IsFree() {
		log.Println("Loading Control Block")
		d.memoryOperation = ReadControlBlock
		d.memory.Read(d.memoryOffset)
	}
}

func (d *Device) CopyControlBlock() {
	// Copy command, deviceaddress, memoryaddress, and parameter from memory bus
	mb := d.memory.Bus()
	buffer := make([]byte, flamego.DeviceControlBlockSize)
	for i := 0; i < flamego.DeviceControlBlockSize; i++ {
		buffer[i] = mb.Read(i)
	}
	d.command = binary.BigEndian.Uint64(buffer[0:8])
	d.deviceAddress = binary.BigEndian.Uint64(buffer[8:16])
	d.memoryAddress = binary.BigEndian.Uint64(buffer[16:24])
	d.parameter = binary.BigEndian.Uint64(buffer[24:32])

	// Split command into controller and operation
	d.controller = int((d.command >> 32) & 0xffffffff)
	d.operation = flamego.DeviceOperation(d.command & 0xffffffff)
	log.Println("Loaded Control Block")
	log.Println("Controller:", d.controller)
	log.Println("Operation:", d.operation)
	log.Println("Device Address:", d.deviceAddress)
	log.Println("Memory Address:", d.memoryAddress)
	log.Println("Parameter:", d.parameter)
}

func (d *Device) SignalController() {
	log.Println("Signalling Controller:", d.controller)
	d.OnSignal(d.controller)
}
