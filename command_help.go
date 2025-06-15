package main

import (
	"fmt"
)

func commandHelp(c *config, args []string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage: \n\n")
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)	
	}
	return nil
}
