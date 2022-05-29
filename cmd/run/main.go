package main

import (
	"encoding/json"
	"fmt"

	"github.com/0xvesion/go-parser/parser"
	jsonastfactory "github.com/0xvesion/go-parser/parser/json_ast_factory"
	"github.com/0xvesion/go-parser/tokenizer"
)

func main() {
	ast, err := parser.New(tokenizer.New(`
	1+1;
	2+2*2;
	"Hello World!";
	`), jsonastfactory.New()).Parse()
	if err != nil {
		panic(err)
	}

	j, _ := json.MarshalIndent(ast, "", "  ")
	fmt.Println(string(j))
}
