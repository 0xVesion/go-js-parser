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
	case tokenizer.VariableDeclarationKeyword:
		return p.variableDeclaration()
	default:
		return p.expressionStatment()
	}
}

// VariableDeclaration
// 	: VARIABLE_DECLARATION_KEYWORD VariableDeclaratorList ';'
// 	;
func (p *parser) variableDeclaration() interface{} {
	kind := p.consume(tokenizer.VariableDeclarationKeyword)
	declarations := p.variableDeclaratorList()
	p.consume(tokenizer.Semicolon)

	return p.factory.VariableDeclaration(kind.Value, declarations)
}

// VariableDeclaratorList
// 	: VariableDeclarator
//	| VariableDeclaratorList ',' VariableDeclarator
// 	;
func (p *parser) variableDeclaratorList() []interface{} {
	declarations := []interface{}{p.variableDeclarator()}

	for p.lookAhead.Type == tokenizer.Comma {
		p.consume(tokenizer.Comma)
		declarations = append(declarations, p.variableDeclarator())
	}

	return declarations
}

// VariableDeclarator
// 	: Identifier OptVariableInitializer
// 	;
func (p *parser) variableDeclarator() interface{} {
	id := p.identifier()

	var init interface{}
	if p.lookAhead.Type != tokenizer.Semicolon && p.lookAhead.Type != tokenizer.Comma {

		p.consume(tokenizer.SimpleAssignmentOperator)
		init = p.assignmentExpression()
	}

	return p.factory.VariableDeclarator(id, init)
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
// 	: RelationalExpression
// 	| LeftHandSideExpression ASSIGNMENT_OPERATOR AssignmentExpression
// 	;
func (p *parser) assignmentExpression() interface{} {
	left := p.relationalExpression()

	if !p.isLookaheadAssignmentOperator() {
		return left
	}

	if !p.factory.IsIdentifier(left) {
		panic(fmt.Errorf("Invalid left-hand side expression: %v", left))
	}

	op := p.consume(p.lookAhead.Type).Value
	right := p.assignmentExpression()

	return p.factory.AssignmentExpression(op, left, right)
}

// RelationalExpression
// 	: AdditiveExpression
// 	| RelationalExpression RELATIONAL_OPERATOR AdditiveExpression
// 	;
func (p *parser) relationalExpression() interface{} {
	return p.binaryExpression(
		p.additiveExpression,
		tokenizer.RelationalOperator,
	)
}

// AdditiveExpression
// 	: MultiplicativeExpression
// 	| AdditiveExpression ADDITIVE_OPERATOR MultiplicativeExpression
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
