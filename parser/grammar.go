package parser

import (
	"fmt"
	"strconv"

	"github.com/0xvesion/go-parser/tokenizer"
)

// Entry point of the program.
// Program
// 	: StatementList
// 	;
func (p *parser) program() (Node, error) {
	node, err := p.statementList()
	if err != nil {
		return Node{}, err
	}

	return Node{Program, node}, nil
}

// StatementList
// 	: Statement
// 	| StatementList Statement
// 	;
func (p *parser) statementList() ([]Node, error) {
	sl := []Node{}

	for p.lookAhead.Type != tokenizer.None {
		statement, _ := p.statement()
		sl = append(sl, statement)
	}

	return sl, nil
}

// Statement
// 	: ExpressionStatment
// 	;
func (p *parser) statement() (Node, error) {
	return p.expressionStatment()
}

// ExpressionStatment
// 	: Expression ';'
// 	;
func (p *parser) expressionStatment() (Node, error) {
	e, _ := p.expression()

	p.consume(tokenizer.Semicolon)

	return Node{ExpressionStatement, e}, nil

}

// Expression
// 	: Literal
// 	;
func (p *parser) expression() (Node, error) {
	return p.literal()
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
