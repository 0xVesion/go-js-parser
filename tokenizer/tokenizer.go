package tokenizer

import (
	"fmt"
	"regexp"
)

type Token struct {
	Type
	Value string
}

type tokenizer struct {
	src    string
	cursor int
}

type Tokenizer interface {
	HasNext() bool
	Next() (Token, error)
	All() ([]Token, error)
}

func New(src string) Tokenizer {
	return &tokenizer{
		src: src,
	}
}

func (t *tokenizer) All() ([]Token, error) {
	tokens := []Token{}
	for t.HasNext() {
		token, err := t.Next()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)

	}

	return tokens, nil
}

func (t *tokenizer) HasNext() bool {
	return t.cursor < len(t.src)
}

func (t *tokenizer) Next() (Token, error) {
	if !t.HasNext() {
		return Token{}, nil
	}

	s := t.src[t.cursor:]
	for _, current := range grammar {
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

			return Token{current.Type, match}, nil

		}
	}

	return Token{}, fmt.Errorf("Unknown token")
}
