package main

import (
	"fmt"
	"os"

	"github.com/hamdan-khan/interpreter/interpreter"
	"github.com/hamdan-khan/interpreter/parser"
)

func RunFile(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file: %v", err.Error())
		return
	}
	fileContet := string(file[:])
	scanner := NewScanner(fileContet)
	scanner.Scan()
	for t := range scanner.tokens {
		fmt.Printf("%v\n", scanner.tokens[t].Lexeme)
	}

	tokens := scanner.tokens
	parser := parser.NewParser(tokens)
	statements, parseErr := parser.Parse()
	if parseErr != nil {
		fmt.Printf("Error parsing syntax: %v\n", parseErr)
		os.Exit(65)
		return;
	}

	// printer := syntax.NewAstPrinter()	
	// formatted, printErr := printer.Print(expr)
	// if printErr != nil {
	// 	fmt.Printf("Error printing expression: %v\n", err)
	// 	os.Exit(65)
	// 	return
	// }
	// fmt.Printf("Expression: %v\n", formatted)

	interpreter := interpreter.NewInterpreter()
	iError := interpreter.Interpret(statements)
	if iError != nil {
		fmt.Printf("Error evaluating: %v\n", iError)
	}
}

func RunPrompt (){
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
