package tokenizer

type Type string

const (
	None                       Type = ""
	Number                          = "Number"
	String                          = "String"
	Semicolon                       = "Semicolon"
	OpeningCurlyBrace               = "OpeningCurlyBrace"
	ClosingCurlyBrace               = "ClosingCurlyBrace"
	AdditiveOperator                = "AdditiveOperator"
	MultiplicativeOperator          = "MultiplicativeOperator"
	OpeningParenthesis              = "OpeningParenthesis"
	ClosingParenthesis              = "ClosingParenthesis"
	Identifier                      = "Identifier"
	SimpleAssignmentOperator        = "SimpleAssignmentOperator"
	ComplexAssignmentOperator       = "ComplexAssignmentOperator"
	VariableDeclarationKeyword      = "VariableDeclarationKeyword"
	Comma                           = "Comma"
	IfKeyword                       = "IfKeyword"
	ElseKeyword                     = "ElseKeyword"
	RelationalOperator              = "RelationalOperator"
	BooleanLiteral                  = "BooleanLiteral"
	NullLiteral                     = "NullLiteral"
	EqualityOperator                = "EqualityOperator"
	LogicalOrOperator               = "LogicalOrOperator"
	LogicalAndOperator              = "LogicalAndOperator"
	LogicalNotOperator              = "LogicalNotOperator"
	WhileKeyword                    = "WhileKeyword"
	DoKeyword                       = "DoKeyword"
	ForKeyword                      = "ForKeyword"
	FunctionKeyword                 = "FunctionKeyword"
)

type specEntry struct {
	Type
	Regexp []string
}

var spec = []specEntry{
	{None, []string{`\s+`, `\/\*[\s\S]*\*\/`, `\/\/.*`}},
	{Number, []string{`\d+`}},
	{String, []string{`\".*\"`, `\'.*\'`}},
	{Semicolon, []string{`;`}},
	{Comma, []string{`,`}},
	{OpeningCurlyBrace, []string{`{`}},
	{ClosingCurlyBrace, []string{`}`}},
	{LogicalOrOperator, []string{`\|\|`}},
	{LogicalAndOperator, []string{`&&`}},
	{EqualityOperator, []string{`[!=]==?`}},
	{LogicalNotOperator, []string{`!`}},
	{RelationalOperator, []string{`[><]=?`}},
	{SimpleAssignmentOperator, []string{`=`}},
	{ComplexAssignmentOperator, []string{`[-+*/]=`}},
	{AdditiveOperator, []string{`[\+-]`}},
	{MultiplicativeOperator, []string{`[\/*]`}},
	{OpeningParenthesis, []string{`\(`}},
	{ClosingParenthesis, []string{`\)`}},
	{VariableDeclarationKeyword, []string{`\b((let)|(const))\b`}},
	{IfKeyword, []string{`\bif\b`}},
	{ElseKeyword, []string{`\belse\b`}},
	{WhileKeyword, []string{`\bwhile\b`}},
	{DoKeyword, []string{`\bdo\b`}},
	{ForKeyword, []string{`\bfor\b`}},
	{FunctionKeyword, []string{`\bfunction\b`}},
	{BooleanLiteral, []string{`\b((true)|(false))\b`}},
	{NullLiteral, []string{`\bnull\b`}},
	{Identifier, []string{`[a-zA-Z_$]\w*`}},
}
