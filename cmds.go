package main

import (
	"fmt"
	"os"
	"math/rand/v2"

	"github.com/dylansbyrd/godex/internal/pokeapi"
)

const  (
	rollMin = 10
	rollMax = 350
)

type cliCommand struct {
	name string
	description string
	callback func(*commandContext, ...string) error
}

type commandContext struct {
	client pokeapi.Client
	nextLocationArea *string
	prevLocationArea *string
	pokedex map[string]pokeapi.PokemonDetails
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
		"map": {
			name: "map",
			description: "Get the next page of locations",
			callback: commandMapf,
		},
		"mapb": {
			name: "mapb",
			description: "Get the previous page of locations",
			callback: commandMapb,
		},
		"explore": {
			name: "explore <location_name>",
			description: "Explore a location",
			callback: commandExplore,
		},
		"catch": {
			name: "catch <pokemon_name>",
			description: "Attempts to catch a pokemon",
			callback: commandCatch,
		},
		"dev_map": {
			name: "dev_map",
			description: "Dev command: prints the next and previous map locations",
			callback: devCommandPrintLocationContext,
		},
		"dev_pokedex": {
			name: "dev_pokedex",
			description: "Dev command: prints the list of available pokedex entries",
			callback: devCommandPrintPokedex,
		},
	}
}

func commandExit(*commandContext, ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(*commandContext, ...string) error {
	fmt.Println("")
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")

	for _, cmd := range getCommands() {
		fmt.Printf("%v: %v\n", cmd.name, cmd.description)
	}

	return nil
}

func commandMapf(context *commandContext, _ ...string) error {
	resourceList, err := context.client.RequestLocationArea(context.nextLocationArea)
	if err != nil {
		return err
	}
	
	for _, resource := range resourceList.Results {
		fmt.Println(resource.Name)
	}

	context.nextLocationArea = resourceList.Next
	context.prevLocationArea = resourceList.Previous
	return nil
}

func commandMapb(context *commandContext, _ ...string) error {
	if context.prevLocationArea == nil {
		fmt.Println("You're on the first page")
		return nil
	}

	resourceList, err := context.client.RequestLocationArea(context.prevLocationArea)
	if err != nil {
		return err
	}
	
	for _, resource := range resourceList.Results {
		fmt.Println(resource.Name)
	}

	context.nextLocationArea = resourceList.Next
	context.prevLocationArea = resourceList.Previous

	return nil
}

func commandExplore(context *commandContext, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Not enough arguments to 'explore' command. Please provide an area name.")
	}

	fmt.Printf("Exploring %s...\n", args[0])

	locationArea, err := context.client.RequestLocationAreaDetails(args[0])
	if err != nil {
		return err
	}

	if len(locationArea.PokemonEncounters) == 0 {
		fmt.Println("No Pokemon found")
		return nil
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range locationArea.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(context* commandContext, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Nnot enough arguments to 'catch' command. Please provide a pokemon name.")
	}

	pokemon, err := context.client.RequestPokemonDetails(args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	
	roll := rand.IntN(rollMax - rollMin) + rollMin
	fmt.Printf("Rolled %v against %v\n", roll, pokemon.BaseExperience)

	if roll > pokemon.BaseExperience {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		context.pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func devCommandPrintLocationContext(context *commandContext, _ ...string) error {
	var next string
	if context.nextLocationArea != nil {
		next = *context.nextLocationArea
	} else {
		next = "nil"
	}

	var prev string
	if context.prevLocationArea != nil {
		prev = *context.prevLocationArea
	} else {
		prev = "nil"
	}

	fmt.Println("Current context:")
	fmt.Printf("Next Location Area: %v\n", next)
	fmt.Printf("Prev Location Area: %v\n", prev)
	return nil
}

func devCommandPrintPokedex(context *commandContext, _ ...string) error {
	fmt.Println("Current pokedex:")
	for name, _ := range context.pokedex {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}

