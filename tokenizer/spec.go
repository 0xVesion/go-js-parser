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
	{RelationalOperator, []string{`[><]=?`}},
	{SimpleAssignmentOperator, []string{`=`}},
	{ComplexAssignmentOperator, []string{`[-+*/]=`}},
	{AdditiveOperator, []string{`[\+-]`}},
	{MultiplicativeOperator, []string{`[\/*]`}},
	{OpeningParenthesis, []string{`\(`}},
	{ClosingParenthesis, []string{`\)`}},
	{VariableDeclarationKeyword, []string{`\b(let)|(const)\b`}},
	{IfKeyword, []string{`\bif\b`}},
	{ElseKeyword, []string{`\belse\b`}},
	{BooleanLiteral, []string{`\b(true)|(false)\b`}},
	{NullLiteral, []string{`\bnull\b`}},
	{Identifier, []string{`[a-zA-Z_$]\w*`}},
}
