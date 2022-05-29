package parser

import (
	"fmt"

	"github.com/0xvesion/go-js-parser/tokenizer"
)

type Parser interface {
	Parse() (interface{}, error)
}

type parser struct {
	t         tokenizer.Tokenizer
	factory   AstFactory
	lookAhead tokenizer.Token
}

func New(t tokenizer.Tokenizer, factory AstFactory) Parser {
	lookAhead, err := t.Next()
	if err != nil {
		panic(err)
	}

	return &parser{
		t:         t,
		lookAhead: lookAhead,
		factory:   factory,
	}
}

func (p *parser) Parse() (n interface{}, err error) {
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
