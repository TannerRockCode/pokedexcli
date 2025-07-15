package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/TannerRockCode/pokedexcli/internal/pokecache"
)

type LocationAreaResponse struct {
	Count    int                  `json:"count"`
	Next     *string              `json:"next"`
	Previous *string              `json:"previous"`
	Results  []LocationAreaResult `json:"results"`
}

type LocationAreaResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func GetLocationAreas(limit Limit, mapEnum int) (LocationAreaResponse, error) {
	pokecache.NewCache(time.Duration(5))
	var locationAreas LocationAreaResponse
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?limit=%d&offset=%d", limit, mapEnum)

	// dat, exists := Cache.Get(url)
	// if exists {
	// 	err := json.Unmarshal(dat, &locationAreas)
	// 	if err != nil {
	// 		return locationAreas, err
	// 	}
	// 	return locationAreas, nil
	// }

	dat, err := Get(url)
	if err != nil {
		return locationAreas, err
	}
	//fmt.Printf("dat: %v", string(dat))
	err = json.Unmarshal(dat, &locationAreas)
	if err != nil {
		return locationAreas, err
	}
	return locationAreas, nil
}

// type LocationArea struct {
// 	ID                   int                   `json:"id"`
// 	Name                 string                `json:"name"`
// 	GameIndex            int                   `json:"game_index"`
// 	EncounterMethodRates []EncounterMethodRate `json:"encounter_method_rates"`
// 	Location             NamedResource         `json:"location"`
// 	Names                []Name                `json:"names"`
// 	PokemonEncounters    []PokemonEncounter    `json:"pokemon_encounters"`
// }

//url := "https://pokeapi.co/api/v2/location-area/?limit=20"
