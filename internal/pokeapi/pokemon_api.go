package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

type LocationsData struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func CallLocations(url string) (LocationsData, error) {

	res, err := http.Get(url)
	if err != nil {
		return LocationsData{}, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationsData{}, err
	}
	defer res.Body.Close()
	locData := LocationsData{}
	err = json.Unmarshal(body, &locData)
	if err != nil {
		return LocationsData{}, nil
	}
	return locData, nil

}

// TODO: need an interface to make one function do all calls
type ExploredLocationsData struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func ExploreLocations(url string) (ExploredLocationsData, error) {
	res, err := http.Get(url)
	if err != nil {
		return ExploredLocationsData{}, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ExploredLocationsData{}, err
	}
	defer res.Body.Close()
	locData := ExploredLocationsData{}
	err = json.Unmarshal(body, &locData)
	if err != nil {
		return ExploredLocationsData{}, nil
	}
	return locData, nil

}
