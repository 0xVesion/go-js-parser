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
	parserTest(t, `123`, Node{Program, Node{NumericLiteral, 123}})
}

func TestRecognizesStrings(t *testing.T) {
	parserTest(t, `"Hello World!"`, Node{Program, Node{StringLiteral, "Hello World!"}})
}
