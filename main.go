package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/thihxm/gopokedex/internal/pokeapi"
)

type config struct {
	Next     *string
	Previous *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *config, params ...string) error
}

var commands map[string]cliCommand
var cfg = config{
	Next:     nil,
	Previous: nil,
}

func main() {
	commands = map[string]cliCommand{
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
			description: "Displays the map of the region",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous page of the map of the region",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explores a location area\n" + "Usage: explore <area>",
			callback:    commandExplore,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Pokedex > ")
	for scanner.Scan() {
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		if len(cleanedInput) == 0 {
			fmt.Print("Pokedex > ")
			continue
		}

		command := cleanedInput[0]
		params := cleanedInput[1:]

		if cmd, ok := commands[command]; ok {
			if err := cmd.callback(&cfg, params...); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}

		fmt.Print("Pokedex > ")
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit(config *config, params ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *config, params ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage:\n\n")
	for _, cmd := range commands {
		fmt.Printf("%s, %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(config *config, params ...string) error {
	if config.Next == nil && config.Previous != nil {
		fmt.Println("you're on the last page")
		return nil
	}

	locationArea, err := pokeapi.GetLocationArea(config.Next)
	if err != nil {
		return fmt.Errorf("failed to get location area: %w", err)
	}

	config.Next = locationArea.Next
	config.Previous = locationArea.Previous

	for _, location := range locationArea.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapb(config *config, params ...string) error {
	if config.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	locationArea, err := pokeapi.GetLocationArea(config.Previous)
	if err != nil {
		return fmt.Errorf("failed to get location area: %w", err)
	}

	config.Next = locationArea.Next
	config.Previous = locationArea.Previous

	for _, location := range locationArea.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandExplore(config *config, params ...string) error {
	if len(params) == 0 {
		return fmt.Errorf("missing area")
	}
	area := params[0]

	locationAreaDetails, err := pokeapi.GetLocationAreaDetails(area)
	if err != nil {
		return fmt.Errorf("failed to get location area (%s) details: %w", area, err)
	}

	fmt.Printf("Exploring %s...\n", locationAreaDetails.Name)
	fmt.Println("Found Pokemon:")
	for _, pokemonEncounters := range locationAreaDetails.PokemonEncounters {
		fmt.Printf("- %s\n", pokemonEncounters.Pokemon.Name)
	}

	return nil
}
