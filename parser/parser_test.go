package parser_test

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func cache(src string, producer func(string) ([]byte, error)) ([]byte, error) {
	path := fmt.Sprintf("/tmp/%s", hash(src))

	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		bytes, err := os.ReadFile(path)
		if err != nil {
			return []byte{}, fmt.Errorf("cannot read file: %w", err)
		}

		return bytes, nil
	}

	bytes, err := producer(src)
	if err != nil {
		return []byte{}, err
	}

	err = os.WriteFile(path, bytes, os.ModePerm)
	if err != nil {
		return []byte{}, fmt.Errorf("cannot write file: %w", err)
	}

	return bytes, nil
}

func acornRaw(exp string) ([]byte, error) {
	cmd := exec.Command("npx", "acorn", "--ecma13")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return []byte{}, fmt.Errorf("cannot spawn acorn: %w", err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, exp)
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		return []byte{}, fmt.Errorf("acorn error: %w", err)
	}

	return out, nil
}

func acorn(exp string) (interface{}, error) {
	out, err := cache(exp, acornRaw)
	if err != nil {
		return nil, err
	}

	x := &map[string]interface{}{}

	err = json.Unmarshal(out, x)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal acorn result: %w", err)
	}

	return x, nil
}

func goParser(src string) (interface{}, error) {
	actualAst, err := parser.New(tokenizer.New(src)).Parse()
	if err != nil {
		return nil, err
	}
	actualAstJson, _ := json.MarshalIndent(actualAst, "", "  ")
	x := &map[string]interface{}{}
	err = json.Unmarshal(actualAstJson, x)
	if err != nil {
		return nil, err
	}

	return x, nil
}

func test(t *testing.T, src string) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatal(r)
		}
	}()

	referenceAst, err := acorn(src)
	if err != nil {
		t.Error(err)
		return
	}
	actualAst, err := goParser(src)
	if err != nil {
		t.Error(err)
		return
	}

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

func TestClassDeclaration(t *testing.T) {
	test(t, `class Test {}`)
	test(t, `class Rectangle extends Drawable {}`)
	test(t, `class Test {
		test = 123;
		test2;
	}`)
	test(t, `class Test {
		constructor(foo, bar) {}
		test(foo, bar) {}
	}`)
	test(t, `class Point extends Vector2D {
		constructor(x, y, color) {
			super(x, y);

			this.color = color;
		}
	}`)
	test(t, `class Test {
		set name(str) {
			this.myName = name;
		}
	  
		get name() {
			return this.myName;
		}
	}`)

	// TOOD: Add support for getters/setters
	// TODO: Add support for static/async modifiers
}
