package parser_test

import (
	"encoding/json"
	"testing"

	"github.com/0xvesion/go-js-parser/parser"
	"github.com/0xvesion/go-js-parser/tokenizer"
)

func parserTest(t *testing.T, src string, expected interface{}) {
	result, err := parser.New(tokenizer.New(src)).Parse()
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

func TestNumber(t *testing.T) {
	parserTest(t, `123;`, parser.NewProgram(0, 4, parser.NewExpressionStatement(parser.NewLiteral(123, 0, 3), 0, 4)))
}

func TestStrings(t *testing.T) {
	parserTest(t, `"Hello World!";`, parser.NewProgram(0, 15, parser.NewExpressionStatement(parser.NewLiteral("Hello World!", 0, 14), 0, 15)))
}

func TestStatements(t *testing.T) {
	parserTest(
		t,
		`1;2;3;`,
		parser.NewProgram(
			0, 6,
			parser.NewExpressionStatement(parser.NewLiteral(1, 0, 1), 0, 2),
			parser.NewExpressionStatement(parser.NewLiteral(2, 2, 3), 2, 4),
			parser.NewExpressionStatement(parser.NewLiteral(3, 4, 5), 4, 6),
		))
}

func TestBlockStatement(t *testing.T) {
	parserTest(
		t,
		`{}`,
		parser.NewProgram(0, 2, parser.NewBlockStatement(0, 2, []interface{}{}...)))

	parserTest(
		t,
		`{
			"Hello World!";
			{
				123;
			}
		}`,
		parser.NewProgram(0, 43, parser.NewBlockStatement(0, 43,
			parser.NewExpressionStatement(parser.NewLiteral("Hello World!", 5, 19), 5, 20),
			parser.NewBlockStatement(24, 39,
				parser.NewExpressionStatement(parser.NewLiteral(123, 30, 33), 30, 34),
			),
		)))

	parserTest(
		t,
		`{
			123;
			"Hello World!";
		}`,
		parser.NewProgram(
			0, 32,
			parser.NewBlockStatement(0, 32,
				parser.NewExpressionStatement(parser.NewLiteral(123, 5, 8), 5, 9),
				parser.NewExpressionStatement(parser.NewLiteral("Hello World!", 13, 27), 13, 28),
			),
		),
	)
}

func TestEmptyStatement(t *testing.T) {
	parserTest(
		t,
		`;`,
		parser.NewProgram(
			0, 1,
			parser.NewEmptyStatement(0, 1),
		))
}

func TestAdditiveExpression(t *testing.T) {
	parserTest(
		t,
		`1+1;`,
		parser.NewProgram(
			0, 4,
			parser.NewExpressionStatement(parser.NewBinaryExpression(
				"+",
				parser.NewLiteral(1, 0, 1),
				parser.NewLiteral(1, 2, 3),
			), 0, 4),
		))

	parserTest(
		t,
		`1-1;`,
		parser.NewProgram(
			0, 4,
			parser.NewExpressionStatement(parser.NewBinaryExpression(
				"-",
				parser.NewLiteral(1, 0, 1),
				parser.NewLiteral(1, 2, 3),
			), 0, 4),
		))

	parserTest(
		t,
		`1+1-2;`,
		parser.NewProgram(
			0, 6,
			parser.NewExpressionStatement(parser.NewBinaryExpression(
				"-",
				parser.NewBinaryExpression(
					"+",
					parser.NewLiteral(1, 0, 1),
					parser.NewLiteral(1, 2, 3),
				),
				parser.NewLiteral(2, 4, 5),
			), 0, 6),
		))
}

func TestMultiplicativeExpression(t *testing.T) {
	parserTest(
		t,
		`1*1;`,
		parser.NewProgram(
			0, 4,
			parser.NewExpressionStatement(parser.NewBinaryExpression(
				"*",
				parser.NewLiteral(1, 0, 1),
				parser.NewLiteral(1, 2, 3),
			), 0, 4),
		))

	parserTest(
		t,
		`1/1;`,
		parser.NewProgram(
			0, 4,
			parser.NewExpressionStatement(parser.NewBinaryExpression(
				"/",
				parser.NewLiteral(1, 0, 1),
				parser.NewLiteral(1, 2, 3),
			), 0, 4),
		))

	parserTest(
		t,
		`2+2*2;`,
		parser.NewProgram(
			0, 6,
			parser.NewExpressionStatement(parser.NewBinaryExpression(
				"+",
				parser.NewLiteral(2, 0, 1),
				parser.NewBinaryExpression(
					"*",
					parser.NewLiteral(2, 2, 3),
					parser.NewLiteral(2, 4, 5),
				),
			), 0, 6),
		))

	parserTest(
		t,
		`2*2*2;`,
		parser.NewProgram(
			0, 6,
			parser.NewExpressionStatement(parser.NewBinaryExpression(
				"*",
				parser.NewBinaryExpression(
					"*",
					parser.NewLiteral(2, 0, 1),
					parser.NewLiteral(2, 2, 3),
				),
				parser.NewLiteral(2, 4, 5),
			), 0, 6),
		))
}

func TestMultiplicativeExpressionPrecedence(t *testing.T) {
	parserTest(
		t,
		`(2+2)*2;`,
		parser.NewProgram(
			0, 8,
			parser.NewExpressionStatement(parser.NewBinaryExpression(
				"*",
				parser.NewBinaryExpression(
					"+",
					parser.NewLiteral(2, 1, 2),
					parser.NewLiteral(2, 3, 4),
				),
				parser.NewLiteral(2, 6, 7),
			), 0, 8),
		))
}
