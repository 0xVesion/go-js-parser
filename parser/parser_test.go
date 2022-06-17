package parser_test

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/0xvesion/go-js-parser/parser"
	"github.com/0xvesion/go-js-parser/tokenizer"
)

func hash(s string) string {
	hasher := sha1.New()
	hasher.Write([]byte(s))
	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash)
}

func cache(src string, producer func(string) []byte) []byte {
	path := fmt.Sprintf("/tmp/%s", hash(src))

	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		bytes, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}

		return bytes
	}

	bytes := producer(src)
	err := os.WriteFile(path, bytes, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return bytes
}

func acornRaw(exp string) []byte {
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

	return out
}

func acorn(exp string) interface{} {
	out := cache(exp, acornRaw)

	x := &map[string]interface{}{}

	json.Unmarshal(out, x)

	return x
}

func goParser(src string) interface{} {
	actualAst, err := parser.New(tokenizer.New(src)).Parse()
	if err != nil {
		log.Fatal(err)
	}
	actualAstJson, _ := json.MarshalIndent(actualAst, "", "  ")
	x := &map[string]interface{}{}
	json.Unmarshal(actualAstJson, x)

	return x
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
	test(t, `(1+2)*3;`)
	test(t, `3*(1+2);`)
	test(t, `((5+5) * (6+4)) / 5;`)
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
	test(t, `a>b>c;`)
}

func TestBooleanParity(t *testing.T) {
	test(t, `true;`)
	test(t, `false;`)
}

func TestNullParity(t *testing.T) {
	test(t, `null;`)
}

func TestEqualityParity(t *testing.T) {
	test(t, `1==2;`)
	test(t, `a+1!=2+2*3;`)
	test(t, `if(1!=2) res = 200;`)
}

func TestLogicalExpressionParity(t *testing.T) {
	test(t, `a||b;`)
	test(t, `a&&b;`)
	test(t, `(a&&b)||c;`)
	test(t, `a||(b&&c);`)
	test(t, `((a&&b)||(c&&d))||e;`)
	test(t, `a = 1||2==2&&3;`)
	test(t, `a||1!=2+2||3;`)
	test(t, `a||1&&2+2||3;`)
}

func TestUnarityExpressionParity(t *testing.T) {
	test(t, `-1;`)
	test(t, `+1;`)
	test(t, `!true;`)
	test(t, `!!true;`)
}

func TestLoops(t *testing.T) {
	test(t, `while (i > 0) {
		i-=1;
	}`)

	test(t, `for (let i = 0; i < 10; i+=1) {
		result = 10;
	}`)

	test(t, `for (;;) {
		result = 10;
	}`)

	test(t, `do {
		i-=1;
	} while (i > 0);`)
}

func TestFunctions(t *testing.T) {
	test(t, `function test(a, b, c) {
		result = 123;
	}`)

	test(t, `function noArgs() {
		return "i don't have any args";
	}`)

	test(t, `function square(x) {
		return x * x;
	}`)

	test(t, `function optExpr(x) {
		return;
	}`)
}

func TestMemberExpressions(t *testing.T) {
	test(t, `x.y;`)
	test(t, `x.y.z;`)
	test(t, `x.y['test'];`)
	test(t, `x.y = 2;`)
	test(t, `x[0] = 2;`)
	test(t, `x.z.y['test'];`)
}

func TestCallExpression(t *testing.T) {
	test(t, `test();`)
	test(t, `test('12343');`)
	test(t, `test('12343', 321);`)
	test(t, `console.log('1235');`)
	test(t, `log()();`)
}
