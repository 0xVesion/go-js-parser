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

func TestSkipComments(t *testing.T) {
	tokenizerTest(
		t,
		"123/* foobar */321",
		[]Token{{Number, "123"}, {Number, "321"}},
	)

	tokenizerTest(
		t,
		`
		123 // foobar
		// foobar
		321`,
		[]Token{{Number, "123"}, {Number, "321"}},
	)
}

func TestRecognizesStrings(t *testing.T) {
	tokenizerTest(t, `'Hello World!'`, []Token{{String, `'Hello World!'`}})
	tokenizerTest(t, `"Hello World!"`, []Token{{String, `"Hello World!"`}})
}

func TestRecognizesSemicolon(t *testing.T) {
	tokenizerTest(t, `;`, []Token{{Semicolon, `;`}})
}

func TestRecognizesCurlyBraces(t *testing.T) {
	tokenizerTest(t, `{`, []Token{{OpeningCurlyBrace, `{`}})
	tokenizerTest(t, `}`, []Token{{ClosingCurlyBrace, `}`}})
}

func TestRecognizesAdditiveOperators(t *testing.T) {
	tokenizerTest(t, `+`, []Token{{AdditiveOperator, `+`}})
	tokenizerTest(t, `-`, []Token{{AdditiveOperator, `-`}})
}
