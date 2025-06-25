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

func commandExit() {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
}

func main() {
	var line string
	var words []string
	inputScanner := bufio.NewScanner(os.Stdin)

	for {
		// prompt user for input
		fmt.Print("Pokedex > ")
		inputScanner.Scan()

		// check for scanner errors
		err := inputScanner.Err()
		if err != nil {
			fmt.Errorf("inputScanner returned error: %w", err)
		}

		// get text input from the user
		line = inputScanner.Text()

		// handle scenario where user only hits the enter key
		if len(line) == 0 {
			fmt.Print()
		} else {
			// clean up the input and act on the commands
			words = cleanInput(line)
			if words == nil {
				fmt.Println()
			}
			switch words[0] {
			case "exit":
				commandExit()
			case "help":
				printUsage()
			default:
				fmt.Printf("Unknown command\n")
			}
		}
	}
}
