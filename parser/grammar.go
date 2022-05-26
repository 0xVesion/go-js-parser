package parser

import (
	"fmt"
	"strconv"

	"github.com/0xvesion/go-parser/tokenizer"
)

func (p *parser) program() (Node, error) {
	node, err := p.literal()
	if err != nil {
		return Node{}, err
	}

	return Node{Program, node}, nil
}

func (p *parser) literal() (Node, error) {
	switch p.lookAhead.Type {
	case tokenizer.Number:
		return p.numericLiteral()
	case tokenizer.String:
		return p.stringLiteral()

	}

	return Node{}, fmt.Errorf("invalid literal type %v", p.lookAhead.Type)
}

func (p *parser) numericLiteral() (Node, error) {
	token, err := p.consume(tokenizer.Number)
	if err != nil {
		return Node{}, err
	}

	value, err := strconv.Atoi(token.Value)
	if err != nil {
		return Node{}, err
	}

	return Node{
		NumericLiteral,
		value,
	}, nil
}

func (p *parser) stringLiteral() (Node, error) {
	token, err := p.consume(tokenizer.String)

	if err != nil {
		return Node{}, err
	}

	return Node{
		StringLiteral,
		token.Value[1 : len(token.Value)-1],
	}, nil
}
