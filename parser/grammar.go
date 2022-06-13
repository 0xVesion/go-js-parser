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
	sl := p.statementList(tokenizer.None)
	sl = p.addDirectives(sl)

	return NewProgram(0, len(p.t.Src()), sl...)
}

// StatementList
// 	: Statement
// 	| StatementList Statement
// 	;
func (p *parser) statementList(endLookahead tokenizer.Type) []interface{} {
	sl := []interface{}{}

	for p.lookAhead.Not(endLookahead) {
		sl = append(sl, p.statement())
	}

	return sl
}

// Statement
// 	: ExpressionStatment
// 	| BlockStatement
// 	| EmptyStatement
// 	| VariableDeclaration
// 	| IfStatement
// 	;
func (p *parser) statement() interface{} {
	switch p.lookAhead.Type {
	case tokenizer.OpeningCurlyBrace:
		return p.blockStatement()
	case tokenizer.Semicolon:
		return p.emptyStatement()
	case tokenizer.VariableDeclarationKeyword:
		return p.variableDeclaration()
	case tokenizer.IfKeyword:
		return p.ifStatement()
	default:
		return p.expressionStatment()
	}
}

// IfStatement
// 	: 'if' ParenthesizedExpression Statement
// 	| 'if' ParenthesizedExpression Statement 'else' Statement
// 	;
func (p *parser) ifStatement() interface{} {
	p.consume(tokenizer.IfKeyword)

	test := p.parenthesizedExpression()

	consequent := p.statement()

	if p.lookAhead.Not(tokenizer.ElseKeyword) {
		return NewIfStatement(test, consequent, nil)
	}

	p.consume(tokenizer.ElseKeyword)

	alternate := p.statement()

	return NewIfStatement(test, consequent, alternate)
}

// VariableDeclaration
// 	: VARIABLE_DECLARATION_KEYWORD VariableDeclaratorList ';'
// 	;
func (p *parser) variableDeclaration() interface{} {
	kind := p.consume(tokenizer.VariableDeclarationKeyword)
	declarations := p.variableDeclaratorList()
	p.consume(tokenizer.Semicolon)

	return NewVariableDeclaration(kind.Value, declarations)
}

// VariableDeclaratorList
// 	: VariableDeclarator
//	| VariableDeclaratorList ',' VariableDeclarator
// 	;
func (p *parser) variableDeclaratorList() []interface{} {
	declarations := []interface{}{p.variableDeclarator()}

	for p.lookAhead.Is(tokenizer.Comma) {
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
	if p.lookAhead.Not(tokenizer.Semicolon, tokenizer.Comma) {
		p.consume(tokenizer.SimpleAssignmentOperator)
		init = p.assignmentExpression()
	}

	return NewVariableDeclarator(id, init)
}

// EmptyStatement
// 	: ';'
// 	;
func (p *parser) emptyStatement() interface{} {
	t := p.consume(tokenizer.Semicolon)

	return NewEmptyStatement(t.Start, t.End)
}

// BlockStatement
// 	: '{' StatementList '}'
// 	;
func (p *parser) blockStatement() interface{} {
	start := p.consume(tokenizer.OpeningCurlyBrace).Start

	sl := p.statementList(tokenizer.ClosingCurlyBrace)

	end := p.consume(tokenizer.ClosingCurlyBrace).End

	return NewBlockStatement(start, end, sl...)
}

// ExpressionStatment
// 	: Expression ';'
// 	;
func (p *parser) expressionStatment() interface{} {
	start := p.lookAhead.Start
	exp := p.expression()

	semi := p.consume(tokenizer.Semicolon)

	return NewExpressionStatement(start, semi.End, exp)
}

// Expression
// 	: AssignmentExpression
// 	;
func (p *parser) expression() interface{} {
	return p.assignmentExpression()
}

// AssignmentExpression
// 	: LogicalOrExpression
// 	| LeftHandSideExpression ASSIGNMENT_OPERATOR AssignmentExpression
// 	;
func (p *parser) assignmentExpression() interface{} {
	left := p.logicalOrExpression()

	if !p.isLookaheadAssignmentOperator() {
		return left
	}

	if _, isIdentifier := left.(IdentifierNode); !isIdentifier {
		panic(fmt.Errorf("invalid left-hand side expression: %v", left))
	}

	op := p.consumeAny().Value
	right := p.assignmentExpression()

	return NewAssignmentExpression(op, left, right)
}

// LogicalOrExpression
// 	: LogicalAndExpression
// 	| LogicalOrExpression '||' LogicalAndExpression
// 	;
func (p *parser) logicalOrExpression() interface{} {
	return p.logicalExpression(
		p.logicalAndExpression,
		tokenizer.LogicalOrOperator,
	)
}

// LogicalAndExpression
// 	: EqualityExpression
// 	| LogicalAndExpression '&&' EqualityExpression
// 	;
func (p *parser) logicalAndExpression() interface{} {
	return p.logicalExpression(
		p.equalityExpression,
		tokenizer.LogicalAndOperator,
	)
}

// EqualityExpression
// 	: RelationalExpression
// 	| EqualityExpression EQUALITY_OPERATOR RelationalExpression
// 	;
func (p *parser) equalityExpression() interface{} {
	return p.binaryExpression(
		p.relationalExpression,
		tokenizer.EqualityOperator,
	)
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
// 	: UnaryExpression
// 	| MultiplicativeExpression MULTIPLICATIVE_OPERATOR UnaryExpression
// 	;
func (p *parser) multiplicativeExpression() interface{} {
	return p.binaryExpression(
		p.unaryExpression,
		tokenizer.MultiplicativeOperator,
	)
}

// UnaryExpression
// 	: LeftHandSideExpression
// 	| ADDITIVE_OPERATOR UnaryExpression
//	| LOGICAL_NOT UnaryExpression
// 	;
func (p *parser) unaryExpression() interface{} {
	if p.lookAhead.Not(tokenizer.LogicalNotOperator, tokenizer.AdditiveOperator) {
		return p.leftHandSideExpression()
	}

	return NewUnaryExpression(
		p.consumeAny().Value,
		p.unaryExpression(),
	)
}

// LeftHandSideExpression
// 	: PrimaryExpression
//	;
func (p *parser) leftHandSideExpression() interface{} {
	return p.primaryExpression()
}

// PrimaryExpression
// 	: Literal
//  | ParenthesizedExpression
//  | Identifier
// 	;
func (p *parser) primaryExpression() interface{} {
	if p.isLookaheadLiteral() {
		return p.literal()
	}

	switch p.lookAhead.Type {
	case tokenizer.OpeningParenthesis:
		return p.parenthesizedExpression()
	case tokenizer.Identifier:
		return p.identifier()
	default:
		panic(fmt.Errorf("invalid token: %s", p.lookAhead.Type))
	}
}

// Identifier
// 	: IDENTIFIER
// 	;
func (p *parser) identifier() interface{} {
	return NewIdentifier(p.consume(tokenizer.Identifier).Value)
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
//	| BooleanLiteral
//  | NullLiteral
// 	;
func (p *parser) literal() interface{} {
	switch p.lookAhead.Type {
	case tokenizer.Number:
		return p.numericLiteral()
	case tokenizer.String:
		return p.stringLiteral()
	case tokenizer.BooleanLiteral:
		return p.booleanLiteral()
	case tokenizer.NullLiteral:
		return p.nullLiteral()
	}

	panic(fmt.Errorf("invalid literal type %v", p.lookAhead.Type))
}

// BooleanLiteral
// 	: 'true'
// 	| 'false'
// 	;
func (p *parser) booleanLiteral() interface{} {
	token := p.consume(tokenizer.BooleanLiteral)

	return NewLiteral(token.Start, token.End, token.Value == "true", token.Value)
}

// NullLiteral
// 	: 'null'
// 	;
func (p *parser) nullLiteral() interface{} {
	token := p.consume(tokenizer.NullLiteral)

	return NewLiteral(token.Start, token.End, nil, token.Value)
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

	return NewLiteral(token.Start, token.End, value, token.Value)
}

// StringLiteral
// 	: STRING
// 	;
func (p *parser) stringLiteral() interface{} {
	token := p.consume(tokenizer.String)

	return NewLiteral(token.Start, token.End, token.Value[1:len(token.Value)-1], token.Value)
}
