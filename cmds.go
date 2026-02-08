package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"

	"github.com/dylansbyrd/godex/internal/pokeapi"
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
	pokeballType PokeballType

	// debug
	debugCatchRolls bool
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
			description: "Explores a location",
			callback: commandExplore,
		},
		"catch": {
			name: "catch <pokemon_name>",
			description: "Attempts to catch a Pokemon",
			callback: commandTryCatchPokemon,
		},
		"inspect": {
			name: "inspect <pokemon_name>",
			description: "Shows details about a Pokemon if caught",
			callback: commandInspect,
		},
		"pokedex": {
			name: "pokedex",
			description: "Prints the list of inspectable Pokedex entries",
			callback: commandPrintPokedex,
		},
		"pokeball": {
			name: "pokeball",
			description: "Shows details about your current Pokeballs",
			callback: commandPrintPokeballDetails,
		},
		"dev_map": {
			name: "dev_map",
			description: "Dev command: prints the next and previous map endpoints",
			callback: devCommandPrintLocationContext,
		},
		"dev_rolls": {
			name: "dev_rolls",
			description: "Dev command: toggles display of dice roll when attempting to catch Pokemon",
			callback: devCommandShowCatchRolls,
		},
		"dev_pokeball": {
			name: "dev_pokeball <pokeball_type>",
			description: "Dev command: sets your current Pokeball type to <pokeball_type>",
			callback: devCommandForcePokeballType,
		},
		"dev_simulate": {
			name: "dev_simulate <pokemon> <(optional)num_attempts=1>",
			description: "Dev command: simulates <num_attempts> catch attempts against <pokemon>",
			callback: devCommandSimulateMultipleCatchAttempts,
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
	if context.nextLocationArea == nil && context.prevLocationArea != nil {
		fmt.Println("You're on the last page")
		return nil
	}

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

func commandTryCatchPokemon(context* commandContext, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Nnot enough arguments to 'catch' command. Please provide a Pokemon name.")
	}

	pokemon, err := context.client.RequestPokemonDetails(args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a %s at %s...\n", context.pokeballType, pokemon.Name)
	
	rollMin, rollMax := context.pokeballType.getRollRange()
	roll := rand.IntN(rollMax - rollMin) + rollMin

	if context.debugCatchRolls {
		fmt.Printf("Rolling between %v - %v...", rollMin, rollMax)
		fmt.Printf("Rolled %v against %v\n", roll, pokemon.BaseExperience)
	}

	if roll > pokemon.BaseExperience {
		context.pokedex[pokemon.Name] = pokemon
		fmt.Printf("%s was caught!\n", pokemon.Name)
		fmt.Println("You may now inspect it with the inspect command.")

		// Attempt to upgrade pokeball type
		numRequiredForPokeballUpgrade := context.pokeballType.getNumPokemonRequiredForUpgrade()
		if numRequiredForPokeballUpgrade > 0 {
			numPokemonEntries := len(context.pokedex)
			if numPokemonEntries > numRequiredForPokeballUpgrade {
				fmt.Printf("You've upgraded your %ss to %ss!\n", context.pokeballType, context.pokeballType + 1)
				context.pokeballType = context.pokeballType + 1
			}
		}
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func commandInspect(context* commandContext, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Not enough arguments to 'inspect' command. Please provide a Pokemon name.")
	}

	if pokemon, exists := context.pokedex[args[0]]; exists {
		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("Height: %v\n", pokemon.Height)
		fmt.Printf("Weight: %v\n", pokemon.Weight)

		fmt.Println("Stats:")
		for _, pokemonStat := range pokemon.Stats {
			fmt.Printf(" -%s: %v\n", pokemonStat.Stat.Name, pokemonStat.BaseStat)
		}

		fmt.Println("Types:")
		for _, pokemonType := range pokemon.Types {
			fmt.Printf(" -%s\n", pokemonType.Type.Name)
		}
	} else {
		fmt.Println("You have not caught that Pokemon")
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

func commandPrintPokedex(context *commandContext, _ ...string) error {
	fmt.Println("Your Pokedex:")
	for name, _ := range context.pokedex {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}

func commandPrintPokeballDetails(context *commandContext, _ ...string) error {
	fmt.Printf("%s stats:\n", context.pokeballType)

	rollMin, rollMax := context.pokeballType.getRollRange()
	fmt.Printf(" - min catch power: %d\n", rollMin)
	fmt.Printf(" - max catch power: %d\n", rollMax)
	
	numNeededForUpgrade := context.pokeballType.getNumPokemonRequiredForUpgrade()
	if numNeededForUpgrade > 0 {
		numUniqueCaught := len(context.pokedex)
		fmt.Printf(" - upgrades at %v unique pokemon caught (need %v more)\n", 
			numNeededForUpgrade, numNeededForUpgrade - numUniqueCaught)
	} else {
		fmt.Printf("- fully upgraded")
	}

	return nil
}

func devCommandShowCatchRolls(context *commandContext, _ ...string) error {
	context.debugCatchRolls = !context.debugCatchRolls
	if context.debugCatchRolls {
		fmt.Println("Now displaying catch rolls")
	} else {
		fmt.Println("No longer displaying catch rolls")
	}

	return nil
}

func devCommandForcePokeballType(context *commandContext, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Not enough arguments to 'dev_pokeball' command. Please provide a Pokeball type.")
	}

	pokeballType, ok := parseString(args[0])
	if !ok {
		return fmt.Errorf("Invalid Pokeball type\n")
	}

	context.pokeballType = pokeballType
	fmt.Printf("You are now using %ss.\n", context.pokeballType)

	return nil
}

func devCommandSimulateMultipleCatchAttempts(context *commandContext, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Not enough arguments to 'dev_simulate' command. Please provide a Pokemon name. Optionally provide a number of attempts.")
	}

	pokemon, err := context.client.RequestPokemonDetails(args[0])
	if err != nil {
		return err
	}

	numAttempts := 1
	if len(args) > 1 {
		numAttempts, err = strconv.Atoi(args[1])
		if err != nil {
			return err
		}
	}
	fmt.Printf("Simulating %v catch attempt(s) against %s...\n", numAttempts, pokemon.Name)

	numCatches := 0
	for range numAttempts {
		rollMin, rollMax := context.pokeballType.getRollRange()
		roll := rand.IntN(rollMax - rollMin) + rollMin
		if roll > pokemon.BaseExperience {
			numCatches++
		}
	}

	catchPercent := 100.0 * float64(numCatches) / float64(numAttempts)
	fmt.Printf("Result: %v catches out of %v attempts (%v%%)\n", numCatches, numAttempts, catchPercent)
	return nil
}
