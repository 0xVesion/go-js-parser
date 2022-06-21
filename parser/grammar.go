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
func (p *parser) program() Node {
	sl := p.statementList(tokenizer.None)
	sl = p.addDirectives(sl)

	return NewProgram(0, len(p.t.Src()), sl...)
}

// StatementList
// 	: Statement
// 	| StatementList Statement
// 	;
func (p *parser) statementList(endLookahead tokenizer.Type) []Node {
	sl := []Node{}

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
// 	| Iterationtatement
// 	| FunctionDeclaration
// 	| ReturnStatement
// 	| ClassDeclaration
// 	;
func (p *parser) statement() Node {
	switch p.lookAhead.Type {
	case tokenizer.OpeningCurlyBrace:
		return p.blockStatement()
	case tokenizer.Semicolon:
		return p.emptyStatement()
	case tokenizer.VariableDeclarationKeyword:
		return p.variableDeclaration()
	case tokenizer.IfKeyword:
		return p.ifStatement()
	case tokenizer.WhileKeyword:
		fallthrough
	case tokenizer.DoKeyword:
		fallthrough
	case tokenizer.ForKeyword:
		return p.iterationStatement()
	case tokenizer.FunctionKeyword:
		return p.functionDeclaration()
	case tokenizer.ReturnKeyword:
		return p.returnStatement()
	case tokenizer.ClassKeyword:
		return p.classDeclaration()
	default:
		return p.expressionStatment()
	}
}

// ClassDeclaration
// 	: 'class' Identifier ClassBody
// 	| 'class' Identifier 'extends' Identifier ClassBody
// 	;
func (p *parser) classDeclaration() Node {
	start := p.consume(tokenizer.ClassKeyword).Start

	id := p.identifier()
	var superClass Node
	if p.lookAhead.Is(tokenizer.ExtendsKeyword) {
		p.consume(tokenizer.ExtendsKeyword)
		superClass = p.identifier()
	}

	body := p.classBody()

	return NewClassDeclaration(start, body.End(), id, superClass, body)
}

// ClassBody
// 	: '{' OptClassMemberDefinitionList '}'
// 	;
func (p *parser) classBody() Node {
	start := p.consume(tokenizer.OpeningCurlyBrace).Start

	body := []Node{}
	for p.lookAhead.Not(tokenizer.ClosingCurlyBrace) {
		body = append(body, p.classMemberDefinition())
	}

	end := p.consume(tokenizer.ClosingCurlyBrace).End

	return NewClassBody(start, end, body)
}

// ClassMemberDefinition
// 	: NewPropertyDefinition
// 	| NewMethodDefinition
// 	;
func (p *parser) classMemberDefinition() Node {
	key := p.identifier()

	if p.lookAhead.Is(tokenizer.OpeningParenthesis) {
		value := p.functionExpression()

		kind := Method
		if IdentifierNode(key).Name() == string(Constructor) {
			kind = Constructor
		}

		return NewMethodDefinition(key.Start(), value.End(), key, kind, value)
	}

	var value Node
	if p.lookAhead.Is(tokenizer.SimpleAssignmentOperator) {
		p.consume(tokenizer.SimpleAssignmentOperator)

		value = p.expression()
	}

	end := p.consume(tokenizer.Semicolon).End
	return NewPropertyDefinition(key.Start(), end, key, value)
}

// FunctionExpression
// 	: '(' OptParameterList ')' BlockStatement
// 	;
func (p *parser) functionExpression() Node {
	start := p.consume(tokenizer.OpeningParenthesis).Start
	params := []Node{}
	if p.lookAhead.Not(tokenizer.ClosingParenthesis) {
		params = p.parameterList()
	}
	p.consume(tokenizer.ClosingParenthesis)

	body := p.blockStatement()

	return NewFunctionExpression(start, body.End(), params, body)
}

// ReturnStatement
// 	: 'return' OptExpression ';'
// 	;
func (p *parser) returnStatement() Node {
	start := p.consume(tokenizer.ReturnKeyword).Start

	var argument Node
	if p.lookAhead.Not(tokenizer.Semicolon) {
		argument = p.expression()
	}

	end := p.consume(tokenizer.Semicolon).End

	return NewReturnStatement(start, end, argument)
}

// FunctionDeclaration
// 	: 'function' Identifier '(' OptParameterList ')' BlockStatement
// 	;
func (p *parser) functionDeclaration() Node {
	start := p.consume(tokenizer.FunctionKeyword).Start

	id := p.identifier()

	p.consume(tokenizer.OpeningParenthesis)
	params := []Node{}
	if p.lookAhead.Not(tokenizer.ClosingParenthesis) {
		params = p.parameterList()
	}
	p.consume(tokenizer.ClosingParenthesis)

	body := p.blockStatement()

	return NewFunctionDeclaration(start, body.End(), id, params, body)
}

// ParameterList
// 	: Identifier
//	| ParameterList ',' Identifier
// 	;
func (p *parser) parameterList() []Node {
	ids := []Node{p.identifier()}

	for p.lookAhead.Is(tokenizer.Comma) {
		p.consume(tokenizer.Comma)
		ids = append(ids, p.identifier())
	}

	return ids
}

// IfStatement
// 	: 'if' ParenthesizedExpression Statement
// 	| 'if' ParenthesizedExpression Statement 'else' Statement
// 	;
func (p *parser) ifStatement() Node {
	start := p.consume(tokenizer.IfKeyword).Start

	test := p.parenthesizedExpression()

	consequent := p.statement()

	if p.lookAhead.Not(tokenizer.ElseKeyword) {
		return NewIfStatement(start, consequent.End(), test, consequent, nil)
	}

	p.consume(tokenizer.ElseKeyword)

	alternate := p.statement()

	return NewIfStatement(start, alternate.End(), test, consequent, alternate)
}

// IterationStatement
// 	: WhileStatment
// 	| DoWhileStatement
// 	| ForStatement
// 	;
func (p *parser) iterationStatement() Node {
	switch p.lookAhead.Type {
	case tokenizer.WhileKeyword:
		return p.whileStatement()
	case tokenizer.DoKeyword:
		return p.doWhileStatement()
	case tokenizer.ForKeyword:
		return p.forStatement()
	}

	panic("invalid look ahead for iteration statement")
}

// ForStatement
//	: 'for' '(' ForStatementInit ';' Expression ';' Expression ')' Statement
//	;
func (p *parser) forStatement() Node {
	start := p.consume(tokenizer.ForKeyword).Start
	p.consume(tokenizer.OpeningParenthesis)

	var init Node
	if !p.lookAhead.Is(tokenizer.Semicolon) {
		init = p.variableDeclarationInit()
	}
	p.consume(tokenizer.Semicolon)

	var test Node
	if !p.lookAhead.Is(tokenizer.Semicolon) {
		test = p.expression()
	}
	p.consume(tokenizer.Semicolon)

	var update Node
	if !p.lookAhead.Is(tokenizer.ClosingParenthesis) {
		update = p.expression()
	}
	p.consume(tokenizer.ClosingParenthesis)

	body := p.statement()

	return NewForStatement(start, body.End(), init, test, update, body)
}

// DoWhileStatement
//	: 'do' Statement 'while' ParenthesizedExpression ';'
// 	;
func (p *parser) doWhileStatement() Node {
	start := p.consume(tokenizer.DoKeyword).Start

	body := p.statement()

	p.consume(tokenizer.WhileKeyword)

	test := p.parenthesizedExpression()

	end := p.consume(tokenizer.Semicolon).End

	return NewDoWhileStatement(start, end, test, body)
}

// WhileStatement
//	: 'while' ParenthesizedExpression Statement
// 	;
func (p *parser) whileStatement() Node {
	start := p.consume(tokenizer.WhileKeyword).Start

	test := p.parenthesizedExpression()

	body := p.statement()

	return NewWhileStatement(start, body.End(), test, body)
}

// VariableDeclarationInit
// 	: VARIABLE_DECLARATION_KEYWORD VariableDeclaratorList
// 	;
func (p *parser) variableDeclarationInit() Node {
	kind := p.consume(tokenizer.VariableDeclarationKeyword)
	declarations := p.variableDeclaratorList()
	end := declarations[len(declarations)-1].End()

	return NewVariableDeclaration(kind.Start, end, kind.Value, declarations)
}

// VariableDeclaration
// 	: VariableDeclarationInit ';'
// 	;
func (p *parser) variableDeclaration() Node {
	init := p.variableDeclarationInit()
	end := p.consume(tokenizer.Semicolon).End
	init.SetEnd(end)

	return init
}

// VariableDeclaratorList
// 	: VariableDeclarator
//	| VariableDeclaratorList ',' VariableDeclarator
// 	;
func (p *parser) variableDeclaratorList() []Node {
	declarations := []Node{p.variableDeclarator()}

	for p.lookAhead.Is(tokenizer.Comma) {
		p.consume(tokenizer.Comma)
		declarations = append(declarations, p.variableDeclarator())
	}

	return declarations
}

// VariableDeclarator
// 	: Identifier OptVariableInitializer
// 	;
func (p *parser) variableDeclarator() Node {
	id := p.identifier()

	var init Node
	end := id.End()
	if p.lookAhead.Not(tokenizer.Semicolon, tokenizer.Comma) {
		p.consume(tokenizer.SimpleAssignmentOperator)
		init = p.assignmentExpression()
		end = init.End()
	}

	return NewVariableDeclarator(id.Start(), end, id, init)
}

// EmptyStatement
// 	: ';'
// 	;
func (p *parser) emptyStatement() Node {
	t := p.consume(tokenizer.Semicolon)

	return NewEmptyStatement(t.Start, t.End)
}

// BlockStatement
// 	: '{' StatementList '}'
// 	;
func (p *parser) blockStatement() Node {
	start := p.consume(tokenizer.OpeningCurlyBrace).Start

	sl := p.statementList(tokenizer.ClosingCurlyBrace)

	end := p.consume(tokenizer.ClosingCurlyBrace).End

	return NewBlockStatement(start, end, sl...)
}

// ExpressionStatment
// 	: Expression ';'
// 	;
func (p *parser) expressionStatment() Node {
	start := p.lookAhead.Start
	exp := p.expression()

	semi := p.consume(tokenizer.Semicolon)

	return NewExpressionStatement(start, semi.End, exp)
}

// Expression
// 	: AssignmentExpression
// 	;
func (p *parser) expression() Node {
	return p.assignmentExpression()
}

// AssignmentExpression
// 	: LogicalOrExpression
// 	| LeftHandSideExpression ASSIGNMENT_OPERATOR AssignmentExpression
// 	;
func (p *parser) assignmentExpression() Node {
	left := p.logicalOrExpression()

	if !p.isLookaheadAssignmentOperator() {
		return left
	}

	if left.Not(Identifier, MemberExpression) {
		panic(fmt.Errorf("invalid left-hand side expression: %v", left))
	}

	op := p.consumeAny().Value
	right := p.assignmentExpression()

	return NewAssignmentExpression(left.Start(), right.End(), op, left, right)
}

// LogicalOrExpression
// 	: LogicalAndExpression
// 	| LogicalOrExpression '||' LogicalAndExpression
// 	;
func (p *parser) logicalOrExpression() Node {
	return p.binaryExpression(
		p.logicalAndExpression,
		tokenizer.LogicalOrOperator,
		NewLogicalExpression,
	)
}

// LogicalAndExpression
// 	: EqualityExpression
// 	| LogicalAndExpression '&&' EqualityExpression
// 	;
func (p *parser) logicalAndExpression() Node {
	return p.binaryExpression(
		p.equalityExpression,
		tokenizer.LogicalAndOperator,
		NewLogicalExpression,
	)
}

// EqualityExpression
// 	: RelationalExpression
// 	| EqualityExpression EQUALITY_OPERATOR RelationalExpression
// 	;
func (p *parser) equalityExpression() Node {
	return p.binaryExpression(
		p.relationalExpression,
		tokenizer.EqualityOperator,
		NewBinaryExpression,
	)
}

// RelationalExpression
// 	: AdditiveExpression
// 	| RelationalExpression RELATIONAL_OPERATOR AdditiveExpression
// 	;
func (p *parser) relationalExpression() Node {
	return p.binaryExpression(
		p.additiveExpression,
		tokenizer.RelationalOperator,
		NewBinaryExpression,
	)
}

// AdditiveExpression
// 	: MultiplicativeExpression
// 	| AdditiveExpression ADDITIVE_OPERATOR MultiplicativeExpression
// 	;
func (p *parser) additiveExpression() Node {
	return p.binaryExpression(
		p.multiplicativeExpression,
		tokenizer.AdditiveOperator,
		NewBinaryExpression,
	)
}

// MultiplicativeExpression
// 	: UnaryExpression
// 	| MultiplicativeExpression MULTIPLICATIVE_OPERATOR UnaryExpression
// 	;
func (p *parser) multiplicativeExpression() Node {
	return p.binaryExpression(
		p.unaryExpression,
		tokenizer.MultiplicativeOperator,
		NewBinaryExpression,
	)
}

// UnaryExpression
// 	: LeftHandSideExpression
// 	| ADDITIVE_OPERATOR UnaryExpression
//	| LOGICAL_NOT UnaryExpression
// 	;
func (p *parser) unaryExpression() Node {
	if p.lookAhead.Not(tokenizer.LogicalNotOperator, tokenizer.AdditiveOperator) {
		return p.leftHandSideExpression()
	}

	operator := p.consumeAny()
	expr := p.unaryExpression()
	return NewUnaryExpression(
		operator.Start,
		expr.End(),
		operator.Value,
		expr,
	)
}

// LeftHandSideExpression
// 	: CallExpression
//	;
func (p *parser) leftHandSideExpression() Node {
	return p.callExpression()
}

// CallExpression
// 	: MemberExpression
//  | CallExpression '(' OptArgumentList ')'
//	;
func (p *parser) callExpression() Node {
	callee := p.memberExpression()

	for p.lookAhead.Is(tokenizer.OpeningParenthesis) {
		p.consume(tokenizer.OpeningParenthesis)
		arguments := []Node{}
		if p.lookAhead.Not(tokenizer.ClosingParenthesis) {
			arguments = p.argumentList()
		}
		end := p.consume(tokenizer.ClosingParenthesis).End

		callee = NewCallExpression(callee.Start(), end, callee, arguments)
	}

	return callee
}

// ArgumentList
//	: Expression
//	| ArgumentList ',' Expression
func (p *parser) argumentList() []Node {
	arguments := []Node{p.expression()}

	for p.lookAhead.Is(tokenizer.Comma) {
		p.consume(tokenizer.Comma)
		arguments = append(arguments, p.expression())
	}

	return arguments
}

// MemberExpression
// 	: PrimaryExpression
// 	| MemberExpression '.' Identifier
// 	| MemberExpression '[' Expression ']'
//	;
func (p *parser) memberExpression() Node {
	object := p.primaryExpression()

	for p.lookAhead.Is(tokenizer.Dot, tokenizer.OpeningBracket) {
		if p.lookAhead.Is(tokenizer.Dot) {
			p.consume(tokenizer.Dot)
			property := p.identifier()

			object = NewMemberExpression(object.Start(), property.End(), object, property, false)
		} else {
			p.consume(tokenizer.OpeningBracket)
			property := p.expression()
			end := p.consume(tokenizer.ClosingBracket).End

			object = NewMemberExpression(object.Start(), end, object, property, true)
		}

	}

	return object
}

// PrimaryExpression
// 	: Literal
//  | ParenthesizedExpression
//  | Identifier
// 	;
func (p *parser) primaryExpression() Node {
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
func (p *parser) identifier() Node {
	id := p.consume(tokenizer.Identifier)

	return NewIdentifier(id.Start, id.End, id.Value)
}

// ParenthesizedExpression
// 	: '(' Expression ')'
// 	;
func (p *parser) parenthesizedExpression() Node {
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
func (p *parser) literal() Node {
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
func (p *parser) booleanLiteral() Node {
	token := p.consume(tokenizer.BooleanLiteral)

	return NewLiteral(token.Start, token.End, token.Value == "true", token.Value)
}

// NullLiteral
// 	: 'null'
// 	;
func (p *parser) nullLiteral() Node {
	token := p.consume(tokenizer.NullLiteral)

	return NewLiteral(token.Start, token.End, nil, token.Value)
}

// NumericLiteral
// 	: NUMBER
// 	;
func (p *parser) numericLiteral() Node {
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
func (p *parser) stringLiteral() Node {
	token := p.consume(tokenizer.String)

	return NewLiteral(token.Start, token.End, token.Value[1:len(token.Value)-1], token.Value)
}
