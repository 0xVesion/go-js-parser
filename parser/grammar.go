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
// 	: AssignmentExpression
// 	;
func (p *parser) expression() interface{} {
	return p.assignmentExpression()
}

// AssignmentExpression
// 	: AdditiveExpression
// 	| LeftHandSideExpression ASSIGNMENT_OPERATOR AssignmentExpression
// 	;
func (p *parser) assignmentExpression() interface{} {
	left := p.additiveExpression()

	if p.lookAhead.Type != tokenizer.AssignmentOperator {
		return left
	}

	if !p.factory.IsIdentifier(left) {
		panic(fmt.Errorf("Invalid left-hand side expression: %v", left))
	}

	op := p.consume(tokenizer.AssignmentOperator)
	right := p.assignmentExpression()

	return p.factory.AssignmentExpression(op.Value, left, right)
}

// AdditiveExpression
// 	: MultiplicativeExpression
// 	| AdditiveExpression ADDITIVE_OPERATOR Literal
// 	;
func (p *parser) additiveExpression() interface{} {
	return p.binaryExpression(
		p.multiplicativeExpression,
		tokenizer.AdditiveOperator,
	)
}

// MultiplicativeExpression
// 	: PrimaryExpression
// 	| MultiplicativeExpression MULTIPLICATIVE_OPERATOR PrimaryExpression
// 	;
func (p *parser) multiplicativeExpression() interface{} {
	return p.binaryExpression(
		p.primaryExpression,
		tokenizer.MultiplicativeOperator,
	)
}

// PrimaryExpression
// 	: Literal
//  | ParenthesizedExpression
//  | LeftHandSideExpression
// 	;
func (p *parser) primaryExpression() interface{} {
	if p.isLookaheadLiteral() {
		return p.literal()
	}

	switch p.lookAhead.Type {
	case tokenizer.OpeningParenthesis:
		return p.parenthesizedExpression()
	default:
		return p.leftHandSideExpression()
	}
}

// LeftHandSideExpression
// 	: Identifier
// 	;
func (p *parser) leftHandSideExpression() interface{} {
	return p.identifier()
}

// Identifier
// 	: IDENTIFIER
// 	;
func (p *parser) identifier() interface{} {
	return p.factory.Identifier(p.consume(tokenizer.Identifier).Value)
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
// 	: NUMBER
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
// 	: STRING
// 	;
func (p *parser) stringLiteral() interface{} {
	token := p.consume(tokenizer.String)

	return p.factory.Literal(token.Value[1 : len(token.Value)-1])
}
