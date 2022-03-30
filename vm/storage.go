package vm

import (
	"aletheiaware.com/flamego"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

var _ (flamego.Device) = (*FileStorage)(nil)

var operations = map[uint32]flamego.DeviceOperation{
	flamego.DeviceCommandEnable:  flamego.DeviceEnable,
	flamego.DeviceCommandDisable: flamego.DeviceDisable,
	flamego.DeviceCommandRead:    flamego.DeviceRead,
	flamego.DeviceCommandWrite:   flamego.DeviceWrite,
}

func NewFileStorage(m flamego.Memory, o uint64, s func(int)) *FileStorage {
	return &FileStorage{
		memory:   m,
		offset:   o,
		onSignal: s,
	}
}

type FileStorage struct {
	file            *os.File
	memory          flamego.Memory
	memoryOperation flamego.MemoryOperation
	offset          uint64
	onSignal        func(int)
	isBusy          bool
	operation       flamego.DeviceOperation
	command         uint64
	controller      int
	parameter       uint64
	deviceAddress   uint64
	memoryAddress   uint64
}

func (s *FileStorage) Open(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	s.file = f
	return nil
}

func (s *FileStorage) Close() error {
	if s.file == nil {
		return nil
	}
	return s.file.Close()
}

func (s *FileStorage) Signal() {
	var n string
	if s.file != nil {
		n = s.file.Name()
	}
	log.Println("Storage Signal:", n)
	if s.isBusy {
		panic("Device Already Busy")
	}
	s.isBusy = true
	s.operation = flamego.DeviceNone
}

func (s *FileStorage) Clock(cycle int) {
	var n string
	if s.file != nil {
		n = s.file.Name()
	}
	log.Println("Storage Clock:", cycle, n)
	if s.memory.IsBusy() {
		// Do nothing
	} else {
		switch s.memoryOperation {
		case flamego.MemoryNone:
		// Do nothing
		case flamego.MemoryRead:
			if s.memory.IsSuccessful() {
				// Copy command, parameter, deviceaddress, and memoryaddress from memory bus
				mb := s.memory.Bus()
				buffer := make([]byte, flamego.DeviceControlBlockSize)
				for i := 0; i < flamego.DeviceControlBlockSize; i++ {
					buffer[i] = mb.Read(i)
				}
				s.command = binary.BigEndian.Uint64(buffer[0:8])
				s.parameter = binary.BigEndian.Uint64(buffer[8:16])
				s.deviceAddress = binary.BigEndian.Uint64(buffer[16:24])
				s.memoryAddress = binary.BigEndian.Uint64(buffer[24:32])

				// Split command into controller and operation
				s.controller = int((s.command >> 32) & 0xffffffff)
				op := uint32(s.command & 0xffffffff)
				o, ok := operations[op]
				if !ok {
					panic(fmt.Errorf("Unrecognized Device Operation: %v", op))
				}
				s.operation = o
				log.Println("Loaded Control Block")
				log.Println("Controller:", s.controller)
				log.Println("Operation:", s.operation)
				log.Println("Parameter:", s.parameter)
				log.Println("Device Address:", s.deviceAddress)
				log.Println("Memory Address:", s.memoryAddress)
			} else {
				panic("Not Yet Implemented")
			}
			s.memory.Free()
		case flamego.MemoryWrite:
			if s.memory.IsSuccessful() {
				log.Println("Write Successful")
				s.signalController()
			}
			s.memory.Free()
		default:
			panic(fmt.Errorf("Unrecognized Memory Operation: %v", s.memoryOperation))
		}
		s.memoryOperation = flamego.MemoryNone
	}
	if s.isBusy {
		switch s.operation {
		case flamego.DeviceNone:
			log.Println("Loading Control Block")
			s.LoadControlBlock()
		case flamego.DeviceEnable:
			log.Println("Enabling")
			if err := s.Enable(); err != nil {
				panic(err)
			}
			s.isBusy = false
		case flamego.DeviceDisable:
			log.Println("Disabling")
			if err := s.Disable(); err != nil {
				panic(err)
			}
		case flamego.DeviceRead:
			log.Println("Reading")
			if err := s.Read(); err != nil {
				panic(err)
			}
		case flamego.DeviceWrite:
			log.Println("Writing")
			if err := s.Write(); err != nil {
				panic(err)
			}
		default:
			panic(fmt.Errorf("Unrecognized Device Operation: %v", s.operation))
		}
	}
}

func (s *FileStorage) LoadControlBlock() {
	// Read control block from memory
	if !s.memory.IsBusy() && s.memory.IsFree() {
		s.memoryOperation = flamego.MemoryRead
		s.memory.Read(s.offset)
	}
}

func (s *FileStorage) Enable() error {
	s.isBusy = false
	s.operation = flamego.DeviceNone
	s.signalController()
	return nil
}

func (s *FileStorage) Disable() error {
	if err := s.Close(); err != nil {
		return err
	}
	s.isBusy = false
	s.operation = flamego.DeviceNone
	s.signalController()
	return nil
}

func (s *FileStorage) Read() error {
	// Read from file into memory
	if !s.memory.IsBusy() && s.memory.IsFree() {
		if _, err := s.file.Seek(int64(s.deviceAddress), os.SEEK_SET); err != nil {
			return err
		}
		mb := s.memory.Bus()
		limit := s.parameter
		if s := uint64(mb.Size()); limit > s {
			limit = s
		}
		buffer := make([]byte, limit)
		if count, err := s.file.Read(buffer); err != nil {
			return err
		} else if uint64(count) != limit {
			return fmt.Errorf("Expected to read %d bytes from file, actually read %d", limit, count)
		}
		for i, b := range buffer {
			mb.Write(i, b)
		}
		s.memoryOperation = flamego.MemoryWrite
		s.memory.Write(s.memoryAddress)
		// TODO
		// if s.parameter > limit {
		//   s.parameter-=limit
		//   s.deviceAddress+=limit
		//   s.memoryAddress+=limit
		// } else {
		s.isBusy = false
		s.operation = flamego.DeviceNone
		// }
	}
	return nil
}

func (s *FileStorage) Write() error {
	// Write from memory into file
	if !s.memory.IsBusy() && s.memory.IsFree() {
		panic("Not Yet Implemented")
		s.isBusy = false
		s.operation = flamego.DeviceNone
	}
	return nil
}

func (s *FileStorage) signalController() {
	log.Println("Signalling Controller:", s.controller)
	s.onSignal(s.controller)
}
