package main

import (
	"fmt"
	"os"
)

func ReportError(lineNumber int, location string, errorMessage string){
	fmt.Printf("[line %d] Error %s: %s \n", lineNumber, location, errorMessage)
	os.Exit(65)
}

func TestError (){
	ReportError(10, "anywhere", "Test error")
}
