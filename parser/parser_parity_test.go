package parser_test

import (
	"encoding/json"
	"io"
	"log"
	"os/exec"
	"reflect"
	"testing"

	"github.com/0xvesion/go-js-parser/parser"
	jsonastfactory "github.com/0xvesion/go-js-parser/parser/json_ast_factory"
	"github.com/0xvesion/go-js-parser/tokenizer"
)

func sanatize(node map[string]interface{}, keys []string) interface{} {
	for key, value := range node {
		if subNode, is := value.(map[string]interface{}); is {
			node[key] = sanatize(subNode, keys)
			continue
		}

		if list, is := value.([]interface{}); is {
			for i, subNode := range list {
				list[i] = sanatize(subNode.(map[string]interface{}), keys)
			}
			continue
		}

		for _, remove := range keys {
			if key == remove {
				delete(node, key)
			}
		}
	}

	return node
}

func acorn(exp string) interface{} {
	cmd := exec.Command("npx", "acorn", "--ecma9")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, exp)
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	x := &map[string]interface{}{}

	json.Unmarshal(out, x)

	return sanatize(*x, []string{"start", "end", "sourceType", "raw", "directive"})
}

func goParser(src string) interface{} {
	actualAst, err := parser.New(tokenizer.New(src), jsonastfactory.New()).Parse()
	if err != nil {
		log.Fatal(err)
	}
	actualAstJson, _ := json.MarshalIndent(actualAst, "", "  ")
	x := &map[string]interface{}{}
	json.Unmarshal(actualAstJson, x)

	return *x
}

func test(t *testing.T, src string) {
	referenceAst := acorn(src)
	actualAst := goParser(src)

	if !reflect.DeepEqual(referenceAst, actualAst) {
		referenceJson, _ := json.MarshalIndent(referenceAst, "", "  ")
		actualJson, _ := json.MarshalIndent(actualAst, "", "  ")

		t.Errorf("Invalid ast.\nwant: %s\ngot: %s\n", referenceJson, actualJson)
	}
}

func TestNumberParity(t *testing.T) {
	test(t, `123;`)
}

func TestStringsParity(t *testing.T) {
	test(t, `"Hello World!";`)
}

func TestStatementsParity(t *testing.T) {
	test(t, `1;2;3;`)
}

func TestBlockStatementParity(t *testing.T) {
	test(t, `{}`)

	test(t, `{
			"Hello World!";
			{
				123;
			}
		}`)

	test(t,
		`{
			123;
			"Hello World!";
		}`)
}

func TestEmptyStatementParity(t *testing.T) {
	test(t, `;`)
}

func TestAdditiveExpressionParity(t *testing.T) {
	test(t, `1+1;`)

	test(t, `1-1;`)

	test(t, `1+1-2;`)
}

func TestMultiplicativeExpressionParity(t *testing.T) {
	test(t, `1*1;`)

	test(t, `1/1;`)

	test(t, `2+2*2;`)

	test(t, `2*2*2;`)
}

func TestMultiplicativeExpressionPrecedenceParity(t *testing.T) {
	test(t, `(2+2)*2;`)
}

func TestAssignments(t *testing.T) {
	test(t, `a = 1;`)
	test(t, `a = y = 1;`)
	test(t, `a = 1 + 2;`)
}

func TestVariableDeclaration(t *testing.T) {
	test(t, `let a;`)
	test(t, `let a = 1;`)
	test(t, `let a, b;`)
	test(t, `let a, b = 1;`)
}

func TestIfStatement(t *testing.T) {
	test(t, `if (a > b) result = 100; else result = 200;`)
	test(t, `if (a > b) result = 100;`)
	test(t, `if (a > b) if (c > d) result = 123; else result = 321; else result = 111;`)
}

func TestRelationalExpression(t *testing.T) {
	test(t, `1>2;`)
	test(t, `1+1<=2;`)
	test(t, `a>a>a;`)
}

func TestBooleanParity(t *testing.T) {
	test(t, `true;`)
	test(t, `false;`)
}

func TestNullParity(t *testing.T) {
	test(t, `null;`)
}
