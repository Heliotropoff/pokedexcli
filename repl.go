package main

import (
	"fmt"
	"os"
	"strings"
)

var supportedCommands map[string]cliCommand

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	words := strings.Fields(text)

	return words
}

func commandExit() error {
	_, err := fmt.Println("Closing the Pokedex... Goodbye!")
	if err != nil {
		return err
	}
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	for _, cmd := range supportedCommands {
		_, err := fmt.Printf("%s: %s\n", cmd.name, cmd.description)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	supportedCommands = map[string]cliCommand{
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
	}
}
