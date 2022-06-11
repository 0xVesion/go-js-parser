package main

import (
	"encoding/json"
	"fmt"

	"github.com/0xvesion/go-js-parser/parser"
	"github.com/0xvesion/go-js-parser/tokenizer"
)

func main() {
	ast, err := parser.New(tokenizer.New(`
	1+1;
	2+2*2;
	"Hello World!";
	`)).Parse()
	if err != nil {
		panic(err)
	}

	j, _ := json.MarshalIndent(ast, "", "  ")
	fmt.Println(string(j))
}
