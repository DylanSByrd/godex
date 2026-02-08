package main

import "strings"

const  (
	pokeballRollMin = 10
	pokeballRollMax = 120
	greatballRollMin = 70
	greatballRollMax = 200
	ultraballRollMin = 150
	ultraballRollMax = 350
	masterballRollMin = 9999
	masterballRollMax = 10000

	greatballRequirement = 10
	ultraballRequirement = 20
	masterballRequirement = 30
)

type PokeballType int
const (
	pokeball PokeballType = iota
	greatball
	ultraball
	masterball
)

var (
	pokeballTypeMap = map[string]PokeballType {
		"pokeball": pokeball,
		"greatball": greatball,
		"ultraball": ultraball,
		"masterball": masterball,
	}
)

func parseString(str string) (PokeballType, bool) {
	pokeball, ok := pokeballTypeMap[strings.ToLower(str)]
	return pokeball, ok
}

func (pokeballType PokeballType) String() string {
	switch pokeballType {
	case pokeball: return "Pokeball"
	case greatball: return "Greatball"
	case ultraball: return "Ultraball"
	case masterball: return "Masterball"
	default: return "Unknown"
	}
}

func (pokeballType PokeballType) getNumPokemonRequiredForUpgrade() int {
	switch pokeballType {
	case pokeball: return greatballRequirement
	case greatball: return ultraballRequirement
	case ultraball: return masterballRequirement
	default: return 0
	}
}

func (pokeballType PokeballType ) getRollRange() (int, int) {
	switch pokeballType {
	case pokeball: return pokeballRollMin, pokeballRollMax
	case greatball: return greatballRollMin, greatballRollMax
	case ultraball: return ultraballRollMin, ultraballRollMax
	case masterball: return masterballRollMin, masterballRollMax 
	default: return 0,0
	}
}

