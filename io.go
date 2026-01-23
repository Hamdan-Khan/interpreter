package main

import (
	"fmt"
	"os"
)

func RunFile(path string){
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file: %v", err.Error())
		return
	}
	fileContet := string(file[:])
	fmt.Printf("%v",fileContet)
}


func RunPrompt (){
	// TODO: implement REPL
	var input string 
	for {
		println("Write script here:\n>>")
		fmt.Scanln(&input)
		if input == "q" {
			fmt.Println("Quitting interpreter")
			break
		}
	}
}
