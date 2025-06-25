package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	if len(text) == 0 {
		return []string{}
	}

	return strings.Fields(strings.ToLower(text))
}

func printUsage() {
	fmt.Println("-----------------------------------------------------")
	fmt.Println("Just type in a message and we'll echo the first word")
	fmt.Println("-----------------------------------------------------")
}

func main() {
	var line string
	var words []string
	inputScanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		inputScanner.Scan()
		err := inputScanner.Err()
		if err != nil {
			fmt.Errorf("inputScanner returned error: %w", err)
		}
		line = inputScanner.Text()
		words = cleanInput(line)
		switch words[0] {
		case "exit":
			fmt.Println("Goodbye!")
			os.Exit(0)
		case "help":
			printUsage()
		default:
			fmt.Printf("\nYour command was: %s\n", words[0])
		}
	}
}
