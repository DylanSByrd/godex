package main

import (
	"time"

	"github.com/dylansbyrd/godex/internal/pokeapi"
)

func main() {
	context := &commandContext {
		client: pokeapi.NewClient(5 * time.Second, 10 * time.Minute),
		pokedex: map[string]pokeapi.PokemonDetails{},
	}
	startRepl(context)
}
