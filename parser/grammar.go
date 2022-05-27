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
func (p *parser) program() program {
	sl := p.statementList(tokenizer.None)

	return newProgram(sl...)
}

// StatementList
// 	: Statement
// 	| StatementList Statement
// 	;
func (p *parser) statementList(endLookahead tokenizer.Type) []interface{} {
	sl := []interface{}{}

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
func (p *parser) statement() interface{} {
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
func (p *parser) blockStatement() blockStatement {
	p.consume(tokenizer.OpeningCurlyBrace)

	sl := p.statementList(tokenizer.ClosingCurlyBrace)

	p.consume(tokenizer.ClosingCurlyBrace)

	return newBlockStatement(sl...)
}

// ExpressionStatment
// 	: Expression ';'
// 	;
func (p *parser) expressionStatment() expressionStatement {
	e := p.expression()

	p.consume(tokenizer.Semicolon)

	return newExpressionStatement(e)
}

// Expression
// 	: Literal
// 	;
func (p *parser) expression() interface{} {
	return p.literal()
}

// Literal
// 	: NumericLiteral
// 	| StringLiteral
// 	;
func (p *parser) literal() interface{} {
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
func (p *parser) numericLiteral() literal[int] {
	token := p.consume(tokenizer.Number)

	value, err := strconv.Atoi(token.Value)
	if err != nil {
		panic(err)
	}

	return newNumericLiteral(value)
}

// StringLiteral
// 	: String
// 	;
func (p *parser) stringLiteral() literal[string] {
	token := p.consume(tokenizer.String)

	return newStringLiteral(token.Value[1 : len(token.Value)-1])
}
