package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Limit int
type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var Commands map[string]cliCommand
var MapEnum int

func init() {
	MapEnum = -1

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
		//args := words[1:]

		command, ok := Commands[commandName]
		if !ok {
			fmt.Println("Invalid command")
			continue
		}
		command.callback()
		//fmt.Printf("Your command was: %s\n", words[0])
	}
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	words := strings.Fields(text)
	return words
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, cmd := range Commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println("")
	return nil
}

func commandMap() error {
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

func commandMapB() error {
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
