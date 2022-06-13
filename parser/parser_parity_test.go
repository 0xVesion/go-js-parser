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

func acorn(exp string, sanatizeKeys []string) interface{} {
	out := cache(exp, acornRaw)

	x := &map[string]interface{}{}

	json.Unmarshal(out, x)

	return sanatize(*x, sanatizeKeys)
}

func goParser(src string, sanatizeKeys []string) interface{} {
	actualAst, err := parser.New(tokenizer.New(src)).Parse()
	if err != nil {
		log.Fatal(err)
	}
	actualAstJson, _ := json.MarshalIndent(actualAst, "", "  ")
	x := &map[string]interface{}{}
	json.Unmarshal(actualAstJson, x)

	return sanatize(*x, sanatizeKeys)
}

func test(t *testing.T, src string, sanatizeKeys ...string) {
	referenceAst := acorn(src, sanatizeKeys)
	actualAst := goParser(src, sanatizeKeys)

	if !reflect.DeepEqual(referenceAst, actualAst) {
		referenceJson, _ := json.MarshalIndent(referenceAst, "", "  ")
		actualJson, _ := json.MarshalIndent(actualAst, "", "  ")

		t.Errorf("Invalid ast.\nwant: %s\ngot: %s\n", referenceJson, actualJson)
	}
}

func sanatizeTest(t *testing.T, src string) {
	test(t, src, "start", "end")
}

func TestNumberParity(t *testing.T) {
	sanatizeTest(t, `123;`)
}

func TestStringsParity(t *testing.T) {
	sanatizeTest(t, `"Hello World!";`)
}

func TestStatementsParity(t *testing.T) {
	sanatizeTest(t, `1;2;3;`)
}

func TestBlockStatementParity(t *testing.T) {
	sanatizeTest(t, `{}`)

	sanatizeTest(t, `{
			"Hello World!";
			{
				123;
			}
		}`)

	sanatizeTest(t,
		`{
			123;
			"Hello World!";
		}`)
}

func TestEmptyStatementParity(t *testing.T) {
	sanatizeTest(t, `;`)
}

func TestAdditiveExpressionParity(t *testing.T) {
	sanatizeTest(t, `1+1;`)

	sanatizeTest(t, `1-1;`)

	sanatizeTest(t, `1+1-2;`)
}

func TestMultiplicativeExpressionParity(t *testing.T) {
	sanatizeTest(t, `1*1;`)

	sanatizeTest(t, `1/1;`)

	sanatizeTest(t, `2+2*2;`)

	sanatizeTest(t, `2*2*2;`)
}

func TestMultiplicativeExpressionPrecedenceParity(t *testing.T) {
	sanatizeTest(t, `(2+2)*2;`)
}

func TestAssignments(t *testing.T) {
	sanatizeTest(t, `a = 1;`)
	sanatizeTest(t, `a = y = 1;`)
	sanatizeTest(t, `a = 1 + 2;`)
}

func TestVariableDeclaration(t *testing.T) {
	sanatizeTest(t, `let a;`)
	sanatizeTest(t, `let a = 1;`)
	sanatizeTest(t, `let a, b;`)
	sanatizeTest(t, `let a, b = 1;`)
}

func TestIfStatement(t *testing.T) {
	sanatizeTest(t, `if (a > b) result = 100; else result = 200;`)
	sanatizeTest(t, `if (a > b) result = 100;`)
	sanatizeTest(t, `if (a > b) if (c > d) result = 123; else result = 321; else result = 111;`)
}

func TestRelationalExpression(t *testing.T) {
	sanatizeTest(t, `1>2;`)
	sanatizeTest(t, `1+1<=2;`)
	sanatizeTest(t, `a>a>a;`)
}

func TestBooleanParity(t *testing.T) {
	sanatizeTest(t, `true;`)
	sanatizeTest(t, `false;`)
}

func TestNullParity(t *testing.T) {
	sanatizeTest(t, `null;`)
}

func TestEqualityParity(t *testing.T) {
	sanatizeTest(t, `1==2;`)
	sanatizeTest(t, `a+1!=2+2*3;`)
	sanatizeTest(t, `if(1!=2) res = 200;`)
}

func TestLogicalExpressionParity(t *testing.T) {
	sanatizeTest(t, `a = 1||2==2&&3;`)
	sanatizeTest(t, `a||1!=2+2||3;`)
	sanatizeTest(t, `a||1&&2+2||3;`)
}

func TestUnarityExpressionParity(t *testing.T) {
	sanatizeTest(t, `-1;`)
	sanatizeTest(t, `+1;`)
	sanatizeTest(t, `!true;`)
	sanatizeTest(t, `!!true;`)
}
