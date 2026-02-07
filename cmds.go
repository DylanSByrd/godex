package main

import (
	"fmt"
	"os"

	"github.com/dylansbyrd/godex/internal/pokeapi"
)

type cliCommand struct {
	name string
	description string
	callback func(config *commandContext) error
}

type commandContext struct {
	client pokeapi.Client
	nextLocationArea *string
	prevLocationArea *string
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
		"dev_context": {
			name: "dev_context",
			description: "Dev command: prints the current command context",
			callback: devCommandPrintContext,
		},
	}
}

func commandExit(*commandContext) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(*commandContext) error {
	fmt.Println("")
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")

	for _, cmd := range getCommands() {
		fmt.Printf("%v: %v\n", cmd.name, cmd.description)
	}

	return nil
}

func commandMapf(context *commandContext) error {
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

func commandMapb(context *commandContext) error {
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

func devCommandPrintContext(context *commandContext) error {
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


