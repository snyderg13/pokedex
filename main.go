package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cmdConfig struct {
	Next string
	Prev string
}

type cliCommand struct {
	name        string
	description string
	callback    func(cmdConfig) error
}

var pokeCmds = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
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

func commandExit(cfg cmdConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg cmdConfig) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	fmt.Printf("help:\tDisplays a help message\n")

	// @TODO: figure out how to loop over pokeCmds without compiler errors
	//        about initialization cycles when "help" is in the cli registry
	//        for now, it has been removed from the registry
	for _, cmd := range pokeCmds {
		fmt.Printf("%s:\t%s\n", cmd.name, cmd.description)
	}

	return nil
}

func commandMap(cfg cmdConfig) error {

	return nil
}

func main() {
	var line string
	var words []string
	var worldCfg cmdConfig
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
			// handle "help" command separately as it's not in the cli registry
			// due to initialization cycle compiler errors when in the registry
			if command == "help" {
				commandHelp(worldCfg)
			} else if cmd, ok := pokeCmds[command]; !ok {
				fmt.Printf("Unknown command\n")
			} else if err := cmd.callback(worldCfg); err != nil {
				// @TODO: not sure if below is the best way to do this
				//        it looks gross and is most likely not something
				//        that should be delayed to the user
				fmt.Println(fmt.Errorf("command \"%s\" returned error \"%w\"", cmd.name, err))
			} else {
				// @TODO: other logic to be added if needed; empt for now on purpose
			}
		}
	}
}
