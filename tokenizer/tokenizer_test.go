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
	tokenizerTest(t, "123", []Token{{Number, "123", 0, 3}})
}

func TestSkipWhitespace(t *testing.T) {
	tokenizerTest(
		t,
		"   1 2 3   ",
		[]Token{{Number, "1", 3, 4}, {Number, "2", 5, 6}, {Number, "3", 7, 8}, {}},
	)
}

func TestSkipComments(t *testing.T) {
	tokenizerTest(
		t,
		"123/* foobar */321",
		[]Token{{Number, "123", 0, 3}, {Number, "321", 15, 18}},
	)

	tokenizerTest(
		t,
		`
		123 // foobar
		// foobar
		321`,
		[]Token{{Number, "123", 3, 6}, {Number, "321", 31, 34}},
	)
}

func TestRecognizesStrings(t *testing.T) {
	tokenizerTest(t, `'Hello World!'`, []Token{{String, `'Hello World!'`, 0, 14}})
	tokenizerTest(t, `"Hello World!"`, []Token{{String, `"Hello World!"`, 0, 14}})
}

func TestRecognizesSemicolon(t *testing.T) {
	tokenizerTest(t, `;`, []Token{{Semicolon, `;`, 0, 1}})
}

func TestRecognizesCurlyBraces(t *testing.T) {
	tokenizerTest(t, `{`, []Token{{OpeningCurlyBrace, `{`, 0, 1}})
	tokenizerTest(t, `}`, []Token{{ClosingCurlyBrace, `}`, 0, 1}})
}

func TestRecognizesAdditiveOperators(t *testing.T) {
	tokenizerTest(t, `+`, []Token{{AdditiveOperator, `+`, 0, 1}})
	tokenizerTest(t, `-`, []Token{{AdditiveOperator, `-`, 0, 1}})
}

func TestRecognizesMultiplicativeOperators(t *testing.T) {
	tokenizerTest(t, `*`, []Token{{MultiplicativeOperator, `*`, 0, 1}})
	tokenizerTest(t, `/`, []Token{{MultiplicativeOperator, `/`, 0, 1}})
}

func TestRecognizesBraces(t *testing.T) {
	tokenizerTest(t, `(`, []Token{{OpeningParenthesis, `(`, 0, 1}})
	tokenizerTest(t, `)`, []Token{{ClosingParenthesis, `)`, 0, 1}})
}

func TestRecognizesIdentifiers(t *testing.T) {
	tokenizerTest(t, `test`, []Token{{Identifier, `test`, 0, 4}})
	tokenizerTest(t, `_`, []Token{{Identifier, `_`, 0, 1}})
	tokenizerTest(t, `$`, []Token{{Identifier, `$`, 0, 1}})
	tokenizerTest(t, `TEST`, []Token{{Identifier, `TEST`, 0, 4}})
	tokenizerTest(t, `$test123`, []Token{{Identifier, `$test123`, 0, 8}})
	tokenizerTest(t, `_test123`, []Token{{Identifier, `_test123`, 0, 8}})
}

func TestRecognizesAssignmentOperator(t *testing.T) {
	tokenizerTest(t, `=`, []Token{{SimpleAssignmentOperator, `=`, 0, 1}})
	tokenizerTest(t, `+=`, []Token{{ComplexAssignmentOperator, `+=`, 0, 2}})
	tokenizerTest(t, `*=`, []Token{{ComplexAssignmentOperator, `*=`, 0, 2}})
	tokenizerTest(t, `-=`, []Token{{ComplexAssignmentOperator, `-=`, 0, 2}})
	tokenizerTest(t, `/=`, []Token{{ComplexAssignmentOperator, `/=`, 0, 2}})
}

func TestRecognizesKeywords(t *testing.T) {
	tokenizerTest(t, `const`, []Token{{VariableDeclarationKeyword, `const`, 0, 5}})
	tokenizerTest(t, `let`, []Token{{VariableDeclarationKeyword, `let`, 0, 3}})
	tokenizerTest(t, `letter`, []Token{{Identifier, `letter`, 0, 6}})
	tokenizerTest(t, `aconst`, []Token{{Identifier, `aconst`, 0, 6}})
	tokenizerTest(t, `if`, []Token{{IfKeyword, `if`, 0, 2}})
	tokenizerTest(t, `else`, []Token{{ElseKeyword, `else`, 0, 4}})
	tokenizerTest(t, `true`, []Token{{BooleanLiteral, `true`, 0, 4}})
	tokenizerTest(t, `trueism`, []Token{{Identifier, `trueism`, 0, 7}})
	tokenizerTest(t, `sofalse`, []Token{{Identifier, `sofalse`, 0, 7}})
	tokenizerTest(t, `false`, []Token{{BooleanLiteral, `false`, 0, 5}})
	tokenizerTest(t, `null`, []Token{{NullLiteral, `null`, 0, 4}})
	tokenizerTest(t, `while`, []Token{{WhileKeyword, `while`, 0, 5}})
	tokenizerTest(t, `do`, []Token{{DoKeyword, `do`, 0, 2}})
	tokenizerTest(t, `for`, []Token{{ForKeyword, `for`, 0, 3}})
	tokenizerTest(t, `function`, []Token{{FunctionKeyword, `function`, 0, 8}})
	tokenizerTest(t, `return`, []Token{{ReturnKeyword, `return`, 0, 6}})
	tokenizerTest(t, `class`, []Token{{ClassKeyword, `class`, 0, 5}})
	tokenizerTest(t, `new`, []Token{{NewKeyword, `new`, 0, 3}})
	tokenizerTest(t, `this`, []Token{{ThisKeyword, `this`, 0, 4}})
	tokenizerTest(t, `extends`, []Token{{ExtendsKeyword, `extends`, 0, 7}})
	tokenizerTest(t, `super`, []Token{{SuperKeyword, `super`, 0, 5}})
}

func TestRelationalOperators(t *testing.T) {
	tokenizerTest(t, `<`, []Token{{RelationalOperator, `<`, 0, 1}})
	tokenizerTest(t, `<=`, []Token{{RelationalOperator, `<=`, 0, 2}})
	tokenizerTest(t, `>`, []Token{{RelationalOperator, `>`, 0, 1}})
	tokenizerTest(t, `>=`, []Token{{RelationalOperator, `>=`, 0, 2}})
}

func TestEqualityOperators(t *testing.T) {
	tokenizerTest(t, `==`, []Token{{EqualityOperator, `==`, 0, 2}})
	tokenizerTest(t, `===`, []Token{{EqualityOperator, `===`, 0, 3}})
	tokenizerTest(t, `!=`, []Token{{EqualityOperator, `!=`, 0, 2}})
	tokenizerTest(t, `!==`, []Token{{EqualityOperator, `!==`, 0, 3}})
}

func TestLogicalOperators(t *testing.T) {
	tokenizerTest(t, `||`, []Token{{LogicalOrOperator, `||`, 0, 2}})
	tokenizerTest(t, `&&`, []Token{{LogicalAndOperator, `&&`, 0, 2}})
	tokenizerTest(t, `!`, []Token{{LogicalNotOperator, `!`, 0, 1}})
}

func TestRecognizesDot(t *testing.T) {
	tokenizerTest(t, `.`, []Token{{Dot, `.`, 0, 1}})
}

func TestRecognizesBrackets(t *testing.T) {
	tokenizerTest(t, `[`, []Token{{OpeningBracket, `[`, 0, 1}})
	tokenizerTest(t, `]`, []Token{{ClosingBracket, `]`, 0, 1}})
}
