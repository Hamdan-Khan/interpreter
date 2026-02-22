package main

import (
	"fmt"
	"os"

	"github.com/hamdan-khan/interpreter/interpreter"
	"github.com/hamdan-khan/interpreter/parser"
	"github.com/hamdan-khan/interpreter/token"
)

func RunFile(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file: %v", err.Error())
		return
	}
	fileContet := string(file[:])
	scanner := token.NewScanner(fileContet)
	scanner.Scan()

	tokens := scanner.Tokens
	parser := parser.NewParser(tokens)
	statements, parseErr := parser.Parse()
	if parseErr != nil {
		fmt.Printf("Error parsing syntax: %v\n", parseErr)
		os.Exit(65)
		return
	}

	i := interpreter.NewInterpreter()

	resolver := interpreter.NewResolver(i)
	rErr := resolver.ResolveStmts(statements)
	if rErr != nil {
		fmt.Printf("Error resolving: %v\n", rErr)
		os.Exit(65)
		return
	}

	iError := i.Interpret(statements)
	if iError != nil {
		fmt.Printf("Error evaluating: %v\n", iError)
	}
}

func RunPrompt() {
	// TODO: implement REPL
	var input string
	for {
		println("Write script here:\n>>")
		// todo: replace with bufio.Scanner
		fmt.Scanln(&input)
		if input == "q" {
			fmt.Println("Quitting interpreter")
			break
		}
	}
}
