package main

import (
	"encoding/json"
	"fmt"

	"github.com/0xvesion/go-js-parser/parser"
	"github.com/0xvesion/go-js-parser/tokenizer"
)

func main() {
	ast, err := parser.New(tokenizer.New(`function test(a, b, c) {
		result = 123;
	}`)).Parse()
	if err != nil {
		panic(err)
	}

	j, _ := json.MarshalIndent(ast, "", "  ")
	fmt.Println(string(j))
}
