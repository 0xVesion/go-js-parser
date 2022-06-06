package parser

type AstFactory interface {
	Literal(interface{}, int, int) interface{}
	ExpressionStatement(interface{}) interface{}
	BlockStatement(...interface{}) interface{}
	Program(...interface{}) interface{}
	EmptyStatement() interface{}
	BinaryExpression(operator string, left interface{}, right interface{}) interface{}
	Identifier(name string) interface{}
	AssignmentExpression(operator string, left interface{}, right interface{}) interface{}
	IsIdentifier(val interface{}) bool
	VariableDeclaration(kind string, declarations []interface{}) interface{}
	VariableDeclarator(id interface{}, init interface{}) interface{}
	IfStatement(test interface{}, consequent interface{}, alternate interface{}) interface{}
	LogicalExpression(operator string, left interface{}, right interface{}) interface{}
	UnaryExpression(operator string, argument interface{}) interface{}
}
