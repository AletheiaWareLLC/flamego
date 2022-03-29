package assembler

import (
	"aletheiaware.com/flamego/assembler/intermediate"
	"fmt"
	"io"
)

type Assembler interface {
	io.ReaderFrom
	io.WriterTo
	AddConstant(string, *intermediate.Data) error
	AddLabel(string, *intermediate.Label) error
	AddStatement(intermediate.Addressable) error
}

func NewAssembler() Assembler {
	return &assembler{
		constants: make(map[string]*intermediate.Data),
		labels:    make(map[string]*intermediate.Label),
	}
}

type assembler struct {
	constants  map[string]*intermediate.Data
	labels     map[string]*intermediate.Label
	statements []intermediate.Addressable
}

func (a *assembler) Constant(name string) (*intermediate.Data, error) {
	if c, ok := a.constants[name]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("Constant '%s' Not Declared", name)
}

func (a *assembler) Label(name string) (*intermediate.Label, error) {
	if l, ok := a.labels[name]; ok {
		return l, nil
	}
	return nil, fmt.Errorf("Label '%s' Not Declared", name)
}

func (a *assembler) ReadFrom(reader io.Reader) (int64, error) {
	l := NewLexer()
	count, err := l.Read(reader)
	if err != nil {
		return 0, err
	}
	if err := NewParser(a, l).Parse(); err != nil {
		return 0, err
	}
	return count, nil
}

func (a *assembler) AddConstant(n string, d *intermediate.Data) error {
	if _, ok := a.constants[n]; ok {
		return fmt.Errorf("Duplicate Constant '%s'", n)
	}
	a.constants[n] = d
	return nil
}

func (a *assembler) AddLabel(n string, l *intermediate.Label) error {
	if _, ok := a.labels[n]; ok {
		return fmt.Errorf("Duplicate Label '%s'", n)
	}
	a.labels[n] = l
	return nil
}

func (a *assembler) AddStatement(s intermediate.Addressable) error {
	a.statements = append(a.statements, s)
	return nil
}

func (a *assembler) WriteTo(writer io.Writer) (int64, error) {
	var address uint32
	for _, s := range a.statements {
		s.SetAddress(address)
		if e, ok := s.(intermediate.Emittable); ok {
			address += e.EmittedSize()
		}
	}
	for _, s := range a.statements {
		if l, ok := s.(intermediate.Linkable); ok {
			if err := l.Link(a); err != nil {
				return 0, err
			}
		}
	}
	var count int64
	for _, s := range a.statements {
		if e, ok := s.(intermediate.Emittable); ok {
			n, err := writer.Write(e.Emit())
			if err != nil {
				return 0, err
			}
			count += int64(n)
		}
	}
	return count, nil
}

type Error struct {
	Line    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Line: %d: %s", e.Line, e.Message)
}
