package main

import (
	"fmt"
	"os"

	"github.com/hamdan-khan/interpreter/scripts"
)

func main(){
	scripts.GenerateAst()
	argsLen := len(os.Args)

	if argsLen > 2 {
		fmt.Println("Invalid arguments. Usage: interpreter [file]")
	} else {
		fmt.Println("Interpreter starting...")
		if argsLen == 2 {
			RunFile(os.Args[1])
		} else {
			RunPrompt()
		}
	}
}