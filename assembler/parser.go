package assembler

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/assembler/intermediate"
	"aletheiaware.com/flamego/isa"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Parser interface {
	Parse() error
}

func NewParser(a Assembler, l Lexer) Parser {
	return &parser{
		assembler: a,
		lexer:     l,
	}
}

type parser struct {
	assembler Assembler
	lexer     Lexer
}

func (p *parser) Parse() error {
	p.lexer.Move()
	for !p.lexer.CurrentIs(CategoryEOF) {
		if p.lexer.CurrentIs(CategoryPadding) {
			p.lexer.Move()
			count, err := p.matchNumber()
			if err != nil {
				return err
			}
			for i := int64(1); i <= count; i++ {
				p.assembler.AddStatement(intermediate.NewDataWithValue(0, fmt.Sprintf("padding %d/%d", i, count)))
			}
		} else if p.lexer.CurrentIs(CategoryUpperName) {
			name := p.lexer.Current().Value
			p.lexer.Move()
			d, err := p.matchData()
			if err != nil {
				return err
			}
			p.assembler.AddConstant(name, d)
		} else {
			s, err := p.matchStatement()
			if err != nil {
				return err
			}
			p.assembler.AddStatement(s)
		}
	}
	return nil
}

func (p *parser) matchData() (*intermediate.Data, error) {
	if p.lexer.CurrentIs(CategoryLabel) {
		name := p.lexer.Current().Value
		p.lexer.Move()
		return intermediate.NewDataWithName(name, p.matchOptionalComment()), nil
	}
	value, err := p.matchNumber()
	if err != nil {
		return nil, err
	}
	if value < 0 {
		value = ^value + 1
	}
	return intermediate.NewDataWithValue(uint64(value), p.matchOptionalComment()), nil
}

func (p *parser) matchLabel() (string, error) {
	return p.lexer.Match(CategoryLabel)
}

func (p *parser) matchNumber() (int64, error) {
	v, err := p.lexer.Match(CategoryNumber)
	if err != nil {
		return 0, err
	}
	// TODO if v starts with '0x' or '-0x', parse as hexadecimal
	return strconv.ParseInt(v, 10, 64)
}

func (p *parser) matchOptionalComment() string {
	if p.lexer.CurrentIs(CategoryComment) {
		comment := strings.TrimPrefix(p.lexer.Current().Value, "// ")
		p.lexer.Move()
		return comment
	}
	return ""
}

func (p *parser) matchRegister() (flamego.Register, error) {
	r, err := p.lexer.Match(CategoryLowerName)
	if err != nil {
		return 0, err
	}
	ok, err := regexp.MatchString(`r[\\d]*`, r)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, &Error{p.lexer.Line(), fmt.Sprintf("Expected Register, found '%s'", r)}
	}
	v, err := strconv.ParseUint(strings.TrimPrefix(r, "r"), 10, 64)
	if err != nil {
		return 0, err
	}
	if v >= flamego.RegisterCount {
		return 0, &Error{p.lexer.Line(), fmt.Sprintf("Register Index Out of Bounds: '%d'", v)}
	}
	return flamego.Register(v), nil
}

func (p *parser) matchStatement() (intermediate.Addressable, error) {
	if p.lexer.CurrentIs(CategoryLabel) {
		name := p.lexer.Current().Value
		p.lexer.Move()
		label := intermediate.NewLabelWithName(name, p.matchOptionalComment())
		p.assembler.AddLabel(name, label)
		return label, nil
	} else if p.lexer.CurrentIs(CategoryData) {
		p.lexer.Move()
		return p.matchData()
	}
	value, err := p.lexer.Match(CategoryLowerName)
	if err != nil {
		return nil, err
	}
	switch value {
	case "halt":
		return intermediate.NewHalt(p.matchOptionalComment()), nil
	case "noop":
		return intermediate.NewNoop(p.matchOptionalComment()), nil
	case "sleep":
		return intermediate.NewSleep(p.matchOptionalComment()), nil
	case "signal":
		r, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		return intermediate.NewSignal(r, p.matchOptionalComment()), nil
	case "lock":
		return intermediate.NewLock(p.matchOptionalComment()), nil
	case "unlock":
		return intermediate.NewUnlock(p.matchOptionalComment()), nil
	case "interrupt":
		v, err := p.matchNumber()
		if err != nil {
			return nil, err
		}
		if v < 0 || v >= flamego.InterruptCount {
			return nil, &Error{p.lexer.Line(), fmt.Sprintf("Invalid Interrupt Value: '%d'", v)}
		}
		return intermediate.NewInterrupt(flamego.InterruptValue(v), p.matchOptionalComment()), nil
	case "uninterrupt":
		r, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		return intermediate.NewUninterrupt(r, p.matchOptionalComment()), nil
	case "jump":
		l, err := p.matchLabel()
		if err != nil {
			return nil, err
		}
		return intermediate.NewJump(isa.JumpEZ, l, flamego.R0, p.matchOptionalComment()), nil
	case "jez":
		r, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		l, err := p.matchLabel()
		if err != nil {
			return nil, err
		}
		return intermediate.NewJump(isa.JumpEZ, l, r, p.matchOptionalComment()), nil
	case "jnz":
		r, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		l, err := p.matchLabel()
		if err != nil {
			return nil, err
		}
		return intermediate.NewJump(isa.JumpNZ, l, r, p.matchOptionalComment()), nil
	case "jle":
		r, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		l, err := p.matchLabel()
		if err != nil {
			return nil, err
		}
		return intermediate.NewJump(isa.JumpLE, l, r, p.matchOptionalComment()), nil
	case "jlz":
		r, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		l, err := p.matchLabel()
		if err != nil {
			return nil, err
		}
		return intermediate.NewJump(isa.JumpLZ, l, r, p.matchOptionalComment()), nil
	case "call":
		a, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		return intermediate.NewCall(a, p.matchOptionalComment()), nil
	case "return":
		a, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		return intermediate.NewReturn(a, p.matchOptionalComment()), nil
	case "loadc":
		c := p.lexer.Current()
		switch c.Category {
		case CategoryLabel:
			p.lexer.Move()
			r, err := p.matchRegister()
			if err != nil {
				return nil, err
			}
			return intermediate.NewLoadCWithLabel(c.Value, r, p.matchOptionalComment()), nil
		case CategoryUpperName:
			p.lexer.Move()
			r, err := p.matchRegister()
			if err != nil {
				return nil, err
			}
			return intermediate.NewLoadCWithConstant(c.Value, r, p.matchOptionalComment()), nil
		default:
			c, err := p.matchNumber()
			if err != nil {
				return nil, err
			}
			if c < 0 {
				return nil, &Error{p.lexer.Line(), fmt.Sprintf("Negative LoadC Constant: %d", c)}
			}
			r, err := p.matchRegister()
			if err != nil {
				return nil, err
			}
			return intermediate.NewLoadCWithNumber(uint32(c), r, p.matchOptionalComment()), nil
		}
	case "load":
		a, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		c := p.lexer.Current()
		switch c.Category {
		case CategoryLabel:
			p.lexer.Move()
			r, err := p.matchRegister()
			if err != nil {
				return nil, err
			}
			return intermediate.NewLoadWithLabel(a, c.Value, r, p.matchOptionalComment()), nil
		case CategoryUpperName:
			p.lexer.Move()
			r, err := p.matchRegister()
			if err != nil {
				return nil, err
			}
			return intermediate.NewLoadWithConstant(a, c.Value, r, p.matchOptionalComment()), nil
		default:
			o, err := p.matchNumber()
			if err != nil {
				return nil, err
			}
			if o < 0 {
				return nil, &Error{p.lexer.Line(), fmt.Sprintf("Negative Load Offset: %d", o)}
			}
			d, err := p.matchRegister()
			if err != nil {
				return nil, err
			}
			return intermediate.NewLoadWithOffset(a, uint32(o), d, p.matchOptionalComment()), nil
		}
	case "store":
		a, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		c := p.lexer.Current()
		switch c.Category {
		case CategoryLabel:
			p.lexer.Move()
			r, err := p.matchRegister()
			if err != nil {
				return nil, err
			}
			return intermediate.NewStoreWithLabel(a, c.Value, r, p.matchOptionalComment()), nil
		case CategoryUpperName:
			p.lexer.Move()
			r, err := p.matchRegister()
			if err != nil {
				return nil, err
			}
			return intermediate.NewStoreWithConstant(a, c.Value, r, p.matchOptionalComment()), nil
		default:
			o, err := p.matchNumber()
			if err != nil {
				return nil, err
			}
			if o < 0 {
				return nil, &Error{p.lexer.Line(), fmt.Sprintf("Negative Store Offset: %d", o)}
			}
			s, err := p.matchRegister()
			if err != nil {
				return nil, err
			}
			return intermediate.NewStoreWithOffset(a, uint32(o), s, p.matchOptionalComment()), nil
		}
	case "clear":
		a, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		c := p.lexer.Current()
		switch c.Category {
		case CategoryLabel:
			p.lexer.Move()
			return intermediate.NewClearWithLabel(a, c.Value, p.matchOptionalComment()), nil
		case CategoryUpperName:
			p.lexer.Move()
			return intermediate.NewClearWithConstant(a, c.Value, p.matchOptionalComment()), nil
		default:
			o, err := p.matchNumber()
			if err != nil {
				return nil, err
			}
			if o < 0 {
				return nil, &Error{p.lexer.Line(), fmt.Sprintf("Negative Clear Offset: %d", o)}
			}
			return intermediate.NewClearWithOffset(a, uint32(o), p.matchOptionalComment()), nil
		}
	case "flush":
		a, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		c := p.lexer.Current()
		switch c.Category {
		case CategoryLabel:
			p.lexer.Move()
			return intermediate.NewFlushWithLabel(a, c.Value, p.matchOptionalComment()), nil
		case CategoryUpperName:
			p.lexer.Move()
			return intermediate.NewFlushWithConstant(a, c.Value, p.matchOptionalComment()), nil
		default:
			o, err := p.matchNumber()
			if err != nil {
				return nil, err
			}
			if o < 0 {
				return nil, &Error{p.lexer.Line(), fmt.Sprintf("Negative Flush Offset: %d", o)}
			}
			return intermediate.NewFlushWithOffset(a, uint32(o), p.matchOptionalComment()), nil
		}
	case "push":
		// TODO support register lists
		a, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		return intermediate.NewPush(a, p.matchOptionalComment()), nil
	case "pop":
		// TODO support register lists
		a, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		return intermediate.NewPop(a, p.matchOptionalComment()), nil
	case "not":
		s, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		d, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		return intermediate.NewNot(s, d, p.matchOptionalComment()), nil
	case "and":
		s1, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		s2, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		d, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		return intermediate.NewAnd(s1, s2, d, p.matchOptionalComment()), nil
	case "or":
		s1, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		s2, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		d, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		return intermediate.NewOr(s1, s2, d, p.matchOptionalComment()), nil
	case "xor":
		s1, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		s2, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		d, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		return intermediate.NewXor(s1, s2, d, p.matchOptionalComment()), nil
	case "leftshift":
		s1, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		s2, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		d, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		return intermediate.NewLeftShift(s1, s2, d, p.matchOptionalComment()), nil
	case "rightshift":
		s1, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		s2, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		d, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		return intermediate.NewRightShift(s1, s2, d, p.matchOptionalComment()), nil
	case "add":
		s1, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		s2, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		d, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		return intermediate.NewAdd(s1, s2, d, p.matchOptionalComment()), nil
	case "subtract":
		s1, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		s2, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		d, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		return intermediate.NewSubtract(s1, s2, d, p.matchOptionalComment()), nil
	case "multiply":
		s1, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		s2, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		d, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		return intermediate.NewMultiply(s1, s2, d, p.matchOptionalComment()), nil
	case "divide":
		s1, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		s2, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		d, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		return intermediate.NewDivide(s1, s2, d, p.matchOptionalComment()), nil
	case "modulo":
		s1, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		s2, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		d, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		return intermediate.NewModulo(s1, s2, d, p.matchOptionalComment()), nil
	case "copy":
		s, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		d, err := p.matchRegister()
		if err != nil {
			return nil, err
		}
		return intermediate.NewAdd(s, flamego.R0, d, p.matchOptionalComment()), nil
	}
	return nil, &Error{p.lexer.Line(), fmt.Sprintf("Unrecognized Instruction: %s", value)}
}
