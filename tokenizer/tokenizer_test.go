package tokenizer

import (
	"reflect"
	"testing"
)

func tokenizerTest(t *testing.T, src string, expected []Token) {
	result, err := New(src).All()
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

func TestRecognizesStrings(t *testing.T) {
	tokenizerTest(t, `'Hello World!'`, []Token{{String, `'Hello World!'`}})
	tokenizerTest(t, `"Hello World!"`, []Token{{String, `"Hello World!"`}})
}
