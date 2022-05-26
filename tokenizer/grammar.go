package tokenizer

type Type string

const (
	None   Type = ""
	Number      = "Number"
	String      = "String"
)

type Grammar struct {
	Type
	Regexp []string
}

var grammar = []Grammar{
	{Number, []string{`\d+`}},
	{String, []string{`\".*\"`, `\'.*\'`}},
}
