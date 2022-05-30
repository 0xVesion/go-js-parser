package jsonastfactory

import "github.com/0xvesion/go-js-parser/parser"

type Type string

const (
	Program              Type = "Program"
	Literal                   = "Literal"
	ExpressionStatement       = "ExpressionStatement"
	BlockStatement            = "BlockStatement"
	EmptyStatement            = "EmptyStatement"
	BinaryExpression          = "BinaryExpression"
	AssignmentExpression      = "AssignmentExpression"
	Identifier                = "Identifier"
	VariableDeclaration       = "VariableDeclaration"
	VariableDeclarator        = "VariableDeclarator"
	IfStatement               = "IfStatement"
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

func (factory) AssignmentExpression(operator string, left interface{}, right interface{}) interface{} {
	return assignmentExpression{AssignmentExpression, operator, left, right}
}

type assignmentExpression struct {
	Type     `json:"type"`
	Operator string      `json:"operator"`
	Left     interface{} `json:"left"`
	Right    interface{} `json:"right"`
}

func (factory) Identifier(name string) interface{} {
	return identifier{Identifier, name}
}

type identifier struct {
	Type `json:"type"`
	Name string `json:"name"`
}

func (factory) IsIdentifier(val interface{}) bool {
	_, ok := val.(identifier)

	return ok
}

func (factory) VariableDeclaration(kind string, declarations []interface{}) interface{} {
	return variableDeclaration{VariableDeclaration, kind, declarations}
}

type variableDeclaration struct {
	Type         `json:"type"`
	Kind         string        `json:"kind"`
	Declarations []interface{} `json:"declarations"`
}

func (factory) VariableDeclarator(id interface{}, init interface{}) interface{} {
	return variableDeclarator{VariableDeclarator, id, init}
}

type variableDeclarator struct {
	Type `json:"type"`
	Id   interface{} `json:"id"`
	Init interface{} `json:"init"`
}

func (factory) IfStatement(test interface{}, consequent interface{}, alternate interface{}) interface{} {
	return ifStatement{IfStatement, test, consequent, alternate}
}

type ifStatement struct {
	Type       `json:"type"`
	Test       interface{} `json:"test"`
	Consequent interface{} `json:"consequent"`
	Alternate  interface{} `json:"alternate"`
}
