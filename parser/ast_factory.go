package parser

type AstFactory interface {
	Literal(interface{}) interface{}
	ExpressionStatement(interface{}) interface{}
	BlockStatement(...interface{}) interface{}
	Program(...interface{}) interface{}
	EmptyStatement() interface{}
	BinaryExpression(operator string, left interface{}, right interface{}) interface{}
	Identifier(name string) interface{}
	AssignmentExpression(operator string, left interface{}, right interface{}) interface{}
	IsIdentifier(val interface{}) bool
}
