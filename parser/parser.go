package parser

import (
	"fmt"

	"github.com/0xvesion/go-parser/tokenizer"
)

type Type string

const (
	Program        Type = "Program"
	NumericLiteral Type = "NumericLiteral"
	StringLiteral  Type = "StringLiteral"
)

type Node struct {
	Type
	Value interface{}
}

type Parser interface {
	Parse() (Node, error)
}

type parser struct {
	t         tokenizer.Tokenizer
	lookAhead tokenizer.Token
}

func New(t tokenizer.Tokenizer) Parser {
	lookAhead, err := t.Next()
	if err != nil {
		panic(err)
	}

	return &parser{
		t:         t,
		lookAhead: lookAhead,
	}
}

func (p *parser) Parse() (Node, error) {
	return p.program()
}

func (p *parser) consume(t tokenizer.Type) (tokenizer.Token, error) {
	token := p.lookAhead

	if token.Type != t {
		return tokenizer.Token{}, fmt.Errorf("Unexpected token type %s", token.Type)
	}

	lookAhead, err := p.t.Next()
	if err != nil {
		return tokenizer.Token{}, err
	}
	p.lookAhead = lookAhead

	return token, err
}
