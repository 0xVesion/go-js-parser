package parser_test

import (
	"encoding/json"
	"testing"

	"github.com/0xvesion/go-parser/parser"
	jsonastfactory "github.com/0xvesion/go-parser/parser/json_ast_factory"
	"github.com/0xvesion/go-parser/tokenizer"
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
	parserTest(t, `123;`, factory.Program(factory.ExpressionStatement(factory.Literal(123))))
}

func TestStrings(t *testing.T) {
	parserTest(t, `"Hello World!";`, factory.Program(factory.ExpressionStatement(factory.Literal("Hello World!"))))
}

func TestStatements(t *testing.T) {
	parserTest(
		t,
		`1;2;3;`,
		factory.Program(
			factory.ExpressionStatement(factory.Literal(1)),
			factory.ExpressionStatement(factory.Literal(2)),
			factory.ExpressionStatement(factory.Literal(3)),
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
			factory.ExpressionStatement(factory.Literal("Hello World!")),
			factory.BlockStatement(
				factory.ExpressionStatement(factory.Literal(123)),
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
				factory.ExpressionStatement(factory.Literal(123)),
				factory.ExpressionStatement(factory.Literal("Hello World!")),
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
				factory.Literal(1),
				factory.Literal(1),
			)),
		))

	parserTest(
		t,
		`1-1;`,
		factory.Program(
			factory.ExpressionStatement(factory.BinaryExpression(
				"-",
				factory.Literal(1),
				factory.Literal(1),
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
					factory.Literal(1),
					factory.Literal(1),
				),
				factory.Literal(2),
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
				factory.Literal(1),
				factory.Literal(1),
			)),
		))

	parserTest(
		t,
		`1/1;`,
		factory.Program(
			factory.ExpressionStatement(factory.BinaryExpression(
				"/",
				factory.Literal(1),
				factory.Literal(1),
			)),
		))

	parserTest(
		t,
		`2+2*2;`,
		factory.Program(
			factory.ExpressionStatement(factory.BinaryExpression(
				"+",
				factory.Literal(2),
				factory.BinaryExpression(
					"*",
					factory.Literal(2),
					factory.Literal(2),
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
					factory.Literal(2),
					factory.Literal(2),
				),
				factory.Literal(2),
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
					factory.Literal(2),
					factory.Literal(2),
				),
				factory.Literal(2),
			)),
		))
}
