// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package main implements a WSJ script runner
package main

import (
	"flag"
	"fmt"
	"github.com/maloquacious/wxx"
	"github.com/playbymail/otto"
	"github.com/playbymail/otto/wjs/lexer"
	"github.com/playbymail/otto/wjs/parser"
	"github.com/playbymail/otto/wjs/vm"
	"os"
	"runtime/debug"
	"strings"
)

var (
	debugMode = false
)

func main() {
	flag.BoolVar(&debugMode, "debug", debugMode, "enable debugging mode")
	showBuildInfo := flag.Bool("build-info", false, "show version with commit and exit")
	showVersion := flag.Bool("version", false, "show version and exit")
	flag.Parse()

	if showVersion != nil && *showVersion {
		fmt.Printf("%s\n", otto.Version().Short())
		os.Exit(0)
	} else if showBuildInfo != nil && *showBuildInfo {
		fmt.Printf("otto %s\nwxx  %s\n", otto.Version().String(), wxx.Version())
		os.Exit(0)
	}

	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Usage: wjs [--debug] [--version] <script.wjs>")
		fmt.Println("   or: wjs [--debug] <WJS statement>")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  wjs myscript.wjs")
		fmt.Println("  wjs 'print(5)'")
		fmt.Println("  wjs --version")
		os.Exit(1)
	}

	input := args[0]

	var filename string // the filename to use for error reporting

	// TODO: users will be using the shebang ("#!") to execute scripts, do we need to do anything special to handle it?

	// if the argument looks like a `.wjs` script name, then try to load the script
	if strings.HasSuffix(input, ".wjs") {
		data, err := os.ReadFile(input)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", input, err)
			os.Exit(1)
		}
		input = string(data)
		filename = args[0]
	} else { // Treat as direct statement - join all args
		input = strings.Join(args, " ")
	}

	if debugMode {
		fmt.Printf("Executing: %s\n", input)
		fmt.Println("---")
	}

	executeCode(filename, input)
}

// filename will only be set when running from a script
func executeCode(filename, input string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error:", r)
			if debugMode {
				fmt.Println("--- Stack Trace ---")
				debug.PrintStack()
			}
			os.Exit(1)
		}
	}()

	// tokenize the input
	l := lexer.New(filename, input)
	tokens := l.AllTokens()

	if debugMode {
		fmt.Println("Tokens:")
		for i, tok := range tokens {
			if tok.Type == lexer.EOF {
				fmt.Printf("%3d: %s\n", i+1, tok)
			} else {
				fmt.Printf("%3d: %s at %d:%d\n", i+1, tok, tok.Pos.Line, tok.Pos.Column)
			}
		}
		fmt.Println("---")
	}

	p := parser.New(tokens)
	prog := p.ParseProgram()

	if debugMode {
		fmt.Printf("AST: %d statements\n", len(prog.Stmts))
		fmt.Println("---")
	}

	// TODO: if we're going to check semantics, check them here

	svm := vm.New(filename)
	if err := svm.Execute(prog); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if debugMode {
		fmt.Println("Execution completed successfully")
	}
}
