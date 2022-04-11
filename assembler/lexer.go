package assembler

import (
	"fmt"
	"io"
)

type Lexer interface {
	Read(io.Reader) (int64, error)
	Line() int
	Current() *Token
	Move()
	CurrentIs(Category) bool
	Match(Category) (string, error)
}

func NewLexer() Lexer {
	l := &lexer{
		lexems: []Lexem{
			NewLexemLiteral(CategoryNewLine, "\n"),
			NewLexemRegexp(CategoryWhitespace, "^\\s$"),
			NewLexemLiteral(CategoryColon, ":"),
			NewLexemLiteral(CategoryFullStop, "."),
			NewLexemRegexp(CategoryComment, "^//.*\n?$"),
			NewLexemLiteral(CategoryAlign, "align"),
			NewLexemLiteral(CategoryAllocate, "allocate"),
			NewLexemLiteral(CategoryData, "data"),
			NewLexemLiteral(CategoryComma, ","),
			NewLexemRegexp(CategoryLabel, "^#[\\.:_a-zA-Z0-9\\(\\)]+$"),
			NewLexemRegexp(CategoryNumber, "^0$|^[+-]?[1-9][0-9]*$|^0x[0-9A-Fa-f]*$"),
			NewLexemRegexp(CategoryUpperName, "^[A-Z][_a-zA-Z0-9]*$"),
			NewLexemRegexp(CategoryLowerName, "^[_a-z][_a-zA-Z0-9]*$"),
		},
		line: 1,
	}
	return l
}

type lexer struct {
	lexems  []Lexem
	tokens  []*Token
	line    int
	index   int
	current *Token
}

func (l *lexer) Read(r io.Reader) (int64, error) {
	var (
		count     int64
		input     byte
		character string
		current   string
		next      string
	)
	buffer := make([]byte, 1)
	for {
		i, err := r.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}
		if i != 1 {
			break
		}
		input = buffer[0]
		character = string(input)
		next = current + character
		tc := l.token(current)
		tn := l.token(next)
		if tc != nil && tn == nil {
			l.tokens = append(l.tokens, tc)
			current = character
		} else {
			current = next
		}
	}
	if current != "" {
		t := l.token(current)
		l.tokens = append(l.tokens, t)
	}
	l.tokens = append(l.tokens, NewToken(CategoryEOF))
	return count, nil
}

func (l *lexer) Line() int {
	return l.line
}

func (l *lexer) Current() *Token {
	return l.current
}

func (l *lexer) Move() {
	l.current = l.next()
}

func (l *lexer) CurrentIs(c Category) bool {
	return l.current.Category == c
}

func (l *lexer) Match(c Category) (string, error) {
	if !l.CurrentIs(c) {
		return "", &Error{l.line, fmt.Sprintf("Expected '%s', found '%s'", c, l.current)}
	}
	v := l.current.Value
	l.Move()
	return v, nil
}

func (l *lexer) next() *Token {
	if l.index < len(l.tokens) {
		t := l.tokens[l.index]
		l.index++
		switch t.Category {
		case CategoryComment, CategoryNewLine:
			l.line++
			fallthrough
		case CategoryWhitespace:
			t = l.next()
		}
		return t
	}
	return nil
}

func (l *lexer) token(value string) *Token {
	for _, x := range l.lexems {
		if x.Matches(value) {
			return NewTokenWithValue(x.Category(), value)
		}
	}
	return nil
}
