package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	if len(text) == 0 {
		return []string{}
	}

	return strings.Fields(strings.ToLower(text))
}

func main() {
	fmt.Println("Hello World!")
	cleanInput("ANOTHER TEST STRING    IN    HERE")
	cleanInput("My TeST Strings hereE")
}
