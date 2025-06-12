package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"encoding/json"
	"net/http"
)

type cliCommand struct {
	name 		string
	description 	string
	callback	func(*config) error
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

func cleanInput(text string) []string {
	cleanText := strings.ToLower(text)
	cleanedWords := strings.Fields(cleanText)	
	return cleanedWords
}

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage: \n\n")
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)	
	}
	return nil
}

func commandMap(c *config) error {
	if c.next == nil {
		fmt.Println("you're on the last page")
		return nil
	}

	res, err := http.Get(*c.next)
	if err != nil {
		fmt.Errorf("Error getting response: %v", err)
		return err
	}
	defer res.Body.Close()

	var page Page 
	decoder := json.NewDecoder(res.Body)
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

func commandMapb(c *config) error {
	if c.previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}	

	res, err := http.Get(*c.previous)
	if err != nil {
		fmt.Errorf("Error getting response: %v", err)
		return err
	}
	defer res.Body.Close()

	var page Page 
	decoder := json.NewDecoder(res.Body)
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
	}  
}

var commands map[string]cliCommand

func main() { 
	var c config
	var startUrl string
	startUrl = "https://pokeapi.co/api/v2/location-area/" 
	c = config{}
	c.next = &startUrl

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
			command.callback(&c)
		}
	}
}
