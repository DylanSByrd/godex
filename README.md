# Godex
Simple toy Pokedex REPL app made in Go with the PokeAPI with local response caching. Part of boot.dev coursework.

## Requirements
Go version 1.22+

## Usage
`help`: Displays a help message

`catch <pokemon_name>`: Attempts to catch a Pokemon

`pokedex`: Prints the list of inspectable Pokedex entries

`pokeball`: Shows details about your current Pokeballs

`map`: Get the next page of locations

`mapb`: Get the previous page of locations

`explore <location_name>`: Explores a location

`inspect <pokemon_name>`: Shows details about a Pokemon if caught

`dev_map`: Dev command: prints the next and previous map endpoints

`dev_pokeball <pokeball_type>`: Dev command: sets your current Pokeball type to `<pokeball_type>`

`dev_rolls`: Dev command: toggles display of dice roll when attempting to catch Pokemon

`dev_simulate <pokemon> <(optional)num_attempts=1>`: Dev command: simulates `<num_attempts>` catch attempts against `<pokemon>`

`exit`: Exit the Pokedex
