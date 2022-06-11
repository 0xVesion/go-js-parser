package parser

type Type string

const (
	ProgramType              Type = "Program"
	LiteralType                   = "Literal"
	ExpressionStatementType       = "ExpressionStatement"
	BlockStatementType            = "BlockStatement"
	EmptyStatementType            = "EmptyStatement"
	BinaryExpressionType          = "BinaryExpression"
	AssignmentExpressionType      = "AssignmentExpression"
	IdentifierType                = "Identifier"
	VariableDeclarationType       = "VariableDeclaration"
	VariableDeclaratorType        = "VariableDeclarator"
	IfStatementType               = "IfStatement"
	LogicalExpressionType         = "LogicalExpression"
	UnaryExpressionType           = "UnaryExpression"
)

type Literal struct {
	Type  `json:"type"`
	Value interface{} `json:"value"`
	Start int         `json:"start"`
	End   int         `json:"end"`
}

func NewLiteral(val interface{}, start int, end int) interface{} {
	return Literal{LiteralType, val, start, end}
}

type BlockStatement struct {
	Type  `json:"type"`
	Body  []interface{} `json:"body"`
	Start int           `json:"start"`
	End   int           `json:"end"`
}

func NewBlockStatement(start int, end int, sl ...interface{}) interface{} {
	return BlockStatement{BlockStatementType, sl, start, end}
}

type Program struct {
	Type  `json:"type"`
	Body  []interface{} `json:"body"`
	Start int           `json:"start"`
	End   int           `json:"end"`
}

func NewProgram(start int, end int, sl ...interface{}) interface{} {
	return Program{ProgramType, sl, start, end}
}

type EmptyStatement struct {
	Type  `json:"type"`
	Start int `json:"start"`
	End   int `json:"end"`
}

func NewEmptyStatement(start int, end int) interface{} {
	return EmptyStatement{EmptyStatementType, start, end}
}

type ExpressionStatement struct {
	Type       `json:"type"`
	Expression interface{} `json:"expression"`
	Start      int         `json:"start"`
	End        int         `json:"end"`
}

func NewExpressionStatement(exp interface{}, start int, end int) interface{} {
	return ExpressionStatement{ExpressionStatementType, exp, start, end}
}

type binaryExpression struct {
	Type     `json:"type"`
	Operator string      `json:"operator"`
	Left     interface{} `json:"left"`
	Right    interface{} `json:"right"`
}

func NewBinaryExpression(operator string, left interface{}, right interface{}) interface{} {
	return binaryExpression{BinaryExpressionType, operator, left, right}
}

type assignmentExpression struct {
	Type     `json:"type"`
	Operator string      `json:"operator"`
	Left     interface{} `json:"left"`
	Right    interface{} `json:"right"`
}

func NewAssignmentExpression(operator string, left interface{}, right interface{}) interface{} {
	return assignmentExpression{AssignmentExpressionType, operator, left, right}
}

type identifier struct {
	Type `json:"type"`
	Name string `json:"name"`
}

func NewIdentifier(name string) interface{} {
	return identifier{IdentifierType, name}
}

func NewIsIdentifier(val interface{}) bool {
	_, ok := val.(identifier)

	return ok
}

type variableDeclaration struct {
	Type         `json:"type"`
	Kind         string        `json:"kind"`
	Declarations []interface{} `json:"declarations"`
}

func NewVariableDeclaration(kind string, declarations []interface{}) interface{} {
	return variableDeclaration{VariableDeclarationType, kind, declarations}
}

type variableDeclarator struct {
	Type `json:"type"`
	Id   interface{} `json:"id"`
	Init interface{} `json:"init"`
}

func NewVariableDeclarator(id interface{}, init interface{}) interface{} {
	return variableDeclarator{VariableDeclaratorType, id, init}
}

type ifStatement struct {
	Type       `json:"type"`
	Test       interface{} `json:"test"`
	Consequent interface{} `json:"consequent"`
	Alternate  interface{} `json:"alternate"`
}

func NewIfStatement(test interface{}, consequent interface{}, alternate interface{}) interface{} {
	return ifStatement{IfStatementType, test, consequent, alternate}
}

type logicalExpression struct {
	Type     `json:"type"`
	Operator string      `json:"operator"`
	Left     interface{} `json:"left"`
	Right    interface{} `json:"right"`
}

func NewLogicalExpression(operator string, left interface{}, right interface{}) interface{} {
	return logicalExpression{LogicalExpressionType, operator, left, right}
}

type unaryExpression struct {
	Type     `json:"type"`
	Operator string      `json:"operator"`
	Prefix   bool        `json:"prefix"`
	Argument interface{} `json:"argument"`
}

func NewUnaryExpression(operator string, argument interface{}) interface{} {
	return unaryExpression{UnaryExpressionType, operator, true, argument}
}
