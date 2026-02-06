package errorHandler

import (
	"fmt"
)

func ReportError(lineNumber int, location string, errorMessage string) error {
	fmt.Printf("[line %d] Error %s: %s \n", lineNumber, location, errorMessage)
	return fmt.Errorf("%s", errorMessage)
}

func TestError (){
	ReportError(10, "anywhere", "Test error")
}
