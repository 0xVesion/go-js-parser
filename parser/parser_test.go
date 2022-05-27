package parser

import (
	"reflect"
	"testing"

	"github.com/0xvesion/go-parser/tokenizer"
)

func parserTest(t *testing.T, src string, expected Node) {
	result, err := New(tokenizer.New(src)).Parse()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Unexpected result. want: %v got: %v", expected, result)
	}
}

func TestRecognizesNumber(t *testing.T) {
	parserTest(t, `123;`, Node{Program, []Node{{ExpressionStatement, Node{NumericLiteral, 123}}}})
}

func TestRecognizesStrings(t *testing.T) {
	parserTest(t, `"Hello World!";`, Node{Program, []Node{{ExpressionStatement, Node{StringLiteral, "Hello World!"}}}})
}

func TestRecognizesStatements(t *testing.T) {
	parserTest(
		t,
		`1;2;3;`,
		Node{Program, []Node{
			{ExpressionStatement, Node{NumericLiteral, 1}},
			{ExpressionStatement, Node{NumericLiteral, 2}},
			{ExpressionStatement, Node{NumericLiteral, 3}}}})
}

func TestRecognizesBlockStatement(t *testing.T) {
	parserTest(
		t,
		`{}`,
		Node{Program, []Node{{BlockStatement, []Node{}}}})

	parserTest(
		t,
		`{
			"Hello World!";
			{
				123;
			}
		}`,
		Node{Program, []Node{{BlockStatement, []Node{
			{ExpressionStatement, Node{StringLiteral, "Hello World!"}},
			{BlockStatement, []Node{
				{ExpressionStatement, Node{NumericLiteral, 123}},
			}},
		}}}})

	parserTest(
		t,
		`{
			123;
			"Hello World!";
		}`,
		Node{Program, []Node{{BlockStatement, []Node{
			{ExpressionStatement, Node{NumericLiteral, 123}},
			{ExpressionStatement, Node{StringLiteral, "Hello World!"}},
		}}}})
}
