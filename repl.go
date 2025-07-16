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

func init() {
	MapEnum = -1
	ExploreEnum = -1

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
	}
}

func repl() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}
		commandName := words[0]
		args := words[1]

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
	fmt.Printf("Entering commandExplore logic with arg: %v\n", arg)
	locationArea, ok := arg.(string)
	if !ok {
		fmt.Printf("commandExplore argument not string\n")
		return errors.New("explore command requires string as argument. Argument provided was not a string")
	}
	ExploreEnum++
	byteArr, err := GetExploreLocationAreas(Limit(20), ExploreEnum, locationArea)
	if err != nil {
		fmt.Printf("An error occurred on http request to explore location area: %v\n", err)
		return err
	}
	fmt.Println(string(byteArr))
	return nil
}
