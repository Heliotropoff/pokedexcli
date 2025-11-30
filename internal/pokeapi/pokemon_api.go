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
