package tokenizer

type Type string

const (
	None                   Type = ""
	Number                      = "Number"
	String                      = "String"
	Semicolon                   = "Semicolon"
	OpeningCurlyBrace           = "OpeningCurlyBrace"
	ClosingCurlyBrace           = "ClosingCurlyBrace"
	AdditiveOperator            = "AdditiveOperator"
	MultiplicativeOperator      = "MultiplicativeOperator"
	OpeningParenthesis          = "OpeningParenthesis"
	ClosingParenthesis          = "ClosingParenthesis"
	Identifier                  = "Identifier"
	AssignmentOperator          = "AssignmentOperator"
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
	{OpeningCurlyBrace, []string{`{`}},
	{ClosingCurlyBrace, []string{`}`}},
	{AdditiveOperator, []string{`[\+-]`}},
	{MultiplicativeOperator, []string{`[\/*]`}},
	{OpeningParenthesis, []string{`\(`}},
	{ClosingParenthesis, []string{`\)`}},
	{Identifier, []string{`[a-zA-Z_$]\w*`}},
	{AssignmentOperator, []string{`[-+*/]?=`}},
}
