package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"pokedex/internal/pokecache"
	"time"
	"math/rand"
)

type cliCommand struct {
	name 		string
	description 	string
	callback	func(*config, []string) error
}

type config struct {
	next		*string
	previous	*string
}

type Area struct {
	Name		string	`json:"name"`
	Url		string	`json:"url"`
}

type Page struct {
	Count		int	`json:"count"`
	Next		*string	`json:"next"`
	Previous	*string	`json:"previous"`
	Results		[]Area	`json:"results"`	
}

type TypeDetails struct {
	Name		string	`json:"name"`
}

type Type struct {
	Slot		int		`json:"slot"`
	Details		TypeDetails	`json:"type"`
}

type StatDetails struct {
	Name		string `json:"name"`
}

type Stat struct {
	BaseStat	int		`json:"base_stat"`	
	Stat		StatDetails	`json:"stat"`
}

type Pokemon struct {
	Name		string 	`json:"name"`
	BaseExperience	int	`json:"base_experience"`
	Weight		int	`json:"weight"`
	Height		int	`json:"height"`
	Stats		[]Stat	`json:"stats"`
	Types		[]Type	`json:"types"`
}

type PokemonDetails struct {
	Name		string	`json:"name"`
	Url		string	`json:"url"`
}

type PokemonEncounter struct {
	Pokemon		PokemonDetails	`json:"pokemon"`
	//VersionDetails	
}

type AreaEncounters struct {
	Encounters	[]PokemonEncounter	`json:"pokemon_encounters"`	
}

func cleanInput(text string) []string {
	cleanText := strings.ToLower(text)
	cleanedWords := strings.Fields(cleanText)	
	return cleanedWords
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:		"exit",
			description:	"Exit the Pokedex",
			callback:	commandExit,
		},
		"help": {
			name:		"help",
			description:	"Display a help message",
			callback:	commandHelp,
		},
		"map": {
			name:  		"map",
			description: 	"Displays the next 20 location areas in Pokemon",
			callback: 	commandMap,
		},
		"mapb": {
			name:  		"mapb",
			description: 	"Displays the previous 20 location areas in Pokemon",
			callback: 	commandMapb,
		},
		"explore": {
			name:		"explore",
			description:	"Displays a list of all pokemon in a location",
			callback: 	commandExplore,
		},
		"catch": {
			name:		"catch",
			description:	"Try to catch a pokemon",
			callback:	commandCatch,
		},
		"inspect": {
			name:		"inspect",
			description:	"Inspects a pokemon you have caught",
			callback:	commandInspect,
		},
		"pokedex": {
			name:		"pokedex",
			description:	"Displays a list of all the pokemon you have caught",
			callback:	commandPokedex,
		},
	}  
}

var Cache *pokecache.Cache
var CaughtPokemon map[string]Pokemon

func main() { 
	var c config
	var startUrl string
	startUrl = "https://pokeapi.co/api/v2/location-area/" 
	c = config{}
	c.next = &startUrl

	Cache = pokecache.NewCache(5 * time.Second)
	CaughtPokemon = make(map[string]Pokemon)
	rand.Seed(time.Now().UnixNano())

	var userInputTokens []string
	scanner := bufio.NewScanner(os.Stdin)
	for ;; {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInputTokens = cleanInput(scanner.Text())
		command, ok := getCommands()[userInputTokens[0]]
		if !ok {
			fmt.Println("Unknown command")
		} else {
			command.callback(&c, userInputTokens[1:])
		}
	}
}
