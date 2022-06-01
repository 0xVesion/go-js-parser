package tokenizer

import (
	"fmt"
	"regexp"
)

type Token struct {
	Type
	Value string
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
}

func New(src string) Tokenizer {
	return &tokenizer{
		src: src,
	}
}

func (t *tokenizer) HasNext() bool {
	return t.cursor < len(t.src)
}

func (t *tokenizer) Next() (Token, error) {
	if !t.HasNext() {
		return Token{}, nil
	}

	s := t.src[t.cursor:]
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

			return Token{current.Type, match}, nil
		}
	}

	return Token{}, fmt.Errorf("Unknown token: %s", s)
}
