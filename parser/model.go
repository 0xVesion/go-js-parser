package parser

type literal[T any] struct {
	Type
	Value T
}

func newNumericLiteral(val int) literal[int] {
	return literal[int]{NumericLiteral, val}
}

func newStringLiteral(val string) literal[string] {
	return literal[string]{StringLiteral, val}
}

type expressionStatement struct {
	Type
	Expression interface{}
}

func newExpressionStatement(exp interface{}) expressionStatement {
	return expressionStatement{ExpressionStatement, exp}
}

type blockStatement struct {
	Type
	Body []interface{}
}

func newBlockStatement(sl ...interface{}) blockStatement {
	return blockStatement{BlockStatement, sl}
}

type program struct {
	Type
	Body []interface{}
}

func newProgram(sl ...interface{}) program {
	return program{Program, sl}
}