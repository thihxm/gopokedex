package main

import (
	"bufio"
	"fmt"
	"math/rand"
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

const (
	PokeballBaseRate int = 255
	MinCatchRate     int = 75
)

var commands map[string]cliCommand
var cfg = config{
	Next:     nil,
	Previous: nil,
}
var pokedex = map[string]pokeapi.PokemonDTO{}

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
		"catch": {
			name:        "catch",
			description: "Tries to catch a Pokemon\n" + "Usage: catch <Pokemon name>",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspects a caught Pokemon\n" + "Usage: inspect <Pokemon name>",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays the caught Pokemon",
			callback:    commandPokedex,
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

func commandCatch(config *config, params ...string) error {
	if len(params) == 0 {
		return fmt.Errorf("missing Pokemon name")
	}
	pokemonName := params[0]

	pokemon, err := pokeapi.GetPokemon(pokemonName)
	if err != nil {
		return fmt.Errorf("failed to get Pokemon (%s): %w", pokemonName, err)
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	catchRate := (rand.Intn(PokeballBaseRate) * 100) / pokemon.BaseExperience

	if catchRate >= MinCatchRate {
		pokedex[pokemonName] = pokemon
		fmt.Printf("%s was caught!\n", pokemonName)
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}

func commandInspect(config *config, params ...string) error {
	if len(params) == 0 {
		return fmt.Errorf("missing Pokemon name")
	}
	pokemonName := params[0]

	pokemon, ok := pokedex[pokemonName]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf(" - %s\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(config *config, params ...string) error {
	if len(pokedex) == 0 {
		fmt.Println("you have not caught any Pokemon")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for pokemon, _ := range pokedex {
		fmt.Printf(" - %s\n", pokemon)
	}

	return nil
}
