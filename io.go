package main

import (
	"fmt"
	"os"

	"github.com/hamdan-khan/interpreter/syntax"
	"github.com/hamdan-khan/interpreter/token"
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

	// test printing
	as := syntax.Binary{
		Left: &syntax.Unary{
			Operator: token.Token{
				TokenType: token.MINUS,
				Lexeme: "-",
				LineNumber: 10,
				Literal: nil,
			},
			Right: &syntax.Literal{Value: 123},
		},
		Operator: token.Token{
			TokenType: token.STAR,
			Lexeme: "*",
			LineNumber: 10,
			Literal: nil,
		},
		Right: &syntax.Grouping{
			Expression: &syntax.Literal{Value: 34},
		},
	}
	printer := syntax.NewAstPrinter()
	
	formatted := printer.Print(&as)
	fmt.Printf("%v\n", formatted)
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
