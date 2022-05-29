package jsonastfactory

import "github.com/0xvesion/go-parser/parser"

type Type string

const (
	Program             Type = "Program"
	Literal                  = "Literal"
	ExpressionStatement      = "ExpressionStatement"
	BlockStatement           = "BlockStatement"
	EmptyStatement           = "EmptyStatement"
	BinaryExpression         = "BinaryExpression"
)

type factory struct{}

func New() parser.AstFactory {
	return factory{}
}

type literal struct {
	Type  `json:"type"`
	Value interface{} `json:"value"`
}

func (factory) Literal(val interface{}) interface{} {
	return literal{Literal, val}
}

type blockStatement struct {
	Type `json:"type"`
	Body []interface{} `json:"body"`
}

func (factory) BlockStatement(sl ...interface{}) interface{} {
	return blockStatement{BlockStatement, sl}
}

type program struct {
	Type `json:"type"`
	Body []interface{} `json:"body"`
}

func (factory) Program(sl ...interface{}) interface{} {
	return program{Program, sl}
}

type emptyStatement struct {
	Type `json:"type"`
}

func (factory) EmptyStatement() interface{} {
	return emptyStatement{EmptyStatement}
}

type expressionStatement struct {
	Type       `json:"type"`
	Expression interface{} `json:"expression"`
}

func (factory) ExpressionStatement(exp interface{}) interface{} {
	return expressionStatement{ExpressionStatement, exp}
}

func (factory) BinaryExpression(operator string, left interface{}, right interface{}) interface{} {
	return binaryExpression{BinaryExpression, operator, left, right}
}

type binaryExpression struct {
	Type     `json:"type"`
	Operator string      `json:"operator"`
	Left     interface{} `json:"left"`
	Right    interface{} `json:"right"`
}
