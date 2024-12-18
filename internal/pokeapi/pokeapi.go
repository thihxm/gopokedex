package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
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

func GetLocationArea(offset *string) (LocationAreaDTO, error) {
	locationUrl := baseURL + "/location-area?limit=20"
	if offset != nil {
		locationUrl += "&offset=" + *offset
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

	var locationArea LocationAreaDTO
	err = json.Unmarshal(data, &locationArea)
	if err != nil {
		return LocationAreaDTO{}, err
	}

	if locationArea.Next != nil {
		parsedURL, _ := url.Parse(*locationArea.Next)
		if parsedURL != nil {
			q := parsedURL.Query()
			offset := q.Get("offset")
			if offset != "" {
				locationArea.Next = &offset
			} else {
				locationArea.Next = nil
			}
		}
	} else {
		locationArea.Next = nil
	}

	if locationArea.Previous != nil {
		parsedURL, _ := url.Parse(*locationArea.Previous)
		if parsedURL != nil {
			q := parsedURL.Query()
			offset := q.Get("offset")
			if offset != "" {
				locationArea.Previous = &offset
			} else {
				locationArea.Previous = nil
			}
		}
	} else {
		locationArea.Previous = nil
	}

	return locationArea, nil
}
