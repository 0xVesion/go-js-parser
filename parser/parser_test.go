package parser

import (
	"encoding/json"
	"testing"

	"github.com/0xvesion/go-parser/tokenizer"
)

func parserTest(t *testing.T, src string, expected program) {
	result, err := New(tokenizer.New(src)).Parse()
	if err != nil {
		t.Error(err)
	}

	toJson := func(x interface{}) string {
		j, _ := json.MarshalIndent(x, "", "  ")

		return string(j)
	}

	expectedJson := toJson(expected)
	resultJson := toJson(result)

	if expectedJson != resultJson {
		t.Errorf("Unexpected result. \nwant: %v \ngot: %v", expectedJson, resultJson)
	}
}

func TestRecognizesNumber(t *testing.T) {
	parserTest(t, `123;`, newProgram(newExpressionStatement(newNumericLiteral(123))))
}

func TestRecognizesStrings(t *testing.T) {
	parserTest(t, `"Hello World!";`, newProgram(newExpressionStatement(newStringLiteral("Hello World!"))))
}

func TestRecognizesStatements(t *testing.T) {
	parserTest(
		t,
		`1;2;3;`,
		newProgram(
			newExpressionStatement(newNumericLiteral(1)),
			newExpressionStatement(newNumericLiteral(2)),
			newExpressionStatement(newNumericLiteral(3)),
		))
}

func TestRecognizesBlockStatement(t *testing.T) {
	parserTest(
		t,
		`{}`,
		newProgram(newBlockStatement([]interface{}{}...)))

	parserTest(
		t,
		`{
			"Hello World!";
			{
				123;
			}
		}`,
		newProgram(newBlockStatement(
			newExpressionStatement(newStringLiteral("Hello World!")),
			newBlockStatement(
				newExpressionStatement(newNumericLiteral(123)),
			),
		)))

	parserTest(
		t,
		`{
			123;
			"Hello World!";
		}`,
		newProgram(
			newBlockStatement(
				newExpressionStatement(newNumericLiteral(123)),
				newExpressionStatement(newStringLiteral("Hello World!")),
			),
		),
	)
}
