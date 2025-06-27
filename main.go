package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/snyderg13/pokedex/internal/pokeapi"
)

type cmdConfig struct {
	Next string
	Prev string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*cmdConfig) error
}

var pokeCmds map[string]cliCommand

func initCmds() {
	pokeCmds = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
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
		"mapb": {
			name:        "mapb",
			description: "Displays world locations",
			callback:    commandMapb,
		},
	}
}

// sanitize user input by taking input text
// make it lowercase and split into a slice
func cleanInput(text string) []string {
	if len(text) == 0 {
		return []string{}
	}

	return strings.Fields(strings.ToLower(text))
}

func commandExit(cfg *cmdConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *cmdConfig) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")

	for _, cmd := range pokeCmds {
		fmt.Printf("%s:\t%s\n", cmd.name, cmd.description)
	}

	return nil
}

// @TODO: handle appending the offset and limit query params to end of URL
// @TODO: for adding the ability to go forward and back and setting the
// @TODO: next/prev fields of the cmdConfig param
func commandMap(cfg *cmdConfig) error {
	debug := false

	var results pokeapi.LocAreaResp
	results, err := pokeapi.GetLocationAreas(cfg.Next)
	if err != nil {
		return err
	}

	cfg.Next = results.Next
	cfg.Prev = results.Prev

	if debug {
		fmt.Printf("cfg.Next = %s, cfg.Prev = %s\n", cfg.Next, cfg.Prev)
	}

	return nil
}

// @TODO: split out HTTP GET functionality to be called
// @TODO: by both commandMap and commandMapb
func commandMapb(cfg *cmdConfig) error {
	if len(cfg.Prev) == 0 {
		fmt.Println("You're on the first page")
		return nil
	}

	debug := false

	var results pokeapi.LocAreaResp
	results, err := pokeapi.GetLocationAreas(cfg.Prev)
	if err != nil {
		return err
	}

	cfg.Next = results.Next
	cfg.Prev = results.Prev

	if debug {
		fmt.Printf("cfg.Next = %s, cfg.Prev = %s\n", cfg.Next, cfg.Prev)
	}

	return nil
}

func main() {
	var line string
	var words []string
	worldCfg := cmdConfig{}
	initCmds()
	pokeapi.Init()
	mainDebug := false
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
			} else if err := cmd.callback(&worldCfg); err != nil {
				// @TODO: not sure if below is the best way to do this
				//        it looks gross and is most likely not something
				//        that should be delayed to the user
				fmt.Println(fmt.Errorf("command \"%s\" returned error \"%w\"", cmd.name, err))
			} else {
				// @TODO: other logic to be added if needed
				//        intentionally empty for now on purpose
				if mainDebug {
					fmt.Printf("worldCfg.Next = %s, worldCfg.Prev = %s\n", worldCfg.Next, worldCfg.Prev)
				}
			}
		}
	}
}
