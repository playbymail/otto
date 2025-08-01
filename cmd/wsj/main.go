// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/playbymail/otto/wsj/ast"
	"github.com/playbymail/otto/wsj/parser"
)

func main() {
	input := `print(5);`

	result, err := parser.Parse("", []byte(input))
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}

	prog, ok := result.(*ast.Program)
	if !ok {
		fmt.Fprintf(os.Stderr, "unexpected AST type: %T\n", result)
		os.Exit(1)
	}
	spew.Dump(prog)

	fmt.Println("Parse successful!")
	for i, stmt := range prog.Statements {
		fmt.Printf("Statement %d: %#v\n", i+1, stmt)
	}
}
