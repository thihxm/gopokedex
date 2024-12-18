package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Pokedex > ")
	for scanner.Scan() {
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		if len(cleanedInput) == 0 {
			fmt.Print("Pokedex > ")
			continue
		}
		firstWord := cleanedInput[0]

		fmt.Printf("Your command was: %s\n", firstWord)
		fmt.Print("Pokedex > ")
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
