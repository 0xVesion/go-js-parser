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
}
