package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/thihxm/gopokedex/internal/pokecache"
)

const (
	baseURL = "https://pokeapi.co/api/v2"
)

var cache = pokecache.NewCache(5 * time.Minute)

func GetLocationArea(url *string) (LocationAreaDTO, error) {
	locationUrl := baseURL + "/location-area"
	if url != nil {
		locationUrl = *url
	}

	var locationArea LocationAreaDTO
	if cacheEntry, ok := cache.Get(locationUrl); ok {
		err := json.Unmarshal(cacheEntry, &locationArea)
		if err != nil {
			return LocationAreaDTO{}, err
		}
		return locationArea, nil
	}

	res, err := http.Get(locationUrl)
	if err != nil {
		return LocationAreaDTO{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreaDTO{}, err
	}

	cache.Add(locationUrl, data)

	err = json.Unmarshal(data, &locationArea)
	if err != nil {
		return LocationAreaDTO{}, err
	}

	return locationArea, nil
}

func GetLocationAreaDetails(areaOrID string) (LocationAreaDetailsDTO, error) {
	locationUrl := baseURL + "/location-area/" + areaOrID

	var locationAreaDetails LocationAreaDetailsDTO
	if cacheEntry, ok := cache.Get(locationUrl); ok {
		err := json.Unmarshal(cacheEntry, &locationAreaDetails)
		if err != nil {
			return LocationAreaDetailsDTO{}, err
		}
		return locationAreaDetails, nil
	}

	res, err := http.Get(locationUrl)
	if err != nil {
		return LocationAreaDetailsDTO{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreaDetailsDTO{}, err
	}

	cache.Add(locationUrl, data)

	err = json.Unmarshal(data, &locationAreaDetails)
	if err != nil {
		return LocationAreaDetailsDTO{}, err
	}

	return locationAreaDetails, nil
}

func GetPokemon(pokemonNameOrID string) (PokemonDTO, error) {
	locationUrl := baseURL + "/pokemon/" + pokemonNameOrID

	var pokemon PokemonDTO
	if cacheEntry, ok := cache.Get(locationUrl); ok {
		err := json.Unmarshal(cacheEntry, &pokemon)
		if err != nil {
			return PokemonDTO{}, err
		}
		return pokemon, nil
	}

	res, err := http.Get(locationUrl)
	if err != nil {
		return PokemonDTO{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return PokemonDTO{}, err
	}

	cache.Add(locationUrl, data)

	err = json.Unmarshal(data, &pokemon)
	if err != nil {
		return PokemonDTO{}, err
	}

	return pokemon, nil
}
