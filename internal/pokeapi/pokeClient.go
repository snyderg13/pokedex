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
		err := json.Unmarshal(cacheData, &results)
		if err != nil {
			fmt.Println("failed to unmarshal cache data", err)
		}
		fmt.Println("CLIENT: Cache get was used")
		return results, err //nil
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

	// convert results to []byte and add to the cache
	bytesBody, err := io.ReadAll(res.Body)
	fmt.Println("len(bytesBody): ", len(bytesBody))
	if err != nil {
		fmt.Println(err)
		return LocAreaResp{}, err
	}

	// add data byte slice to cache
	fmt.Println("CLIENT: Cache add was used")
	pokeAPICache.Add(locURL, bytesBody)

	// unmarshal into the LocAreaResp to return to the caller
	err = json.Unmarshal(bytesBody, &results)
	if err != nil {
		fmt.Println("failed to unmarshal cache data", err)
	}

	return results, nil
}
