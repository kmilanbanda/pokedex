package main

import (
	"fmt"
)

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
