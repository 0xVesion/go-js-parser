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
	WhileStatement            = "WhileStatement"
	ForStatement              = "ForStatement"
	DoWhileStatement          = "DoWhileStatement"
	FunctionDeclaration       = "FunctionDeclaration"
	ReturnStatement           = "ReturnStatement"
	MemberExpression          = "MemberExpression"
	CallExpression            = "CallExpression"
	ClassDeclaration          = "ClassDeclaration"
	ClassBody                 = "ClassBody"
	PropertyDefinition        = "PropertyDefinition"
	MethodDefinition          = "MethodDefinition"
	FunctionExpression        = "FunctionExpression"
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

func (n Node) SetEnd(end int) {
	n["end"] = end
}

func (n Node) Is(types ...Type) bool {
	for _, t := range types {
		if n.Type() == t {
			return true
		}
	}

	return false
}

func (n Node) Not(types ...Type) bool {
	return !n.Is(types...)
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

func (n ExpressionStatementNode) SetDirective(directive string) {
	n["directive"] = directive
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

func NewLogicalExpression(start int, end int, operator string, left Node, right Node) Node {
	n := NewNode(LogicalExpression, start, end)

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

func NewWhileStatement(start int, end int, test Node, body Node) Node {
	n := NewNode(WhileStatement, start, end)

	n["test"] = test
	n["body"] = body

	return n
}

func NewForStatement(start int, end int, init Node, test Node, update Node, body Node) Node {
	n := NewNode(ForStatement, start, end)

	n["init"] = init
	n["test"] = test
	n["update"] = update
	n["body"] = body

	return n
}

func NewDoWhileStatement(start int, end int, test Node, body Node) Node {
	n := NewNode(DoWhileStatement, start, end)

	n["test"] = test
	n["body"] = body

	return n
}

func NewFunctionDeclaration(start int, end int, id Node, params []Node, body Node) Node {
	n := NewNode(FunctionDeclaration, start, end)

	n["id"] = id
	n["expression"] = false
	n["async"] = false
	n["generator"] = false
	n["params"] = params
	n["body"] = body

	return n
}

func NewReturnStatement(start int, end int, argument Node) Node {
	n := NewNode(ReturnStatement, start, end)

	n["argument"] = argument

	return n
}

func NewMemberExpression(start int, end int, object Node, property Node, computed bool) Node {
	n := NewNode(MemberExpression, start, end)

	n["object"] = object
	n["property"] = property
	n["computed"] = computed

	return n
}

func NewCallExpression(start int, end int, callee Node, arguments []Node) Node {
	n := NewNode(CallExpression, start, end)

	n["callee"] = callee
	n["arguments"] = arguments

	return n
}

func NewClassDeclaration(start int, end int, id Node, superClass Node, body Node) Node {
	n := NewNode(ClassDeclaration, start, end)

	n["id"] = id
	n["superClass"] = superClass
	n["body"] = body

	return n
}

func NewClassBody(start int, end int, body []Node) Node {
	n := NewNode(ClassBody, start, end)

	n["body"] = body

	return n
}

func NewPropertyDefinition(start int, end int, key Node, value Node) Node {
	n := NewNode(PropertyDefinition, start, end)

	n["static"] = false
	n["key"] = key
	n["value"] = value

	return n
}

func NewMethodDefinition(start int, end int, key Node, value Node) Node {
	n := NewNode(MethodDefinition, start, end)

	n["static"] = false
	n["kind"] = "constructor"
	n["key"] = key
	n["value"] = value

	return n
}

func NewFunctionExpression(start int, end int, body Node) Node {
	n := NewNode(FunctionExpression, start, end)

	n["expression"] = false
	n["generator"] = false
	n["async"] = false
	n["params"] = []Node{}
	n["body"] = body

	return n
}
