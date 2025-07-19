package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Limit int
type cliCommand struct {
	name        string
	description string
	callback    func(interface{}) error
}

var Commands map[string]cliCommand
var MapEnum int
var ExploreEnum int
var Pokedex map[string]PokemonInfo

func init() {
	MapEnum = -1
	Pokedex = make(map[string]PokemonInfo, 200)

	Commands = map[string]cliCommand{
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
			description: "Displays names of next 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "bmap",
			description: "Displays names of previous 20 location areas in the Pokemon world",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Displays names of Pokemon found in provided location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Displays attempt to catch provided pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Display information about provided pokemon in inventory",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays names of all pokemon in pokedex inventory",
			callback:    commandPokedex,
		},
	}
}

func repl() {
	var args string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		if len(words) > 1 {
			args = words[1]
		}

		command, ok := Commands[commandName]
		if !ok {
			fmt.Println("Invalid command")
			continue
		}

		command.callback(args)
		//fmt.Printf("Your command was: %s\n", words[0])
	}
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	words := strings.Fields(text)
	return words
}

func commandExit(arg interface{}) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(arg interface{}) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, cmd := range Commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println("")
	return nil
}

func commandMap(arg interface{}) error {
	MapEnum++
	las, err := GetLocationAreas(Limit(20), MapEnum)

	if err != nil {
		return err
	}

	for _, la := range las.Results {
		fmt.Printf("%s\n", la.Name)
	}
	fmt.Printf("MapEnum: %d\n", MapEnum)
	return nil
}

func commandMapB(arg interface{}) error {
	if MapEnum > 0 {
		MapEnum--
	}
	las, err := GetLocationAreas(Limit(20), MapEnum)
	if err != nil {
		return err
	}
	fmt.Printf("MapEnum: %d\n", MapEnum)

	for _, la := range las.Results {
		fmt.Printf("%s\n", la.Name)
	}
	return nil
}

func commandExplore(arg interface{}) error {
	locationArea, ok := arg.(string)
	if !ok {
		return errors.New("explore command requires string as argument. Argument provided was not a string")
	}

	fmt.Printf("Exploring %s...\n", locationArea)
	pokemonNames, err := GetExploreLocationAreas(locationArea)
	if err != nil {
		fmt.Printf("An error occurred on http request to explore location area: %v\n", err)
		return err
	}

	fmt.Printf("Found Pokemon: \n")
	for _, name := range pokemonNames {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}

func commandCatch(arg interface{}) error {
	pokemon, ok := arg.(string)
	if !ok {
		return errors.New("catch command requires string as argument. Argument provided was not a string")
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)
	catchAttempt, err := GetPokemonCatchAttempt(pokemon)
	if catchAttempt {
		fmt.Printf("%s was caught!\n", pokemon)
		fmt.Printf("You may now inspect %s with the inspect command.", pokemon)
	} else {
		fmt.Printf("%s escaped!\n", pokemon)
	}
	return err
}

func commandInspect(arg interface{}) error {
	pokemon, ok := arg.(string)
	if !ok {
		return errors.New("inspect command requires string as argument. Argument provided was not a string")
	}
	pokemonInfo, err := InspectPokemon(pokemon)
	if err != nil {
		fmt.Printf("Pokemon has not been caught before\n")
		return err
	}

	fmt.Printf("Name: %s\nHeight: %d, \nWeight: %d,\n", pokemonInfo.Name, pokemonInfo.Height, pokemonInfo.Weight)
	fmt.Printf("Stats: \n")
	for _, val := range pokemonInfo.Stats {
		fmt.Printf("  -%s: %d\n", val.Stat.Name, val.BaseStat)
	}
	fmt.Printf("Types: \n")
	for _, val := range pokemonInfo.Types {
		fmt.Printf("  - %s\n", val.Type.Name)
	}

	return nil
}

func commandPokedex(arg interface{}) error {
	fmt.Println("Your Pokedex:")
	if len(Pokedex) == 0 {
		fmt.Println("Pokedex empty")
	} else {
		for _, pokemon := range Pokedex {
			fmt.Printf(" - %s\n", pokemon.Name)
		}
	}
	return nil
}
