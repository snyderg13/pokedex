package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var pokeCmds = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"test_err": {
		name:        "error",
		description: "test error return handling",
		callback:    commandTestError,
	},
	"help": {
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	},
	"map": {
		name:        "map",
		description: "Displays world locations",
		callback:    commandMap,
	},
}

// sanitize user input by taking input text
// make it lowercase and split into a slice
func cleanInput(text string) []string {
	if len(text) == 0 {
		return []string{}
	}

	return strings.Fields(strings.ToLower(text))
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandTestError() error {
	return fmt.Errorf("test error returning")
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	fmt.Printf("help: Displays a help message\n")
	fmt.Printf("exit: Exit the pokedex\n")

	// @TODO: figure out how to loop over pokeCmds without compiler errors
	//        about references cycles
	// for _, cmd := range pokeCmds {
	// 	fmt.Printf("$s: %s\n", cmd.name, cmd.description)
	// }

	return nil
}

func commandMap() error {

	return nil
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
			command := words[0]
			if cmd, ok := pokeCmds[command]; !ok {
				fmt.Printf("Unknown command\n")
			} else {
				if err := cmd.callback(); err != nil {
					fmt.Printf("Command %s returned error %w\n", cmd.name, err)
				}
			}
		}
	}
}
