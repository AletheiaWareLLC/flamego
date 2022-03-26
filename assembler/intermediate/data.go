package intermediate

import (
	"aletheiaware.com/flamego"
	"encoding/binary"
	"fmt"
)

type Data struct {
	Statement
	name  string
	value uint64
	label *Label
}

func NewDataWithName(n string, c string) *Data {
	return &Data{
		Statement: Statement{
			comment: c,
		},
		name: n,
	}
}

func NewDataWithValue(v uint64, c string) *Data {
	return &Data{
		Statement: Statement{
			comment: c,
		},
		value: v,
	}
}

func NewDataWithLabel(l *Label, c string) *Data {
	return &Data{
		Statement: Statement{
			comment: c,
		},
		label: l,
	}
}

func (d *Data) Value() uint64 {
	return d.value
}

func (d *Data) EmittedSize() uint32 {
	return flamego.DataSize
}

func (d *Data) Emit() []byte {
	buffer := make([]byte, 8)
	binary.BigEndian.PutUint64(buffer, d.value)
	return buffer
}

func (d *Data) Link(l Linker) error {
	if d.name != "" {
		label, err := l.Label(d.name)
		if err != nil {
			return err
		}
		d.label = label
		d.value = uint64(d.label.AbsoluteAddress())
	}
	return nil
}

func (d *Data) String() string {
	if d.label != nil {
		return fmt.Sprintf("data %s%s", d.label.name, d.Statement.String())
	}
	return fmt.Sprintf("data %d%s", d.value, d.Statement.String())
}
