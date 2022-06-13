package parser

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/0xvesion/go-js-parser/tokenizer"
)

type Parser interface {
	Parse() (interface{}, error)
}

type parser struct {
	t         tokenizer.Tokenizer
	lookAhead tokenizer.Token
}

func New(t tokenizer.Tokenizer) Parser {
	lookAhead, err := t.Next()
	if err != nil {
		panic(err)
	}

	return &parser{
		t:         t,
		lookAhead: lookAhead,
	}
}

func (p *parser) Parse() (n interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = p.formatError(r)
		}
	}()

	n = p.program()

	return
}

func (p *parser) formatError(err any) error {
	res := ""
	cursor := p.t.Cursor()
	resLine := 0
	resCol := 0

	for i, line := range strings.Split(p.t.Src(), "\n") {
		res += fmt.Sprintf("%02d  %s\n", i, line)
		if cursor-len(line) <= 0 {
			res += "  "
			for ii := 0; ii < cursor; ii++ {
				res += " "
			}
			res += "  ^"

			resLine = i
			resCol = cursor
			break
		}
		cursor -= len(line)
	}

	res = fmt.Sprintf("Ln %02d, Col %02d\n%s", resLine, resCol, res)

	return fmt.Errorf("%v\n%s\n%s", err, res, debug.Stack())
}

func (p *parser) consume(t tokenizer.Type) tokenizer.Token {
	token := p.lookAhead

	if token.Type != t {
		panic(fmt.Errorf("unexpected token type. want: %s got: %s", t, token.Type))
	}

	lookAhead, err := p.t.Next()
	if err != nil {
		panic(err)
	}
	p.lookAhead = lookAhead

	return token
}

func (p *parser) binaryExpression(builder func() interface{}, operator tokenizer.Type) interface{} {
	left := builder()

	for p.lookAhead.Type == operator {
		operator := p.consume(operator)
		right := builder()

		left = NewBinaryExpression(operator.Value, left, right)
	}

	return left
}

func (p *parser) isLookaheadLiteral() bool {
	return p.lookAhead.Type == tokenizer.Number ||
		p.lookAhead.Type == tokenizer.String ||
		p.lookAhead.Type == tokenizer.NullLiteral ||
		p.lookAhead.Type == tokenizer.BooleanLiteral
}

func (p *parser) isLookaheadAssignmentOperator() bool {
	return p.lookAhead.Type == tokenizer.SimpleAssignmentOperator ||
		p.lookAhead.Type == tokenizer.ComplexAssignmentOperator
}

func (p *parser) logicalExpression(builder func() interface{}, operator tokenizer.Type) interface{} {
	left := builder()

	for p.lookAhead.Type == operator {
		operator := p.consume(operator)
		right := builder()

		left = NewLogicalExpression(operator.Value, left, right)
	}

	return left
}

func (p *parser) consumeAny() tokenizer.Token {
	return p.consume(p.lookAhead.Type)
}

func (p *parser) addDirectives(sl []interface{}) []interface{} {
	for k, v := range sl {
		if exp, ok := v.(ExpressionStatementNode); ok {
			if literal, ok := exp.Expression.(LiteralNode); ok {
				if str, ok := literal.Value.(string); ok {
					exp.Directive = str
					sl[k] = exp
				}
			}
		}
	}

	return sl
}
