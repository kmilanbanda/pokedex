package main

import (
	"fmt"
)

func commandPokedex(c *config, args []string) error {
	fmt.Println("Your Pokedex:")
	for key, _ := range CaughtPokemon {
		fmt.Printf(" - %s\n", key)
	}
	
	return nil
}
