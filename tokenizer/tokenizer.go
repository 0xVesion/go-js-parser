package tokenizer

import (
	"fmt"
	"regexp"
)

type Token struct {
	Type
	Value string
}

type Tokenizer struct {
	src string
}

func New(src string) *Tokenizer {
	return &Tokenizer{
		src: src,
	}
}

func (t *Tokenizer) All() ([]Token, error) {
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

func (t *Tokenizer) HasNext() bool {
	return len(t.src) > 0
}

func (t *Tokenizer) Next() (Token, error) {
	for k, v := range grammar {
		r, _ := regexp.Compile(v)
		loc := r.FindStringIndex(t.src)
		if loc == nil || loc[0] != 0 {
			continue
		}

		match := t.src[loc[0]:loc[1]]
		t.src = t.src[loc[1]:]

		return Token{k, match}, nil
	}

	return Token{}, fmt.Errorf("Unknown token")
}
