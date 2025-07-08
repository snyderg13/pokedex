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
	pokemonEndpoint      = "https://pokeapi.co/api/v2/pokemon/"
	cacheReapRate        = 10 * time.Second
)

var pokeAPIDebug = false
var pokeAPICache pokecache.Cache

func Init() {
	pokeAPICache = pokecache.NewCache(cacheReapRate)
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

type PokeFetch interface {
	DoGetData(string) (any, error)
}

func (l LocAreaResp) DoGetData(url string) (LocAreaResp, error) {
	if url == "" {
		url = locationAreaEndpoint
	}

	var results LocAreaResp
	cacheData, found := pokeAPICache.Get(url)
	if found {
		err := json.Unmarshal(cacheData, &results)
		if err != nil {
			fmt.Println("failed to unmarshal cache data", err)
		}

		if pokeAPIDebug {
			fmt.Println("CLIENT: Cache get was used")
		}
		return results, err
	}

	res, err := http.Get(url)
	if err != nil {
		// @TODO cleanup below lines
		fmt.Println("http req failed")
		fmt.Println(fmt.Errorf("http req failed: %w", err))
		return results, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return results, fmt.Errorf("status code (%d) > 299", res.StatusCode)
	} else if res.StatusCode != 200 {
		fmt.Printf("status code (%d) != 200\n", res.StatusCode)
	}

	// convert results to []byte and add to the cache
	bytesBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return results, err
	}

	// add data byte slice to cache
	pokeAPICache.Add(url, bytesBody)
	if pokeAPIDebug {
		fmt.Println("CLIENT: Cache add was used")
	}

	// unmarshal into the LocAreaResp to return to the caller
	err = json.Unmarshal(bytesBody, &results)
	if err != nil {
		fmt.Println("failed to unmarshal cache data", err)
	}

	return results, nil
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

		if pokeAPIDebug {
			fmt.Println("CLIENT: Cache get was used")
		}
		return results, err
	}

	res, err := http.Get(locURL)
	if err != nil {
		// @TODO cleanup below lines
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
	if err != nil {
		fmt.Println(err)
		return LocAreaResp{}, err
	}

	// add data byte slice to cache
	pokeAPICache.Add(locURL, bytesBody)
	if pokeAPIDebug {
		fmt.Println("CLIENT: Cache add was used")
	}

	// unmarshal into the LocAreaResp to return to the caller
	err = json.Unmarshal(bytesBody, &results)
	if err != nil {
		fmt.Println("failed to unmarshal cache data", err)
	}

	return results, nil
}

type PokemonEncounters struct {
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
}

type LocationDetails struct {
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
	PokemonList []PokemonEncounters `json:"pokemon_encounters"`
}

// @TODO add test cases for different commands
func (l LocationDetails) DoGetData(locName string) (LocationDetails, error) {
	url := locationAreaEndpoint + locName + "/"
	var results LocationDetails
	cacheData, found := pokeAPICache.Get(url)
	if found {
		err := json.Unmarshal(cacheData, &results)
		if err != nil {
			fmt.Println("failed to unmarshal cache data", err)
		}

		if pokeAPIDebug {
			fmt.Println("CLIENT: Cache get was used")
		}
		return results, err
	}

	res, err := http.Get(url)
	if err != nil {
		// @TODO cleanup below lines
		fmt.Println("http req failed")
		fmt.Println(fmt.Errorf("http req failed: %w", err))
		return results, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return results, fmt.Errorf("status code (%d) > 299", res.StatusCode)
	} else if res.StatusCode != 200 {
		fmt.Printf("status code (%d) != 200\n", res.StatusCode)
	}

	// convert results to []byte and add to the cache
	bytesBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return results, err
	}

	// add data byte slice to cache
	pokeAPICache.Add(url, bytesBody)
	if pokeAPIDebug {
		fmt.Println("CLIENT: Cache add was used")
	}

	// unmarshal into the LocAreaResp to return to the caller
	err = json.Unmarshal(bytesBody, &results)
	if err != nil {
		fmt.Println("failed to unmarshal cache data", err)
	}

	return results, nil
}

type PokemonStats struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height    int `json:"height"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order        any `json:"order"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []struct {
		Abilities []struct {
			Ability  any  `json:"ability"`
			IsHidden bool `json:"is_hidden"`
			Slot     int  `json:"slot"`
		} `json:"abilities"`
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
	} `json:"past_abilities"`
	PastTypes []any `json:"past_types"`
	Species   struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       any    `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  any    `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      any    `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale any    `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default"`
				FrontFemale  any    `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string `json:"back_default"`
				BackFemale       any    `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           string `json:"back_default"`
					BackShiny             string `json:"back_shiny"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					BackTransparent       string `json:"back_transparent"`
					FrontDefault          string `json:"front_default"`
					FrontShiny            string `json:"front_shiny"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
					FrontTransparent      string `json:"front_transparent"`
				} `json:"crystal"`
				Gold struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"gold"`
				Silver struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default"`
						BackFemale       any    `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  any    `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      any    `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale any    `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

func (p PokemonStats) DoGetData(pokemonName string) (PokemonStats, error) {
	url := pokemonEndpoint + pokemonName + "/"
	if pokeAPIDebug {
		fmt.Println("url = ", url)
	}
	var results PokemonStats
	cacheData, found := pokeAPICache.Get(url)
	if found {
		err := json.Unmarshal(cacheData, &results)
		if err != nil {
			fmt.Println("failed to unmarshal cache data", err)
		}

		if pokeAPIDebug {
			fmt.Println("CLIENT: Cache get was used")
		}
		return results, err
	}

	res, err := http.Get(url)
	if err != nil {
		// @TODO cleanup below lines
		fmt.Println("http req failed")
		fmt.Println(fmt.Errorf("http req failed: %w", err))
		return results, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return results, fmt.Errorf("status code (%d) > 299", res.StatusCode)
	} else if res.StatusCode != 200 {
		fmt.Printf("status code (%d) != 200\n", res.StatusCode)
	}

	// convert results to []byte and add to the cache
	bytesBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return results, err
	}

	// add data byte slice to cache
	pokeAPICache.Add(url, bytesBody)
	if pokeAPIDebug {
		fmt.Println("CLIENT: Cache add was used")
	}

	// unmarshal into the LocAreaResp to return to the caller
	err = json.Unmarshal(bytesBody, &results)
	if err != nil {
		fmt.Println("failed to unmarshal cache data", err)
	}

	return results, nil
}
