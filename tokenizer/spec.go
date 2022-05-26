package tokenizer

type Type string

const (
	None   Type = ""
	Number      = "Number"
	String      = "String"
)

type specEntry struct {
	Type
	Regexp []string
}

var spec = []specEntry{
	{None, []string{`\s+`}},
	{Number, []string{`\d+`}},
	{String, []string{`\".*\"`, `\'.*\'`}},
}
