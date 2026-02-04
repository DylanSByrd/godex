package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		tokens := cleanInput(input)
		if len(tokens) != 0 {
			fmt.Printf("Your command was: %v", tokens[0])
		}
		fmt.Println("")
	}
}
