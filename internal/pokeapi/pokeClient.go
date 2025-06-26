package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	pokeAPIBaseURL       = "https://pokeapi.co/api/v2/"
	locationAreaEndpoint = "https://pokeapi.co/api/v2/location-area/"
)

// {
//   "count": 1089,
//   "next": "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20",
//   "previous": null,
//   "results": [
//     {
//       "name": "canalave-city-area",
//       "url": "https://pokeapi.co/api/v2/location-area/1/"
//     },
//     {

type locArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocAreaResp struct {
	Count   int       `json:"count"`
	Next    string    `json:"next"`
	Prev    string    `json:"previous"`
	Results []locArea `json:"results"`
}

func GetLocationAreas(locURL string) (LocAreaResp, error) {
	if locURL == "" {
		locURL = locationAreaEndpoint
	}

	res, err := http.Get(locURL)
	if err != nil {
		fmt.Println("http req failed")
		fmt.Println(fmt.Errorf("http req failed: %w", err))
		return LocAreaResp{}, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return LocAreaResp{}, fmt.Errorf("status code (%d) > 299", res.StatusCode)
	} else if res.StatusCode != 200 {
		fmt.Printf("status code (%d) != 200\n", res.StatusCode)
	}

	var results LocAreaResp
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&results)
	if err != nil {
		fmt.Println(fmt.Errorf("json decode returned %w", err))
		return LocAreaResp{}, fmt.Errorf("json decode returned %w", err)
	}

	for _, name := range results.Results {
		fmt.Println(name.Name)
	}

	return results, nil
}
