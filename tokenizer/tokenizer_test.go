package tokenizer

import (
	"reflect"
	"testing"
)

func all(t *tokenizer) ([]Token, error) {
	tokens := []Token{}
	for t.HasNext() {
		token, err := t.Next()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)

	}

	return tokens, nil
}

func tokenizerTest(t *testing.T, src string, expected []Token) {
	result, err := all(New(src).(*tokenizer))
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Unexpected result. want: %v got: %v", expected, result)
	}
}

func TestRecognizesNumber(t *testing.T) {
	tokenizerTest(t, "123", []Token{{Number, "123"}})
}

func TestSkipWhitespace(t *testing.T) {
	tokenizerTest(
		t,
		"   1 2 3   ",
		[]Token{{Number, "1"}, {Number, "2"}, {Number, "3"}, {}},
	)
}

func TestRecognizesStrings(t *testing.T) {
	tokenizerTest(t, `'Hello World!'`, []Token{{String, `'Hello World!'`}})
	tokenizerTest(t, `"Hello World!"`, []Token{{String, `"Hello World!"`}})
}
