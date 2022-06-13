package parser

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
	LogicalExpression         = "LogicalExpression"
	UnaryExpression           = "UnaryExpression"
)

type Node map[string]interface{}

func (n Node) Type() Type {
	return n["type"].(Type)
}

func (n Node) Start() int {
	return n["start"].(int)
}

func (n Node) End() int {
	return n["end"].(int)
}

func NewNode(t Type, start int, end int) Node {
	return Node{
		"type":  t,
		"start": start,
		"end":   end,
	}
}

type LiteralNode Node

func (n LiteralNode) Value() interface{} {
	return n["value"]
}

func NewLiteral(start int, end int, value interface{}, raw string) Node {
	n := NewNode(Literal, start, end)

	n["value"] = value
	n["raw"] = raw

	return n
}

func NewBlockStatement(start int, end int, body ...Node) Node {
	n := NewNode(BlockStatement, start, end)

	n["body"] = body

	return n
}

func NewProgram(start int, end int, body ...Node) Node {
	n := NewNode(Program, start, end)

	n["body"] = body
	n["sourceType"] = "script"

	return n
}
func NewEmptyStatement(start int, end int) Node {
	return NewNode(EmptyStatement, start, end)
}

type ExpressionStatementNode Node

func (n ExpressionStatementNode) Expression() Node {
	return n["expression"].(Node)
}

func NewExpressionStatement(start int, end int, expression Node) Node {
	n := NewNode(ExpressionStatement, start, end)

	n["expression"] = expression

	return n
}

func NewBinaryExpression(start int, end int, operator string, left Node, right Node) Node {
	n := NewNode(BinaryExpression, start, end)

	n["left"] = left
	n["operator"] = operator
	n["right"] = right

	return n
}

func NewIdentifier(start int, end int, name string) Node {
	n := NewNode(Identifier, start, end)

	n["name"] = name

	return n
}

func NewAssignmentExpression(start int, end int, operator string, left Node, right Node) Node {
	n := NewNode(AssignmentExpression, start, end)

	n["left"] = left
	n["operator"] = operator
	n["right"] = right

	return n
}

func NewVariableDeclaration(start int, end int, kind string, declarations []Node) Node {
	n := NewNode(VariableDeclaration, start, end)

	n["kind"] = kind
	n["declarations"] = declarations

	return n
}

func NewVariableDeclarator(start int, end int, id Node, init Node) Node {
	n := NewNode(VariableDeclarator, start, end)

	n["id"] = id
	n["init"] = init

	return n
}

func NewIfStatement(start int, end int, test Node, consequent Node, alternate Node) Node {
	n := NewNode(IfStatement, start, end)

	n["test"] = test
	n["consequent"] = consequent
	n["alternate"] = alternate

	return n
}

func NewLogicalExpression(operator string, left Node, right Node) Node {
	n := NewNode(LogicalExpression, 0, 0)

	n["left"] = left
	n["operator"] = operator
	n["right"] = right

	return n
}

func NewUnaryExpression(start int, end int, operator string, argument Node) Node {
	n := NewNode(UnaryExpression, start, end)

	n["operator"] = operator
	n["argument"] = argument
	n["prefix"] = true

	return n
}
