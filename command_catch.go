package main

import  (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
	"bytes"
	"math/rand"
)

func commandCatch(c *config, args []string) error {	
	baseUrl :=  "https://pokeapi.co/api/v2/pokemon/"
	fullUrl := baseUrl + args[0]

	res, err := http.Get(fullUrl)
	if err != nil {
		fmt.Errorf("Error getting response: %v", err)
		return err
	}
	defer res.Body.Close()
	
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Errorf("Error reading body: %v", err)
		return err
	}
	
	var pokemon Pokemon
	decoder := json.NewDecoder(bytes.NewReader(body))
	if err := decoder.Decode(&pokemon); err != nil {
		fmt.Errorf("Error unmarshaling data: %v", err)
		return err
	}
	
	fmt.Printf("Throwing a Pokeball at %s...\n", args[0])

	maxRoll := 1000
	master := false
	if len(args) > 1 {
		switch args[1]  {
		case "master":
			master = true
		case "ultra":
			maxRoll = 3000
		case "great":
			maxRoll = 2000
		default:
			maxRoll = 1000
		}
	}

	if roll := rand.Intn(maxRoll); roll > pokemon.BaseExperience || master {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		CaughtPokemon[pokemon.Name] = pokemon

	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}
