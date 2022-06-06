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
	parserTest(t, `123;`, factory.Program(factory.ExpressionStatement(factory.Literal(123, 0, 3))))
}

func TestStrings(t *testing.T) {
	parserTest(t, `"Hello World!";`, factory.Program(factory.ExpressionStatement(factory.Literal("Hello World!", 0, 14))))
}

func TestStatements(t *testing.T) {
	parserTest(
		t,
		`1;2;3;`,
		factory.Program(
			factory.ExpressionStatement(factory.Literal(1, 0, 1)),
			factory.ExpressionStatement(factory.Literal(2, 2, 3)),
			factory.ExpressionStatement(factory.Literal(3, 4, 5)),
		))
}

func TestBlockStatement(t *testing.T) {
	parserTest(
		t,
		`{}`,
		factory.Program(factory.BlockStatement([]interface{}{}...)))

	parserTest(
		t,
		`{
			"Hello World!";
			{
				123;
			}
		}`,
		factory.Program(factory.BlockStatement(
			factory.ExpressionStatement(factory.Literal("Hello World!", 5, 19)),
			factory.BlockStatement(
				factory.ExpressionStatement(factory.Literal(123, 30, 33)),
			),
		)))

	parserTest(
		t,
		`{
			123;
			"Hello World!";
		}`,
		factory.Program(
			factory.BlockStatement(
				factory.ExpressionStatement(factory.Literal(123, 5, 8)),
				factory.ExpressionStatement(factory.Literal("Hello World!", 13, 27)),
			),
		),
	)
}

func TestEmptyStatement(t *testing.T) {
	parserTest(
		t,
		`;`,
		factory.Program(
			factory.EmptyStatement(),
		))
}

func TestAdditiveExpression(t *testing.T) {
	parserTest(
		t,
		`1+1;`,
		factory.Program(
			factory.ExpressionStatement(factory.BinaryExpression(
				"+",
				factory.Literal(1, 0, 1),
				factory.Literal(1, 2, 3),
			)),
		))

	parserTest(
		t,
		`1-1;`,
		factory.Program(
			factory.ExpressionStatement(factory.BinaryExpression(
				"-",
				factory.Literal(1, 0, 1),
				factory.Literal(1, 2, 3),
			)),
		))

	parserTest(
		t,
		`1+1-2;`,
		factory.Program(
			factory.ExpressionStatement(factory.BinaryExpression(
				"-",
				factory.BinaryExpression(
					"+",
					factory.Literal(1, 0, 1),
					factory.Literal(1, 2, 3),
				),
				factory.Literal(2, 4, 5),
			)),
		))
}

func TestMultiplicativeExpression(t *testing.T) {
	parserTest(
		t,
		`1*1;`,
		factory.Program(
			factory.ExpressionStatement(factory.BinaryExpression(
				"*",
				factory.Literal(1, 0, 1),
				factory.Literal(1, 2, 3),
			)),
		))

	parserTest(
		t,
		`1/1;`,
		factory.Program(
			factory.ExpressionStatement(factory.BinaryExpression(
				"/",
				factory.Literal(1, 0, 1),
				factory.Literal(1, 2, 3),
			)),
		))

	parserTest(
		t,
		`2+2*2;`,
		factory.Program(
			factory.ExpressionStatement(factory.BinaryExpression(
				"+",
				factory.Literal(2, 0, 1),
				factory.BinaryExpression(
					"*",
					factory.Literal(2, 2, 3),
					factory.Literal(2, 4, 5),
				),
			)),
		))

	parserTest(
		t,
		`2*2*2;`,
		factory.Program(
			factory.ExpressionStatement(factory.BinaryExpression(
				"*",
				factory.BinaryExpression(
					"*",
					factory.Literal(2, 0, 1),
					factory.Literal(2, 2, 3),
				),
				factory.Literal(2, 4, 5),
			)),
		))
}

func TestMultiplicativeExpressionPrecedence(t *testing.T) {
	parserTest(
		t,
		`(2+2)*2;`,
		factory.Program(
			factory.ExpressionStatement(factory.BinaryExpression(
				"*",
				factory.BinaryExpression(
					"+",
					factory.Literal(2, 1, 2),
					factory.Literal(2, 3, 4),
				),
				factory.Literal(2, 6, 7),
			)),
		))
}
