package main

//DRY it
import (
	"encoding/json"
	"fmt"
	"math/rand"
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

const basePokemonUrl string = "https://pokeapi.co/api/v2/pokemon/"

var Pokedex map[string]pokeapi.Pokemon

func commandCatch(config *Config, args []string) error {
	pokeName := args[0]
	callUrl := basePokemonUrl + pokeName
	data, err := pokeapi.GetPokemonData(callUrl)
	if err != nil {
		return err
	}
	basePokeExp := float64(data.BaseExperience)
	//lowestExp := 36.0   // Sunkern
	highestExp := 608.0 //they say it's Blissey
	currentPokemonChallengeLevel := basePokeExp / highestExp
	easyPokemonChance := 1.0
	hardPokemonChance := 0.01
	challengeRange := easyPokemonChance - hardPokemonChance
	currentPokemonChance := easyPokemonChance - currentPokemonChallengeLevel*challengeRange
	fmt.Printf("Throwing a Pokeball at %s...\n", pokeName)
	playersChance := rand.Float64()
	if playersChance < currentPokemonChance {
		Pokedex[pokeName] = data
		fmt.Printf("%s was caught!\n", pokeName)
	} else {
		fmt.Printf("%s escaped!\n", pokeName)
	}
	return nil
}

func commandInspect(config *Config, args []string) error {
	name := args[0]
	data, ok := Pokedex[name]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\n", data.Name, data.Height, data.Weight)
	fmt.Println("Stats:")
	for idx := range data.Stats {
		name := data.Stats[idx].Stat.Name
		value := data.Stats[idx].BaseStat
		fmt.Printf("	-%s: %d\n", name, value)
	}
	fmt.Println("Types:")
	for idx := range data.Types {
		name := data.Types[idx].Type.Name
		fmt.Printf("	- ", name)
	}

	return nil
}

func commandPokedex(config *Config, _ []string) error {
	if len(Pokedex) == 0 {
		fmt.Println("Your Pokedex is empty!")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for key, _ := range Pokedex {
		fmt.Println("	- %s", key)
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
		"catch": {
			name:        "catch",
			description: "Try to catch a pokenmon. On success it is added to you Pokedex",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Check the Pokedex, if pokemon was caught print it's stats",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists all the caught pokemons",
			callback:    commandPokedex,
		},
	}
}
