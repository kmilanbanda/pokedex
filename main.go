package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"io"
	"encoding/json"
	"net/http"
	"pokedex/internal/pokecache"
	"time"
	"bytes"
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

func commandExit(c *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, args []string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage: \n\n")
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)	
	}
	return nil
}

func commandMap(c *config, args []string) error {
	if c.next == nil {
		fmt.Println("you're on the last page")
		return nil
	}
	
	var body []byte
	var isCached bool
	if body, isCached = Cache.Get(*c.next); !isCached { 
		res, err := http.Get(*c.next)
		if err != nil {
			fmt.Errorf("Error getting response: %v", err)
			return err
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		Cache.Add(*c.next, body)
	}

	var page Page 
	decoder := json.NewDecoder(bytes.NewReader(body))
	if err := decoder.Decode(&page); err != nil {
		return err
	}

	areas := page.Results
	for i := 0; i < len(areas); i++ {
		fmt.Printf("%s\n", areas[i].Name)
	}

	*c = config{
		next: 		page.Next,
		previous: 	page.Previous,
	}	

	return nil
}

func commandMapb(c *config, args []string) error {
	if c.previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}	

	var body []byte
	var isCached bool
	if body, isCached = Cache.Get(*c.previous); !isCached { 
		res, err := http.Get(*c.previous)
		if err != nil {
			fmt.Errorf("Error getting response: %v", err)
			return err
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		Cache.Add(*c.next, body)
	}

	var page Page 
	decoder := json.NewDecoder(bytes.NewReader(body))
	if err := decoder.Decode(&page); err != nil {
		return err
	}
	
	areas := page.Results
	for i := 0; i < len(areas); i++ {
		fmt.Printf("%s\n", areas[i].Name)
	}

	*c = config{
		next: 		page.Next,
		previous: 	page.Previous,
	}

	return nil
}

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

	if roll := rand.Intn(1000); roll > pokemon.BaseExperience {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		CaughtPokemon[pokemon.Name] = pokemon

	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func commandInspect(c *config, args []string) error {
	if pokemon, isCaught := CaughtPokemon[args[0]]; isCaught {
		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("Height: %v\n", pokemon.Height)
		fmt.Printf("Weight: %v\n", pokemon.Weight)

		fmt.Printf("Stats:\n")
		for i := 0; i < len(pokemon.Stats); i++ {
			fmt.Printf("  -%s: %v\n", pokemon.Stats[i].Stat.Name, pokemon.Stats[i].BaseStat)
		}

		fmt.Printf("Types:\n")
		for i := 0; i < len(pokemon.Types); i++ {
			fmt.Printf("  - %s\n", pokemon.Types[i].Details.Name)
		}

	} else {
		fmt.Println("You have not caught this pokemon!")
	}

	return nil
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
