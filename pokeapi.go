package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"

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

type LocationAreaPokeEncounter struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon PokemonResourceInfo `json:"pokemon"`
}

type PokemonResourceInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonInfo struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Abilities      []PokemonAbilityPosition
}

type PokemonAbilityPosition struct {
	IsHidden bool                       `json:"is_hidden"`
	Slot     int                        `json:"slot"`
	Ability  PokemonAbilityResourceInfo `json:"ability"`
}

type PokemonAbilityResourceInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func GetLocationAreas(limit Limit, mapEnum int) (LocationAreaResponse, error) {
	var locationAreas LocationAreaResponse
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?limit=%d&offset=%d", limit, mapEnum)

	dat, exists := pokecache.Cache.Get(url)
	if exists {
		err := json.Unmarshal(dat, &locationAreas)
		if err != nil {
			return locationAreas, err
		}
		return locationAreas, nil
	}

	dat, err := Get(url)
	if err != nil {
		return locationAreas, err
	}

	err = json.Unmarshal(dat, &locationAreas)
	if err != nil {
		return locationAreas, err
	}
	pokecache.Cache.Add(url, dat)
	return locationAreas, nil
}

func GetExploreLocationAreas(areaName string) ([]string, error) {
	var la LocationAreaPokeEncounter
	var pokemonList []string
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", areaName)

	//Get data from Cache
	dat, exists := pokecache.Cache.Get(url)
	if exists {
		err := json.Unmarshal(dat, &la)
		if err != nil {
			return pokemonList, err
		}
		for _, encounter := range la.PokemonEncounters {
			pokemonList = append(pokemonList, encounter.Pokemon.Name)
		}
		return pokemonList, nil
	}

	//Get data from endpoint
	dat, err := Get(url)
	if err != nil {
		return pokemonList, err
	}

	err = json.Unmarshal(dat, &la)
	if err != nil {
		return pokemonList, err
	}

	for _, encounter := range la.PokemonEncounters {
		pokemonList = append(pokemonList, encounter.Pokemon.Name)
	}
	pokecache.Cache.Add(url, dat)
	return pokemonList, nil
}

func GetPokemonCatchAttempt(pokemon string) (bool, error) {
	var pI PokemonInfo
	var caught bool = false
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemon)

	//Get data from endpoint
	dat, err := Get(url)
	if err != nil {
		return caught, err
	}

	err = json.Unmarshal(dat, &pI)
	if err != nil {
		return caught, err
	}

	caught = CatchAttempt(pI)
	return caught, nil

}

func CatchAttempt(pi PokemonInfo) bool {
	var caught bool = false
	fmt.Printf("Pokemon: %s - BaseExperience: %d\n", pi.Name, pi.BaseExperience)
	randNum := rand.IntN(pi.BaseExperience)
	if randNum < 32 {
		caught = true
		Pokedex[pi.Name] = pi
	}
	return caught
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
