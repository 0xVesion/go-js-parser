package jsonastfactory

import "github.com/0xvesion/go-parser/parser"

type Type string

const (
	Program             Type = "Program"
	NumericLiteral           = "NumericLiteral"
	StringLiteral            = "StringLiteral"
	ExpressionStatement      = "ExpressionStatement"
	BlockStatement           = "BlockStatement"
	EmptyStatement           = "EmptyStatement"
	AdditiveExpression       = "AdditiveExpression"
)

type factory struct{}

func New() parser.AstFactory {
	return factory{}
}

type literal[T any] struct {
	Type
	Value T
}

func (factory) NumericLiteral(val int) interface{} {
	return literal[int]{NumericLiteral, val}
}

func (factory) StringLiteral(val string) interface{} {
	return literal[string]{StringLiteral, val}
}

type blockStatement struct {
	Type
	Body []interface{}
}

func (factory) BlockStatement(sl ...interface{}) interface{} {
	return blockStatement{BlockStatement, sl}
}

type program struct {
	Type
	Body []interface{}
}

func (factory) Program(sl ...interface{}) interface{} {
	return program{Program, sl}
}

type emptyStatement struct {
	Type
}

func (factory) EmptyStatement() interface{} {
	return emptyStatement{EmptyStatement}
}

type expressionStatement struct {
	Type
	Expression interface{}
}

func (factory) ExpressionStatement(exp interface{}) interface{} {
	return expressionStatement{ExpressionStatement, exp}
}

func (factory) AdditiveExpression(operator string, left interface{}, right interface{}) interface{} {
	return additiveExpression{AdditiveExpression, operator, left, right}
}

type additiveExpression struct {
	Type
	Operator string
	Left     interface{}
	Right    interface{}
}
