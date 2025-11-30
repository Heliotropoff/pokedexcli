package main

import (
	"fmt"
	"os"
	"strings"

	pokeapi "github.com/Heliotropoff/pokedexcli/internal"
)

type Config struct {
	Previous string
	Next     string
}

var config Config = Config{}
var supportedCommands map[string]cliCommand

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config) error
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	words := strings.Fields(text)

	return words
}

func commandExit(config *Config) error {
	_, err := fmt.Println("Closing the Pokedex... Goodbye!")
	if err != nil {
		return err
	}
	os.Exit(0)
	return nil
}

func commandHelp(config *Config) error {
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

const baseLocationsURL string = "https://pokeapi.co/api/v2/location-area/"

func commandMap(config *Config) error {
	var callUrl string
	if config.Next == "" {
		callUrl = baseLocationsURL
	} else {
		callUrl = config.Next
	}
	data, err := pokeapi.CallLocations(callUrl)
	if err != nil {
		return err
	}
	config.Next = data.Next
	if data.Previous != nil {
		config.Previous = data.Previous.(string)
	} else {
		config.Previous = ""
	}
	for loc_idx := range data.Results {
		fmt.Println(data.Results[loc_idx].Name)
	}
	return nil
}

func commandMapb(config *Config) error {
	var callUrl string
	if config.Previous == "" {
		fmt.Println("you're on the first page")
	} else {
		callUrl = config.Previous
	}
	data, err := pokeapi.CallLocations(callUrl)
	if err != nil {
		return err
	}
	config.Next = data.Next
	if data.Previous != nil {
		config.Previous = data.Previous.(string)
	} else {
		config.Previous = ""
	}
	for loc_idx := range data.Results {
		fmt.Println(data.Results[loc_idx].Name)
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
		"map": {
			name:        "map",
			description: "Displays the names of next 20 location areas in the Pokemon world.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of previous 20 location areas in the Pokemon world.",
			callback:    commandMapb,
		},
	}
}
