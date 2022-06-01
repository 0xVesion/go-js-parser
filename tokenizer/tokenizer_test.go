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

func TestRecognizesMultiplicativeOperators(t *testing.T) {
	tokenizerTest(t, `*`, []Token{{MultiplicativeOperator, `*`}})
	tokenizerTest(t, `/`, []Token{{MultiplicativeOperator, `/`}})
}

func TestRecognizesBraces(t *testing.T) {
	tokenizerTest(t, `(`, []Token{{OpeningParenthesis, `(`}})
	tokenizerTest(t, `)`, []Token{{ClosingParenthesis, `)`}})
}

func TestRecognizesIdentifiers(t *testing.T) {
	tokenizerTest(t, `test`, []Token{{Identifier, `test`}})
	tokenizerTest(t, `_`, []Token{{Identifier, `_`}})
	tokenizerTest(t, `$`, []Token{{Identifier, `$`}})
	tokenizerTest(t, `TEST`, []Token{{Identifier, `TEST`}})
	tokenizerTest(t, `$test123`, []Token{{Identifier, `$test123`}})
	tokenizerTest(t, `_test123`, []Token{{Identifier, `_test123`}})
}

func TestRecognizesAssignmentOperator(t *testing.T) {
	tokenizerTest(t, `=`, []Token{{SimpleAssignmentOperator, `=`}})
	tokenizerTest(t, `+=`, []Token{{ComplexAssignmentOperator, `+=`}})
	tokenizerTest(t, `*=`, []Token{{ComplexAssignmentOperator, `*=`}})
	tokenizerTest(t, `-=`, []Token{{ComplexAssignmentOperator, `-=`}})
	tokenizerTest(t, `/=`, []Token{{ComplexAssignmentOperator, `/=`}})
}

func TestRecognizesKeywords(t *testing.T) {
	tokenizerTest(t, `const`, []Token{{VariableDeclarationKeyword, `const`}})
	tokenizerTest(t, `let`, []Token{{VariableDeclarationKeyword, `let`}})
	tokenizerTest(t, `if`, []Token{{IfKeyword, `if`}})
	tokenizerTest(t, `else`, []Token{{ElseKeyword, `else`}})
	tokenizerTest(t, `true`, []Token{{BooleanLiteral, `true`}})
	tokenizerTest(t, `false`, []Token{{BooleanLiteral, `false`}})
	tokenizerTest(t, `null`, []Token{{NullLiteral, `null`}})
}

func TestRelationalOperators(t *testing.T) {
	tokenizerTest(t, `<`, []Token{{RelationalOperator, `<`}})
	tokenizerTest(t, `<=`, []Token{{RelationalOperator, `<=`}})
	tokenizerTest(t, `>`, []Token{{RelationalOperator, `>`}})
	tokenizerTest(t, `>=`, []Token{{RelationalOperator, `>=`}})
}

func TestEqualityOperators(t *testing.T) {
	tokenizerTest(t, `==`, []Token{{EqualityOperator, `==`}})
	tokenizerTest(t, `===`, []Token{{EqualityOperator, `===`}})
	tokenizerTest(t, `!=`, []Token{{EqualityOperator, `!=`}})
	tokenizerTest(t, `!==`, []Token{{EqualityOperator, `!==`}})
}

func TestLogicalOperators(t *testing.T) {
	tokenizerTest(t, `||`, []Token{{LogicalOrOperator, `||`}})
	tokenizerTest(t, `&&`, []Token{{LogicalAndOperator, `&&`}})
}
