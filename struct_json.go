package main

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
