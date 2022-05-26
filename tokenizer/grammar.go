package tokenizer

type Type int

const (
	Number Type = iota
	Identifier
	Symbol
)

var grammar = map[Type]string{
	Number:     `\d+`,
	Identifier: `\w+`,
	Symbol:     `[\+-/*]`,
}
