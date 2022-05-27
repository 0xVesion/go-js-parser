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
func (p *parser) program() Node {
	node := p.statementList(tokenizer.None)

	return Node{Program, node}
}

// StatementList
// 	: Statement
// 	| StatementList Statement
// 	;
func (p *parser) statementList(endLookahead tokenizer.Type) []Node {
	sl := []Node{}

	for p.lookAhead.Type != endLookahead {
		statement := p.statement()
		sl = append(sl, statement)
	}

	return sl
}

// Statement
// 	: ExpressionStatment
// 	| BlockStatement
// 	;
func (p *parser) statement() Node {
	switch p.lookAhead.Type {
	case tokenizer.OpeningCurlyBrace:
		return p.blockStatement()
	default:
		return p.expressionStatment()
	}
}

// BlockStatement
// 	: '{' StatementList '}'
// 	;
func (p *parser) blockStatement() Node {
	p.consume(tokenizer.OpeningCurlyBrace)

	sl := p.statementList(tokenizer.ClosingCurlyBrace)

	p.consume(tokenizer.ClosingCurlyBrace)

	return Node{BlockStatement, sl}
}

// ExpressionStatment
// 	: Expression ';'
// 	;
func (p *parser) expressionStatment() Node {
	e := p.expression()

	p.consume(tokenizer.Semicolon)

	return Node{ExpressionStatement, e}

}

// Expression
// 	: Literal
// 	;
func (p *parser) expression() Node {
	return p.literal()
}

// Literal
// 	: NumericLiteral
// 	| StringLiteral
// 	;
func (p *parser) literal() Node {
	switch p.lookAhead.Type {
	case tokenizer.Number:
		return p.numericLiteral()
	case tokenizer.String:
		return p.stringLiteral()
	}

	panic(fmt.Errorf("invalid literal type %v", p.lookAhead.Type))
}

// NumericLiteral
// 	: Number
// 	;
func (p *parser) numericLiteral() Node {
	token := p.consume(tokenizer.Number)

	value, err := strconv.Atoi(token.Value)
	if err != nil {
		panic(err)
	}

	return Node{
		NumericLiteral,
		value,
	}
}

// StringLiteral
// 	: String
// 	;
func (p *parser) stringLiteral() Node {
	token := p.consume(tokenizer.String)

	return Node{
		StringLiteral,
		token.Value[1 : len(token.Value)-1],
	}
}
