package parser

type AstFactory interface {
	NumericLiteral(int) interface{}
	StringLiteral(string) interface{}
	ExpressionStatement(interface{}) interface{}
	BlockStatement(...interface{}) interface{}
	Program(...interface{}) interface{}
	EmptyStatement() interface{}
	BinaryExpression(operator string, left interface{}, right interface{}) interface{}
}
