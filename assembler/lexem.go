package assembler

import (
	"fmt"
	"regexp"
)

type Lexem interface {
	Category() Category
	Matches(string) bool
}

func NewLexemLiteral(c Category, v string) Lexem {
	return &lexemLiteral{
		category: c,
		value:    v,
	}
}

type lexemLiteral struct {
	category Category
	value    string
}

func (l *lexemLiteral) Category() Category {
	return l.category
}

func (l *lexemLiteral) Matches(s string) bool {
	return l.value == s
}

func (l *lexemLiteral) String() string {
	return fmt.Sprintf("%d : %s", l.category, l.value)
}

func NewLexemRegexp(c Category, r string) Lexem {
	return &lexemRegex{
		category: c,
		regexp:   regexp.MustCompile(r),
	}
}

type lexemRegex struct {
	category Category
	regexp   *regexp.Regexp
}

func (l *lexemRegex) Category() Category {
	return l.category
}

func (l *lexemRegex) Matches(s string) bool {
	return l.regexp.MatchString(s)
}

func (l *lexemRegex) String() string {
	return fmt.Sprintf("%d : %s", l.category, l.regexp)
}
