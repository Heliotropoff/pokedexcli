package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Heliotropoff/pokedexcli/internal/cache"
	"github.com/Heliotropoff/pokedexcli/internal/pokeapi"
)

type Config struct {
	Previous string
	Next     string
	Cache    *cache.Cache
}

// var config Config = Config{}
var supportedCommands map[string]cliCommand

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config, args []string) error
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	words := strings.Fields(text)

	return words
}

func commandExit(config *Config, _ []string) error {
	_, err := fmt.Println("Closing the Pokedex... Goodbye!")
	if err != nil {
		return err
	}
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, _ []string) error {
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

func commandMap(config *Config, _ []string) error {
	var callUrl string
	var data pokeapi.LocationsData
	var err error
	if config.Next == "" {
		callUrl = baseLocationsURL
	} else {
		callUrl = config.Next
	}
	if rawData, ok := config.Cache.Get(callUrl); !ok {
		data, err = pokeapi.CallLocations(callUrl)
		if err != nil {
			return err
		}
		rawData, err = json.Marshal(data)
		if err != nil {
			return err
		}
		config.Cache.Add(callUrl, rawData)
	} else {
		if err := json.Unmarshal(rawData, &data); err != nil {
			return err
		}
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

func commandMapb(config *Config, _ []string) error {
	var callUrl string
	var data pokeapi.LocationsData
	var err error
	if config.Previous == "" {
		fmt.Println("you're on the first page")
	} else {
		callUrl = config.Previous
	}
	if rawData, ok := config.Cache.Get(callUrl); !ok {
		data, err = pokeapi.CallLocations(callUrl)
		if err != nil {
			return err
		}
		rawData, err = json.Marshal(data)
		if err != nil {
			return err
		}
		config.Cache.Add(callUrl, rawData)

	} else {
		if err = json.Unmarshal(rawData, &data); err != nil {
			return err
		}
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

func commandExplore(config *Config, args []string) error {
	//TODO: upade with CACHE AND API LOGIC
	callUrl := baseLocationsURL + args[0]
	var data pokeapi.ExploredLocationsData
	var err error
	if rawData, ok := config.Cache.Get(callUrl); !ok {
		data, err = pokeapi.ExploreLocations(callUrl)
		if err != nil {
			return err
		}
		rawData, err = json.Marshal(data)
		if err != nil {
			return err
		}
		config.Cache.Add(callUrl, rawData)
	} else {
		err = json.Unmarshal(rawData, &data)
		if err != nil {
			return err
		}
	}
	for idx := range data.PokemonEncounters {
		fmt.Println(data.PokemonEncounters[idx].Pokemon.Name)
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
		"explore": {
			name:        "explore",
			description: "Lists all pokemon in a given area",
			callback:    commandExplore,
		},
	}
}
