package assembler

import (
	"fmt"
)

type Token struct {
	Category Category
	Value    string
}

func NewToken(c Category) *Token {
	return &Token{
		Category: c,
	}
}

func NewTokenWithValue(c Category, v string) *Token {
	return &Token{
		Category: c,
		Value:    v,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("%s : %s", t.Category, t.Value)
}
