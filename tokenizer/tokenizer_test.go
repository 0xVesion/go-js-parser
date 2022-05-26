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

func TestRecognizesIdentifiers(t *testing.T) {
	tokenizerTest(t, "test", []Token{{Identifier, "test"}})
}

func TestRecognizesIdentifiersWithNumbers(t *testing.T) {
	tokenizerTest(t, "test123", []Token{{Identifier, "test123"}})
}

func TestRecognizesNumbersAndIdentifiers(t *testing.T) {
	tokenizerTest(t, "123test123", []Token{{Number, "123"}, {Identifier, "test123"}})
}

func TestRecognizesSymbols(t *testing.T) {
	tokenizerTest(t, "+", []Token{{Symbol, "+"}})
	tokenizerTest(t, "-", []Token{{Symbol, "-"}})
	tokenizerTest(t, "*", []Token{{Symbol, "*"}})
	tokenizerTest(t, "/", []Token{{Symbol, "/"}})
}

func TestRecognizesNumbersAndSymbols(t *testing.T) {
	tokenizerTest(t, "1+1", []Token{{Number, "1"}, {Symbol, "+"}, {Number, "1"}})
}
