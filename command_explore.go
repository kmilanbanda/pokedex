package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io"
	"bytes"
)

func commandExplore(c *config, args []string) error {
	baseUrl := "https://pokeapi.co/api/v2/location-area/"
	fullUrl := baseUrl + args[0]


	var body []byte
	var isCached bool
	if body, isCached = Cache.Get(fullUrl); !isCached { 
		res, err := http.Get(fullUrl)
		if err != nil {
			fmt.Errorf("Error getting response: %v", err)
			return err
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		Cache.Add(fullUrl, body)
	}

	var encounters AreaEncounters 
	decoder := json.NewDecoder(bytes.NewReader(body))
	if err := decoder.Decode(&encounters); err != nil {
		return err
	}
	
	fmt.Printf("Exploring %s...\nFound Pokemon:\n", args[0])
	pokemon := encounters.Encounters
	for i := 0; i < len(pokemon); i++ {
		fmt.Printf(" - %s\n", pokemon[i].Pokemon.Name)
	}

	return nil	
}
