package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
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
	callback    func(*cmdConfig, ...string) error
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
		"explore": {
			name:        "explore",
			description: "Explore an area for pokemon",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon",
			callback:    commandCatch,
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

func commandExit(cfg *cmdConfig, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *cmdConfig, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")

	for _, cmd := range pokeCmds {
		fmt.Printf("%s:\t%s\n", cmd.name, cmd.description)
	}

	return nil
}

func commandMap(cfg *cmdConfig, args ...string) error {
	debug := false

	var results pokeapi.LocAreaResp
	results, err := results.DoGetData(cfg.Next)
	if err != nil {
		return err
	}

	cfg.Next = results.Next
	cfg.Prev = results.Prev

	for _, name := range results.Results {
		fmt.Println(name.Name)
	}

	if debug {
		fmt.Printf("cfg.Next = %s, cfg.Prev = %s\n", cfg.Next, cfg.Prev)
	}

	return nil
}

func commandMapb(cfg *cmdConfig, args ...string) error {
	if len(cfg.Prev) == 0 {
		fmt.Println("You're on the first page")
		return nil
	}

	debug := false

	var results pokeapi.LocAreaResp
	results, err := results.DoGetData(cfg.Prev)
	if err != nil {
		return err
	}

	cfg.Next = results.Next
	cfg.Prev = results.Prev

	for _, name := range results.Results {
		fmt.Println(name.Name)
	}

	if debug {
		fmt.Printf("cfg.Next = %s, cfg.Prev = %s\n", cfg.Next, cfg.Prev)
	}

	return nil
}

// @TODO add test cases for different commands
func commandExplore(cfg *cmdConfig, args ...string) error {
	fmt.Println("len(args) = ", len(args))
	fmt.Println("args = ", args)
	if len(args) == 0 {
		return fmt.Errorf("not enough args, expected <location_name>")
	}
	fmt.Printf("Exploring %s...\n", args[0])

	var results pokeapi.LocationDetails
	results, err := results.DoGetData(args[0])
	if err != nil {
		fmt.Println("Exp: get data ret: ", err)
		return err
	}

	for _, p := range results.PokemonList {
		fmt.Println(p.Pokemon.Name)
	}

	// note that cfg.next and cfg.prev are not updated
	// since the user only chose to explore an area;
	// next and prev are only really used for map and mapb
	// up to this point in development

	return nil
}

func commandCatch(cfg *cmdConfig, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Did not provide pokemon name")
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", args[0])

	randVal := rand.Float32()
	fmt.Printf("Randomly generated value = %f == %v\n", randVal, randVal)
	// @TODO: fetch the pokemon data and use the chance
	//        plus a random number to determine if the
	//        user's catch attempt is successful
	//        use base_experience field on some sort of scaled chance range to determine success
	// psuedo code:

	var results pokeapi.PokemonStats
	results, err := results.DoGetData(args[0])
	if err != nil {
		fmt.Println("Exp: get data ret: ", err)
		return err
	}

	base_exp := results.BaseExperience
	fmt.Printf("%s base exp is %d\n", args[0], base_exp)

	if base_exp > 600 {
		// 5 % chance to catch
	} else if base_exp > 500 {
		// 7.5 % chance to catch
	} else if base_exp > 400 {
		// 10 % chance to catch
	} else if base_exp > 200 {
		// 15 % chance to catch
	} else if base_exp > 100 {
		// 20 % chance to catch
	} else {
		// 25 % chance to catch
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
			args := words[1:]
			if cmd, ok := pokeCmds[command]; !ok {
				fmt.Printf("Unknown command: %s\n", command)
			} else if err := cmd.callback(&worldCfg, args...); err != nil {
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
