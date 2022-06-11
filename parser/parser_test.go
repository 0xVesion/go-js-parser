package parser_test

import (
	"encoding/json"
	"testing"

	"github.com/0xvesion/go-js-parser/parser"
	jsonastfactory "github.com/0xvesion/go-js-parser/parser/json_ast_factory"
	"github.com/0xvesion/go-js-parser/tokenizer"
)

var factory = jsonastfactory.New()

func parserTest(t *testing.T, src string, expected interface{}) {
	result, err := parser.New(tokenizer.New(src), factory).Parse()
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
	parserTest(t, `123;`, factory.Program(0, 4, factory.ExpressionStatement(factory.Literal(123, 0, 3), 0, 4)))
}

func TestStrings(t *testing.T) {
	parserTest(t, `"Hello World!";`, factory.Program(0, 15, factory.ExpressionStatement(factory.Literal("Hello World!", 0, 14), 0, 15)))
}

func TestStatements(t *testing.T) {
	parserTest(
		t,
		`1;2;3;`,
		factory.Program(
			0, 6,
			factory.ExpressionStatement(factory.Literal(1, 0, 1), 0, 2),
			factory.ExpressionStatement(factory.Literal(2, 2, 3), 2, 4),
			factory.ExpressionStatement(factory.Literal(3, 4, 5), 4, 6),
		))
}

func TestBlockStatement(t *testing.T) {
	parserTest(
		t,
		`{}`,
		factory.Program(0, 2, factory.BlockStatement(0, 2, []interface{}{}...)))

	parserTest(
		t,
		`{
			"Hello World!";
			{
				123;
			}
		}`,
		factory.Program(0, 43, factory.BlockStatement(0, 43,
			factory.ExpressionStatement(factory.Literal("Hello World!", 5, 19), 5, 20),
			factory.BlockStatement(24, 39,
				factory.ExpressionStatement(factory.Literal(123, 30, 33), 30, 34),
			),
		)))

	parserTest(
		t,
		`{
			123;
			"Hello World!";
		}`,
		factory.Program(
			0, 32,
			factory.BlockStatement(0, 32,
				factory.ExpressionStatement(factory.Literal(123, 5, 8), 5, 9),
				factory.ExpressionStatement(factory.Literal("Hello World!", 13, 27), 13, 28),
			),
		),
	)
}

func TestEmptyStatement(t *testing.T) {
	parserTest(
		t,
		`;`,
		factory.Program(
			0, 1,
			factory.EmptyStatement(0, 1),
		))
}

func TestAdditiveExpression(t *testing.T) {
	parserTest(
		t,
		`1+1;`,
		factory.Program(
			0, 4,
			factory.ExpressionStatement(factory.BinaryExpression(
				"+",
				factory.Literal(1, 0, 1),
				factory.Literal(1, 2, 3),
			), 0, 4),
		))

	parserTest(
		t,
		`1-1;`,
		factory.Program(
			0, 4,
			factory.ExpressionStatement(factory.BinaryExpression(
				"-",
				factory.Literal(1, 0, 1),
				factory.Literal(1, 2, 3),
			), 0, 4),
		))

	parserTest(
		t,
		`1+1-2;`,
		factory.Program(
			0, 6,
			factory.ExpressionStatement(factory.BinaryExpression(
				"-",
				factory.BinaryExpression(
					"+",
					factory.Literal(1, 0, 1),
					factory.Literal(1, 2, 3),
				),
				factory.Literal(2, 4, 5),
			), 0, 6),
		))
}

func TestMultiplicativeExpression(t *testing.T) {
	parserTest(
		t,
		`1*1;`,
		factory.Program(
			0, 4,
			factory.ExpressionStatement(factory.BinaryExpression(
				"*",
				factory.Literal(1, 0, 1),
				factory.Literal(1, 2, 3),
			), 0, 4),
		))

	parserTest(
		t,
		`1/1;`,
		factory.Program(
			0, 4,
			factory.ExpressionStatement(factory.BinaryExpression(
				"/",
				factory.Literal(1, 0, 1),
				factory.Literal(1, 2, 3),
			), 0, 4),
		))

	parserTest(
		t,
		`2+2*2;`,
		factory.Program(
			0, 6,
			factory.ExpressionStatement(factory.BinaryExpression(
				"+",
				factory.Literal(2, 0, 1),
				factory.BinaryExpression(
					"*",
					factory.Literal(2, 2, 3),
					factory.Literal(2, 4, 5),
				),
			), 0, 6),
		))

	parserTest(
		t,
		`2*2*2;`,
		factory.Program(
			0, 6,
			factory.ExpressionStatement(factory.BinaryExpression(
				"*",
				factory.BinaryExpression(
					"*",
					factory.Literal(2, 0, 1),
					factory.Literal(2, 2, 3),
				),
				factory.Literal(2, 4, 5),
			), 0, 6),
		))
}

func TestMultiplicativeExpressionPrecedence(t *testing.T) {
	parserTest(
		t,
		`(2+2)*2;`,
		factory.Program(
			0, 8,
			factory.ExpressionStatement(factory.BinaryExpression(
				"*",
				factory.BinaryExpression(
					"+",
					factory.Literal(2, 1, 2),
					factory.Literal(2, 3, 4),
				),
				factory.Literal(2, 6, 7),
			), 0, 8),
		))
}
