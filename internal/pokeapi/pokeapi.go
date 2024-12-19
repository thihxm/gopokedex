package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/thihxm/gopokedex/internal/pokecache"
)

const (
	baseURL = "https://pokeapi.co/api/v2"
)

type LocationAreaDTO struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationAreaDetailsDTO struct {
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

func GetLocationAreaDetails(area string) (LocationAreaDetailsDTO, error) {
	trimmedArea := strings.Trim(area, " ")
	if len(trimmedArea) == 0 {
		return LocationAreaDetailsDTO{}, fmt.Errorf("invalid area")
	}
	locationUrl := baseURL + "/location-area/" + trimmedArea

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
