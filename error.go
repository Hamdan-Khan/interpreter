package main

import (
	"fmt"
	"os"
)

func ReportError(lineNumber int, location string, err error){
	fmt.Printf("[line %d] Error %s: %s \n", lineNumber, location, err.Error())
}

func TestError (){
	_, err := os.ReadDir("any")
	ReportError(10, "anywhere", err)
	os.Exit(65)
}

