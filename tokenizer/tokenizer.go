package tokenizer

import (
	"fmt"
	"regexp"
)

type Token struct {
	Type
	Value string
	Start int
	End   int
}

func (to Token) Is(types ...Type) bool {
	for _, t := range types {
		if to.Type == t {
			return true
		}
	}

	return false
}

func (to Token) Not(types ...Type) bool {
	return !to.Is(types...)
}

type tokenizer struct {
	src    string
	cursor int
}

type Tokenizer interface {
	HasNext() bool
	Next() (Token, error)
	Src() string
	Cursor() int
}

func New(src string) Tokenizer {
	return &tokenizer{
		src: src,
	}
}

func (t *tokenizer) Cursor() int {
	return t.cursor
}

func (t *tokenizer) Src() string {
	return t.src
}

func (t *tokenizer) HasNext() bool {
	return t.cursor < len(t.src)
}

func (t *tokenizer) Next() (Token, error) {
	if !t.HasNext() {
		return Token{}, nil
	}

	s := t.src[t.cursor:]
	tokenStart := t.cursor
	for _, current := range spec {
		for _, expr := range current.Regexp {
			r, err := regexp.Compile(expr)
			if err != nil {
				return Token{}, nil
			}

			loc := r.FindStringIndex(s)
			if loc == nil || loc[0] != 0 {
				continue
			}

			match := s[loc[0]:loc[1]]
			t.cursor += loc[1]

			if current.Type == None {
				return t.Next()
			}

			return Token{current.Type, match, tokenStart, t.cursor}, nil
		}
	}

	return Token{}, fmt.Errorf("unknown token: %s", s)
}
