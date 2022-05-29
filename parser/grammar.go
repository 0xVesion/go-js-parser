package parser

import (
	"fmt"
	"strconv"

	"github.com/0xvesion/go-js-parser/tokenizer"
)

// Entry point of the program.
// Program
// 	: StatementList
// 	;
func (p *parser) program() interface{} {
	return p.factory.Program(p.statementList(tokenizer.None)...)
}

// StatementList
// 	: Statement
// 	| StatementList Statement
// 	;
func (p *parser) statementList(endLookahead tokenizer.Type) []interface{} {
	sl := []interface{}{}

	for p.lookAhead.Type != endLookahead {
		sl = append(sl, p.statement())
	}

	return sl
}

// Statement
// 	: ExpressionStatment
// 	| BlockStatement
// 	| EmptyStatement
// 	;
func (p *parser) statement() interface{} {
	switch p.lookAhead.Type {
	case tokenizer.OpeningCurlyBrace:
		return p.blockStatement()
	case tokenizer.Semicolon:
		return p.emptyStatement()
	default:
		return p.expressionStatment()
	}
}

// EmptyStatement
// 	: ';'
// 	;
func (p *parser) emptyStatement() interface{} {
	p.consume(tokenizer.Semicolon)

	return p.factory.EmptyStatement()
}

// BlockStatement
// 	: '{' StatementList '}'
// 	;
func (p *parser) blockStatement() interface{} {
	p.consume(tokenizer.OpeningCurlyBrace)

	sl := p.statementList(tokenizer.ClosingCurlyBrace)

	p.consume(tokenizer.ClosingCurlyBrace)

	return p.factory.BlockStatement(sl...)
}

// ExpressionStatment
// 	: Expression ';'
// 	;
func (p *parser) expressionStatment() interface{} {
	e := p.expression()

	p.consume(tokenizer.Semicolon)

	return p.factory.ExpressionStatement(e)
}

// Expression
// 	: AdditiveExpression
// 	;
func (p *parser) expression() interface{} {
	return p.additiveExpression()
}

// AdditiveExpression
// 	: MultiplicativeExpression
// 	| AdditiveExpression AdditiveOperator Literal
// 	;
func (p *parser) additiveExpression() interface{} {
	left := p.multiplicativeExpression()

	for p.lookAhead.Type == tokenizer.AdditiveOperator {
		operator := p.consume(tokenizer.AdditiveOperator)
		right := p.multiplicativeExpression()

		left = p.factory.BinaryExpression(operator.Value, left, right)
	}

	return left
}

// MultiplicativeExpression
// 	: PrimaryExpression
// 	| MultiplicativeExpression MultiplicativeOperator PrimaryExpression
// 	;
func (p *parser) multiplicativeExpression() interface{} {
	left := p.primaryExpression()

	for p.lookAhead.Type == tokenizer.MultiplicativeOperator {
		operator := p.consume(tokenizer.MultiplicativeOperator)
		right := p.primaryExpression()

		left = p.factory.BinaryExpression(operator.Value, left, right)
	}

	return left
}

// PrimaryExpression
// 	: Literal
//  | ParenthesizedExpression
// 	;
func (p *parser) primaryExpression() interface{} {
	switch p.lookAhead.Type {
	case tokenizer.OpeningParenthesis:
		return p.parenthesizedExpression()
	default:
		return p.literal()
	}
}

// ParenthesizedExpression
// 	: '(' Expression ')'
// 	;
func (p *parser) parenthesizedExpression() interface{} {
	p.consume(tokenizer.OpeningParenthesis)

	ex := p.expression()

	p.consume(tokenizer.ClosingParenthesis)

	return ex
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
func (p *parser) numericLiteral() interface{} {
	token := p.consume(tokenizer.Number)

	value, err := strconv.Atoi(token.Value)
	if err != nil {
		panic(err)
	}

	return p.factory.Literal(value)
}

// StringLiteral
// 	: String
// 	;
func (p *parser) stringLiteral() interface{} {
	token := p.consume(tokenizer.String)

	return p.factory.Literal(token.Value[1 : len(token.Value)-1])
}
