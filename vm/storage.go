package vm

import (
	"aletheiaware.com/flamego"
	"fmt"
	"log"
	"os"
)

var _ (flamego.Device) = (*FileStorage)(nil)

func NewFileStorage(m flamego.Memory, o uint64) *FileStorage {
	fs := &FileStorage{
		Device: *NewDevice(m, o),
	}
	fs.AddOperation(flamego.DeviceStatus, fs.Status)
	fs.AddOperation(flamego.DeviceEnable, fs.Enable)
	fs.AddOperation(flamego.DeviceDisable, fs.Disable)
	fs.AddOperation(flamego.DeviceRead, fs.Read)
	fs.AddOperation(flamego.DeviceWrite, fs.Write)
	return fs
}

type FileStorage struct {
	Device
	file *os.File
}

func (s *FileStorage) File() *os.File {
	return s.file
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

func (s *FileStorage) Status() error {
	// TODO write to memory
	// MemoryAddress[0]<-Manufacturer
	// MemoryAddress[1]<-Current State
	// MemoryAddress[2]<-Serial Number/Product ID
	// MemoryAddress[3]<-Hardware Version
	// MemoryAddress[4]<-Software Version
	// MemoryAddress[5]<-Storage Size/Capacity
	// MemoryAddress[6]
	// MemoryAddress[7]
	return nil
}

func (s *FileStorage) Enable() error {
	if s.file == nil {
		log.Println("Warning: Storage enabled without file being open")
	}
	s.isBusy = false
	s.operation = flamego.DeviceNone
	s.SignalController()
	return nil
}

func (s *FileStorage) Disable() error {
	if err := s.Close(); err != nil {
		return err
	}
	s.isBusy = false
	s.operation = flamego.DeviceNone
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
		count, err := s.file.Read(buffer)
		if err != nil {
			return err
		}
		if uint64(count) != limit {
			return fmt.Errorf("Expected to read %d bytes from file, actually read %d", limit, count)
		}
		for i, b := range buffer {
			mb.Write(i, b)
		}
		s.memoryOperation = flamego.MemoryWrite
		s.memory.Write(s.memoryAddress)
		if s.parameter > limit {
			s.deviceAddress += limit
			s.memoryAddress += limit
			s.parameter -= limit
		} else {
			s.isBusy = false
			s.operation = flamego.DeviceNone
			s.SignalController()
		}
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
