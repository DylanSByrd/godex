package main

import (
	"time"

	"github.com/dylansbyrd/godex/internal/pokeapi"
)

func main() {
	context := &commandContext {
		client: pokeapi.NewClient(5 * time.Second, 5 * time.Second),
		pokedex: map[string]pokeapi.PokemonDetails{},
	}
	startRepl(context)
}
