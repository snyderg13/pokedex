package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/snyderg13/pokedex/internal/pokecache"
)

const (
	pokeAPIBaseURL       = "https://pokeapi.co/api/v2/"
	locationAreaEndpoint = "https://pokeapi.co/api/v2/location-area/"
	cacheReapRate        = 5 * time.Second
)

var pokeAPICache pokecache.Cache

func Init() {
	pokeAPICache = pokecache.NewCache(cacheReapRate)
	fmt.Println("pokeAPICache: ", pokeAPICache, &pokeAPICache)
}

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

	var results LocAreaResp
	cacheData, found := pokeAPICache.Get(locURL)
	if found {
		fmt.Println("len(cacheData): ", len(cacheData))
		fmt.Println("cache item data: ", cacheData)
		err := json.Unmarshal(cacheData, &results)
		if err != nil {
			fmt.Println("failed to unmarshal cache data", err)
		}
	} else {
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

		// var results LocAreaResp
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&results)
		if err != nil {
			fmt.Println(fmt.Errorf("json decode returned %w", err))
			return LocAreaResp{}, fmt.Errorf("json decode returned %w", err)
		}

		// @TODO: might need to split json/data handling out of GET handling
		// @TODO: since below call adds empty data since above decode
		// @TODO: operation empties the res.Body buffer
		// or do I need to just (re)marshal data and then unmarshal it
		// to get it again or create a copy of it to add to cache?

		// convert results to []byte and add to the cache
		bytesBody, err := io.ReadAll(res.Body)
		fmt.Println("len(bytesBody): ", len(bytesBody))
		if err != nil {
			fmt.Println(err)
			return LocAreaResp{}, err
		}
		pokeAPICache.Add(locURL, bytesBody)
	}

	for _, name := range results.Results {
		fmt.Println(name.Name)
	}

	return results, nil
}
