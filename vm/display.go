package vm

import (
	"aletheiaware.com/flamego"
	"image"
	"image/color"
)

const PixelBytes = 4

var _ (flamego.Device) = (*Display)(nil)

func NewDisplay(m flamego.Memory, o uint64, w, h int) *Display {
	d := &Display{
		Device: *NewDevice(m, o),
		size:   image.Rect(0, 0, w, h),
	}
	d.buffer = image.NewRGBA(d.size)
	d.AddOperation(flamego.DeviceStatus, d.Status)
	d.AddOperation(flamego.DeviceEnable, d.Enable)
	d.AddOperation(flamego.DeviceDisable, d.Disable)
	d.AddOperation(flamego.DeviceWrite, d.FetchFrame)
	d.OnMemoryRead = d.LoadFrame
	return d
}

type Display struct {
	Device
	size   image.Rectangle
	buffer *image.RGBA
}

func (d *Display) Image() image.Image {
	return d.buffer
}

func (d *Display) Status() error {
	// TODO write to memory
	// Command<-Manufacturer
	// Parameter<-Serial Number/Product ID/Hardware & Software Versions
	// DeviceAddress<-Current State
	// MemoryAddress
	return nil
}

func (d *Display) Enable() error {
	c := &color.NRGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}
	for x := 0; x < d.size.Dx(); x++ {
		for y := 0; y < d.size.Dy(); y++ {
			if x == y {
				d.buffer.Set(x, y, c)
			}
		}
	}
	d.isBusy = false
	d.operation = flamego.DeviceNone
	d.SignalController()
	return nil
}

func (d *Display) Disable() error {
	d.isBusy = false
	d.operation = flamego.DeviceNone
	return nil
}

func (d *Display) FetchFrame() error {
	// Read from memory and write to image
	if !d.memory.IsBusy() && d.memory.IsFree() {
		d.memoryOperation = flamego.MemoryRead
		d.memory.Read(d.memoryAddress)
	}
	return nil
}

func (d *Display) LoadFrame() error {
	// Copy from memory bus into image
	width := uint64(d.size.Dx())
	mb := d.memory.Bus()
	size := mb.Size()
	count := 0
	for ; count < size && d.parameter > 0; count += PixelBytes {
		index := d.deviceAddress / PixelBytes
		x := index % width
		y := index / width
		c := &color.NRGBA{
			R: mb.Read(count + 0),
			G: mb.Read(count + 1),
			B: mb.Read(count + 2),
			A: mb.Read(count + 3),
		}
		d.buffer.Set(int(x), int(y), c)
		d.deviceAddress += PixelBytes
		d.memoryAddress += PixelBytes
		d.parameter -= PixelBytes
	}
	if d.parameter == 0 {
		d.isBusy = false
		d.operation = flamego.DeviceNone
		d.SignalController()
	}
	return nil
}
