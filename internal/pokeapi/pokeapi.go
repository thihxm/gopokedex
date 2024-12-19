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

type LocationAreaDTO struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

var cache = pokecache.NewCache(5 * time.Second)

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
