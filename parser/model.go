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

type Node struct {
	Type  `json:"type"`
	Start int `json:"start"`
	End   int `json:"end"`
}

type LiteralNode struct {
	Node
	Value interface{} `json:"value"`
	Raw   interface{} `json:"raw"`
}

func NewLiteral(start int, end int, val interface{}, raw string) LiteralNode {
	return LiteralNode{Node{Literal, start, end}, val, raw}
}

type BlockStatementNode struct {
	Node
	Body []interface{} `json:"body"`
}

func NewBlockStatement(start int, end int, sl ...interface{}) BlockStatementNode {
	return BlockStatementNode{Node{BlockStatement, start, end}, sl}
}

type ProgramNode struct {
	Node
	SourceType string        `json:"sourceType"`
	Body       []interface{} `json:"body"`
}

func NewProgram(start int, end int, sl ...interface{}) ProgramNode {
	return ProgramNode{Node{Program, start, end}, "script", sl}
}

type EmptyStatementNode struct {
	Node
}

func NewEmptyStatement(start int, end int) EmptyStatementNode {
	return EmptyStatementNode{Node{EmptyStatement, start, end}}
}

type ExpressionStatementNode struct {
	Node
	Expression interface{} `json:"expression"`
	Directive  string      `json:"directive,omitempty"`
}

func NewExpressionStatement(start int, end int, exp interface{}) ExpressionStatementNode {
	return ExpressionStatementNode{Node{ExpressionStatement, start, end}, exp, ""}
}

type BinaryExpressionNode struct {
	Node
	Operator string      `json:"operator"`
	Left     interface{} `json:"left"`
	Right    interface{} `json:"right"`
}

func NewBinaryExpression(start int, end int, operator string, left interface{}, right interface{}) BinaryExpressionNode {
	return BinaryExpressionNode{Node{BinaryExpression, start, end}, operator, left, right}
}

type IdentifierNode struct {
	Node
	Name string `json:"name"`
}

func NewIdentifier(start int, end int, name string) IdentifierNode {
	return IdentifierNode{Node{Identifier, start, end}, name}
}

type AssignmentExpressionNode struct {
	Node
	Operator string      `json:"operator"`
	Left     interface{} `json:"left"`
	Right    interface{} `json:"right"`
}

func NewAssignmentExpression(start int, end int, operator string, left interface{}, right interface{}) AssignmentExpressionNode {
	return AssignmentExpressionNode{Node{AssignmentExpression, start, end}, operator, left, right}
}

type VariableDeclarationNode struct {
	Node
	Kind         string        `json:"kind"`
	Declarations []interface{} `json:"declarations"`
}

func NewVariableDeclaration(kind string, declarations []interface{}) VariableDeclarationNode {
	return VariableDeclarationNode{Node{VariableDeclaration, 0, 0}, kind, declarations}
}

type VariableDeclaratorNode struct {
	Node
	Id   interface{} `json:"id"`
	Init interface{} `json:"init"`
}

func NewVariableDeclarator(id interface{}, init interface{}) VariableDeclaratorNode {
	return VariableDeclaratorNode{Node{VariableDeclarator, 0, 0}, id, init}
}

type IfStatementNode struct {
	Node
	Test       interface{} `json:"test"`
	Consequent interface{} `json:"consequent"`
	Alternate  interface{} `json:"alternate"`
}

func NewIfStatement(test interface{}, consequent interface{}, alternate interface{}) IfStatementNode {
	return IfStatementNode{Node{IfStatement, 0, 0}, test, consequent, alternate}
}

type LogicalExpressionNode struct {
	Node
	Operator string      `json:"operator"`
	Left     interface{} `json:"left"`
	Right    interface{} `json:"right"`
}

func NewLogicalExpression(operator string, left interface{}, right interface{}) LogicalExpressionNode {
	return LogicalExpressionNode{Node{LogicalExpression, 0, 0}, operator, left, right}
}

type UnaryExpressionNode struct {
	Node
	Operator string      `json:"operator"`
	Prefix   bool        `json:"prefix"`
	Argument interface{} `json:"argument"`
}

func NewUnaryExpression(operator string, argument interface{}) UnaryExpressionNode {
	return UnaryExpressionNode{Node{UnaryExpression, 0, 0}, operator, true, argument}
}
