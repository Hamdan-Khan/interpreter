package main

import (
	"fmt"
	"os"

	"github.com/hamdan-khan/interpreter/parser"
	"github.com/hamdan-khan/interpreter/syntax"
)

func RunFile(path string){
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
	expr := parser.Parse()

	printer := syntax.NewAstPrinter()	
	formatted := printer.Print(expr) //todo: debug
	fmt.Printf("Expression: %v\n", formatted)
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
