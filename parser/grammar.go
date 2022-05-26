package parser

import (
	"fmt"
	"strconv"

	"github.com/0xvesion/go-parser/tokenizer"
)

// Entry point of the program.
// Program
// 	: Literal
// 	;
func (p *parser) program() (Node, error) {
	node, err := p.literal()
	if err != nil {
		return Node{}, err
	}

	return Node{Program, node}, nil
}

// Literal
// 	: NumericLiteral
// 	| StringLiteral
// 	;
func (p *parser) literal() (Node, error) {
	switch p.lookAhead.Type {
	case tokenizer.Number:
		return p.numericLiteral()
	case tokenizer.String:
		return p.stringLiteral()

	}

	return Node{}, fmt.Errorf("invalid literal type %v", p.lookAhead.Type)
}

// NumericLiteral
// 	: Number
// 	;
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

// StringLiteral
// 	: String
// 	;
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
