package parser

import (
	"fmt"

	"github.com/0xvesion/go-parser/tokenizer"
)

type Type string

const (
	Program             Type = "Program"
	NumericLiteral           = "NumericLiteral"
	StringLiteral            = "StringLiteral"
	ExpressionStatement      = "ExpressionStatement"
	BlockStatement           = "BlockStatement"
)

type Parser interface {
	Parse() (program, error)
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

func (p *parser) Parse() (n program, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	n = p.program()

	return
}

func (p *parser) consume(t tokenizer.Type) tokenizer.Token {
	token := p.lookAhead

	if token.Type != t {
		panic(fmt.Errorf("Unexpected token type. want: %s got: %s", t, token.Type))
	}

	lookAhead, err := p.t.Next()
	if err != nil {
		panic(err)
	}
	p.lookAhead = lookAhead

	return token
}
